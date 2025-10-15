package ocr

import (
	"context"
	"fmt"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"google.golang.org/api/option"
)

// GoogleVisionClient maneja las operaciones con Google Cloud Vision API
type GoogleVisionClient struct {
	client    *vision.ImageAnnotatorClient
	projectID string
}

// GoogleVisionConfig contiene la configuración para Google Vision
type GoogleVisionConfig struct {
	ProjectID           string
	CredentialsFilePath string // Ruta al archivo JSON de credenciales
}

// NewGoogleVisionClient crea una nueva instancia del cliente Vision
func NewGoogleVisionClient(ctx context.Context, cfg GoogleVisionConfig) (*GoogleVisionClient, error) {
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("project ID es requerido")
	}
	if cfg.CredentialsFilePath == "" {
		return nil, fmt.Errorf("credentials file path es requerido")
	}

	// Crear cliente Vision con credenciales del archivo
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(cfg.CredentialsFilePath))
	if err != nil {
		return nil, fmt.Errorf("error creando cliente Vision: %w", err)
	}

	return &GoogleVisionClient{
		client:    client,
		projectID: cfg.ProjectID,
	}, nil
}

// Close cierra el cliente Vision
func (c *GoogleVisionClient) Close() error {
	return c.client.Close()
}

// DetectText detecta texto en una imagen usando OCR
// Retorna el texto completo detectado
func (c *GoogleVisionClient) DetectText(ctx context.Context, imageData []byte) (string, error) {
	if len(imageData) == 0 {
		return "", fmt.Errorf("imageData no puede estar vacía")
	}

	// Crear imagen desde bytes
	image := &visionpb.Image{
		Content: imageData,
	}

	// Crear request para text detection
	request := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type: visionpb.Feature_TEXT_DETECTION,
			},
		},
	}

	// Anotar imagen
	batch := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{request},
	}

	response, err := c.client.BatchAnnotateImages(ctx, batch)
	if err != nil {
		return "", fmt.Errorf("error detectando texto: %w", err)
	}

	if len(response.Responses) == 0 {
		return "", nil
	}

	annotations := response.Responses[0].TextAnnotations
	if len(annotations) == 0 {
		return "", nil // No se detectó texto
	}

	// La primera anotación contiene todo el texto
	return annotations[0].Description, nil
}

// DetectTextDetailed detecta texto con información detallada de cada palabra/bloque
// Retorna anotaciones completas para procesamiento avanzado
func (c *GoogleVisionClient) DetectTextDetailed(ctx context.Context, imageData []byte) ([]*visionpb.EntityAnnotation, error) {
	if len(imageData) == 0 {
		return nil, fmt.Errorf("imageData no puede estar vacía")
	}

	image := &visionpb.Image{
		Content: imageData,
	}

	request := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type:       visionpb.Feature_TEXT_DETECTION,
				MaxResults: 100,
			},
		},
	}

	batch := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{request},
	}

	response, err := c.client.BatchAnnotateImages(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("error detectando texto detallado: %w", err)
	}

	if len(response.Responses) == 0 {
		return nil, nil
	}

	return response.Responses[0].TextAnnotations, nil
}

// DetectDocument detecta texto usando Document Text Detection (mejor para documentos densos)
// Ideal para facturas, boletas, recibos con mucho texto
func (c *GoogleVisionClient) DetectDocument(ctx context.Context, imageData []byte) (*visionpb.TextAnnotation, error) {
	if len(imageData) == 0 {
		return nil, fmt.Errorf("imageData no puede estar vacía")
	}

	image := &visionpb.Image{
		Content: imageData,
	}

	request := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
			},
		},
	}

	batch := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{request},
	}

	response, err := c.client.BatchAnnotateImages(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("error detectando documento: %w", err)
	}

	if len(response.Responses) == 0 {
		return nil, nil
	}

	return response.Responses[0].FullTextAnnotation, nil
}

// OCRResult estructura el resultado del OCR
type OCRResult struct {
	FullText   string                   `json:"full_text"`
	Pages      []*visionpb.Page         `json:"pages,omitempty"`
	Confidence float32                  `json:"confidence"`
	Language   string                   `json:"language,omitempty"`
	RawData    *visionpb.TextAnnotation `json:"-"` // No serializar
}

// AnalyzeReceipt analiza un comprobante/recibo y retorna resultado estructurado
func (c *GoogleVisionClient) AnalyzeReceipt(ctx context.Context, imageData []byte) (*OCRResult, error) {
	// Usar Document Text Detection para mejor precisión en documentos
	annotation, err := c.DetectDocument(ctx, imageData)
	if err != nil {
		return nil, err
	}

	if annotation == nil || annotation.Text == "" {
		return &OCRResult{
			FullText:   "",
			Confidence: 0,
		}, nil
	}

	// Calcular confianza promedio
	var totalConfidence float32
	var count int
	if len(annotation.Pages) > 0 {
		for _, page := range annotation.Pages {
			for _, block := range page.Blocks {
				totalConfidence += block.Confidence
				count++
			}
		}
	}

	avgConfidence := float32(0)
	if count > 0 {
		avgConfidence = totalConfidence / float32(count)
	}

	// Detectar idioma (si está disponible)
	language := ""
	if len(annotation.Pages) > 0 && annotation.Pages[0].Property != nil && len(annotation.Pages[0].Property.DetectedLanguages) > 0 {
		language = annotation.Pages[0].Property.DetectedLanguages[0].LanguageCode
	}

	return &OCRResult{
		FullText:   annotation.Text,
		Pages:      annotation.Pages,
		Confidence: avgConfidence,
		Language:   language,
		RawData:    annotation,
	}, nil
}

// DetectLabels detecta etiquetas/categorías en la imagen (opcional, útil para clasificación)
func (c *GoogleVisionClient) DetectLabels(ctx context.Context, imageData []byte, maxResults int) ([]*visionpb.EntityAnnotation, error) {
	if len(imageData) == 0 {
		return nil, fmt.Errorf("imageData no puede estar vacía")
	}
	if maxResults <= 0 {
		maxResults = 10
	}

	image := &visionpb.Image{
		Content: imageData,
	}

	request := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type:       visionpb.Feature_LABEL_DETECTION,
				MaxResults: int32(maxResults),
			},
		},
	}

	batch := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{request},
	}

	response, err := c.client.BatchAnnotateImages(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("error detectando etiquetas: %w", err)
	}

	if len(response.Responses) == 0 {
		return nil, nil
	}

	return response.Responses[0].LabelAnnotations, nil
}
