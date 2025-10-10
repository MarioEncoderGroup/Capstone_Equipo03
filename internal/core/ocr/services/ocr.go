package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/libraries/ocr"
	"github.com/redis/go-redis/v9"
)

// ocrService implementa el servicio de OCR
type ocrService struct {
	visionClient *ocr.GoogleVisionClient
	parser       *ReceiptParser
	redisClient  *redis.Client
	cacheEnabled bool
	cacheTTL     time.Duration
}

// OCRServiceConfig contiene la configuración del servicio OCR
type OCRServiceConfig struct {
	VisionClient *ocr.GoogleVisionClient
	RedisClient  *redis.Client
	CacheEnabled bool
	CacheTTL     time.Duration
}

// NewOCRService crea una nueva instancia del servicio OCR
func NewOCRService(cfg OCRServiceConfig) *ocrService {
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 24 * time.Hour
	}

	return &ocrService{
		visionClient: cfg.VisionClient,
		parser:       NewReceiptParser(),
		redisClient:  cfg.RedisClient,
		cacheEnabled: cfg.CacheEnabled,
		cacheTTL:     cfg.CacheTTL,
	}
}

// ProcessReceipt procesa una imagen de recibo y retorna datos estructurados
func (s *ocrService) ProcessReceipt(ctx context.Context, imageData []byte) (*ParsedReceipt, error) {
	if len(imageData) == 0 {
		return nil, fmt.Errorf("imageData no puede estar vacía")
	}

	// Generar hash de la imagen para cache
	cacheKey := s.generateCacheKey(imageData)

	// Intentar obtener desde cache
	if s.cacheEnabled && s.redisClient != nil {
		cachedResult, err := s.getFromCache(ctx, cacheKey)
		if err == nil && cachedResult != nil {
			return cachedResult, nil
		}
	}

	// 1. Extraer texto usando Google Vision API
	ocrResult, err := s.visionClient.AnalyzeReceipt(ctx, imageData)
	if err != nil {
		return nil, fmt.Errorf("error en OCR: %w", err)
	}

	if ocrResult.FullText == "" {
		return &ParsedReceipt{
			RawText:       "",
			Confidence:    0,
			ExtractedData: make(map[string]float64),
		}, nil
	}

	// 2. Parsear texto extraído con el parser chileno
	parsedReceipt, err := s.parser.ParseChileanReceipt(ocrResult.FullText)
	if err != nil {
		return nil, fmt.Errorf("error parseando recibo: %w", err)
	}

	// Agregar confianza del OCR al resultado
	if parsedReceipt.ExtractedData == nil {
		parsedReceipt.ExtractedData = make(map[string]float64)
	}
	parsedReceipt.ExtractedData["ocr_confidence"] = float64(ocrResult.Confidence)

	// Recalcular confianza general incluyendo OCR
	parsedReceipt.Confidence = s.calculateCombinedConfidence(parsedReceipt, float64(ocrResult.Confidence))

	// Guardar en cache
	if s.cacheEnabled && s.redisClient != nil {
		if err := s.saveToCache(ctx, cacheKey, parsedReceipt); err != nil {
			// Log error pero no fallar la operación
		}
	}

	return parsedReceipt, nil
}

// ProcessReceiptFromURL descarga imagen desde URL y procesa
func (s *ocrService) ProcessReceiptFromURL(ctx context.Context, imageURL string) (*ParsedReceipt, error) {
	if imageURL == "" {
		return nil, fmt.Errorf("imageURL no puede estar vacía")
	}

	// Crear request HTTP con contexto
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	// Realizar request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error descargando imagen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error HTTP %d descargando imagen", resp.StatusCode)
	}

	// Leer imagen completa (límite 10MB)
	limitedReader := io.LimitReader(resp.Body, 10*1024*1024)
	imageData, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("error leyendo imagen: %w", err)
	}

	// Procesar imagen descargada
	return s.ProcessReceipt(ctx, imageData)
}

// generateCacheKey genera una clave de cache basada en el hash de la imagen
func (s *ocrService) generateCacheKey(imageData []byte) string {
	// Usar primeros y últimos 32 bytes como identificador simple
	prefix := "ocr:receipt:"
	if len(imageData) < 64 {
		return fmt.Sprintf("%s%d", prefix, len(imageData))
	}
	return fmt.Sprintf("%s%x%x", prefix, imageData[:32], imageData[len(imageData)-32:])
}

// getFromCache obtiene resultado desde Redis
func (s *ocrService) getFromCache(ctx context.Context, key string) (*ParsedReceipt, error) {
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var result ParsedReceipt
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// saveToCache guarda resultado en Redis
func (s *ocrService) saveToCache(ctx context.Context, key string, result *ParsedReceipt) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, data, s.cacheTTL).Err()
}

// calculateCombinedConfidence calcula confianza combinada del OCR y parsing
func (s *ocrService) calculateCombinedConfidence(receipt *ParsedReceipt, ocrConfidence float64) float64 {
	// Confianza del parser (calculada anteriormente)
	parserConfidence := receipt.Confidence

	// Promedio ponderado: 60% OCR + 40% Parser
	// OCR tiene más peso porque si el OCR falla, el parser no puede funcionar bien
	combined := (ocrConfidence * 0.6) + (parserConfidence * 0.4)

	return combined
}
