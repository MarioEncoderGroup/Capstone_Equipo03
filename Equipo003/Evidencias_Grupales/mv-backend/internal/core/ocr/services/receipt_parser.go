package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParsedReceipt representa la información estructurada extraída de un recibo chileno
type ParsedReceipt struct {
	Amount        float64            `json:"amount"`
	Date          *time.Time         `json:"date,omitempty"`
	MerchantRUT   string             `json:"merchant_rut,omitempty"`
	MerchantName  string             `json:"merchant_name,omitempty"`
	DocumentType  string             `json:"document_type,omitempty"` // boleta, factura, ticket
	Confidence    float64            `json:"confidence"`              // 0.0 - 1.0
	RawText       string             `json:"raw_text"`
	ExtractedData map[string]float64 `json:"extracted_data,omitempty"` // datos adicionales con confianza
}

// ReceiptParser parsea texto OCR para extraer información de recibos chilenos
type ReceiptParser struct {
	// Regex patterns para extracción de datos
	amountRegex       *regexp.Regexp
	rutRegex          *regexp.Regexp
	dateRegexes       []*regexp.Regexp
	merchantRegex     *regexp.Regexp
	documentTypeRegex *regexp.Regexp
}

// NewReceiptParser crea una nueva instancia del parser
func NewReceiptParser() *ReceiptParser {
	return &ReceiptParser{
		// Montos en formato chileno: $15.000 o 15000 o $15000
		amountRegex: regexp.MustCompile(`\$?\s*(\d{1,3}(?:\.\d{3})+|\d+)`),

		// RUT chileno: 12.345.678-9 o 12345678-9
		rutRegex: regexp.MustCompile(`(\d{1,2}\.?\d{3}\.?\d{3}-[\dkK])`),

		// Múltiples formatos de fecha
		dateRegexes: []*regexp.Regexp{
			regexp.MustCompile(`(\d{1,2})[/-](\d{1,2})[/-](\d{2,4})`),                                             // DD/MM/YYYY o DD-MM-YYYY
			regexp.MustCompile(`(\d{1,2})\s+de\s+(enero|febrero|marzo|abril|mayo|junio|julio|agosto|septiembre|octubre|noviembre|diciembre)\s+de\s+(\d{4})`), // 15 de Octubre de 2024
			regexp.MustCompile(`(\d{4})[/-](\d{1,2})[/-](\d{1,2})`),                                             // YYYY-MM-DD
		},

		// Nombre del comercio (primeras líneas del recibo)
		merchantRegex: regexp.MustCompile(`^([A-ZÁÉÍÓÚÑ][A-ZÁÉÍÓÚÑa-záéíóúñ\s\.]+)`),

		// Tipo de documento
		documentTypeRegex: regexp.MustCompile(`(?i)(boleta|factura|ticket|comprobante)`),
	}
}

// ParseChileanReceipt parsea un texto OCR completo de un recibo chileno
func (p *ReceiptParser) ParseChileanReceipt(text string) (*ParsedReceipt, error) {
	if text == "" {
		return nil, fmt.Errorf("texto vacío")
	}

	receipt := &ParsedReceipt{
		RawText:       text,
		ExtractedData: make(map[string]float64),
	}

	// Extraer monto total
	amount, amountConf, err := p.ExtractAmount(text)
	if err == nil {
		receipt.Amount = amount
		receipt.ExtractedData["amount_confidence"] = amountConf
	}

	// Extraer fecha
	date, dateConf, err := p.ExtractDate(text)
	if err == nil {
		receipt.Date = date
		receipt.ExtractedData["date_confidence"] = dateConf
	}

	// Extraer RUT del comercio
	rut, rutConf, err := p.ExtractRUT(text)
	if err == nil {
		receipt.MerchantRUT = rut
		receipt.ExtractedData["rut_confidence"] = rutConf
	}

	// Extraer nombre del comercio
	merchantName, merchantConf := p.ExtractMerchantName(text)
	if merchantName != "" {
		receipt.MerchantName = merchantName
		receipt.ExtractedData["merchant_confidence"] = merchantConf
	}

	// Detectar tipo de documento
	documentType, docConf := p.DetectDocumentType(text)
	if documentType != "" {
		receipt.DocumentType = documentType
		receipt.ExtractedData["document_type_confidence"] = docConf
	}

	// Calcular confianza general
	receipt.Confidence = p.CalculateOverallConfidence(receipt)

	return receipt, nil
}

// ExtractAmount extrae el monto total del recibo
func (p *ReceiptParser) ExtractAmount(text string) (float64, float64, error) {
	lines := strings.Split(text, "\n")

	// Buscar palabras clave para el total
	keywords := []string{"total", "total a pagar", "total $", "monto total", "importe"}

	var bestAmount float64
	var bestConfidence float64

	for _, line := range lines {
		lineLower := strings.ToLower(line)

		// Verificar si la línea contiene una palabra clave
		containsKeyword := false
		for _, keyword := range keywords {
			if strings.Contains(lineLower, keyword) {
				containsKeyword = true
				break
			}
		}

		// Extraer montos de la línea
		matches := p.amountRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				amountStr := strings.ReplaceAll(match[1], ".", "") // Remover puntos de miles
				amount, err := strconv.ParseFloat(amountStr, 64)
				if err != nil {
					continue
				}

				// Calcular confianza basada en contexto
				confidence := 0.5
				if containsKeyword {
					confidence = 0.9
				}

				// Preferir montos mayores con alta confianza
				if amount > bestAmount && confidence >= bestConfidence {
					bestAmount = amount
					bestConfidence = confidence
				}
			}
		}
	}

	if bestAmount == 0 {
		return 0, 0, fmt.Errorf("no se encontró monto en el recibo")
	}

	return bestAmount, bestConfidence, nil
}

// ExtractDate extrae la fecha del recibo
func (p *ReceiptParser) ExtractDate(text string) (*time.Time, float64, error) {
	// Mapeo de meses en español
	monthMap := map[string]time.Month{
		"enero":      time.January,
		"febrero":    time.February,
		"marzo":      time.March,
		"abril":      time.April,
		"mayo":       time.May,
		"junio":      time.June,
		"julio":      time.July,
		"agosto":     time.August,
		"septiembre": time.September,
		"octubre":    time.October,
		"noviembre":  time.November,
		"diciembre":  time.December,
	}

	// Intentar con cada patrón de fecha
	for i, dateRegex := range p.dateRegexes {
		matches := dateRegex.FindStringSubmatch(strings.ToLower(text))
		if len(matches) > 1 {
			var date time.Time
			var err error

			switch i {
			case 0: // DD/MM/YYYY
				day, _ := strconv.Atoi(matches[1])
				month, _ := strconv.Atoi(matches[2])
				year, _ := strconv.Atoi(matches[3])
				if year < 100 {
					year += 2000
				}
				date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

			case 1: // 15 de Octubre de 2024
				day, _ := strconv.Atoi(matches[1])
				month := monthMap[matches[2]]
				year, _ := strconv.Atoi(matches[3])
				date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

			case 2: // YYYY-MM-DD
				year, _ := strconv.Atoi(matches[1])
				month, _ := strconv.Atoi(matches[2])
				day, _ := strconv.Atoi(matches[3])
				date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			}

			if err == nil && !date.IsZero() {
				confidence := 0.8 - float64(i)*0.2 // Mayor confianza para formatos más comunes
				return &date, confidence, nil
			}
		}
	}

	return nil, 0, fmt.Errorf("no se encontró fecha válida")
}

// ExtractRUT extrae y valida el RUT del comercio
func (p *ReceiptParser) ExtractRUT(text string) (string, float64, error) {
	matches := p.rutRegex.FindAllString(text, -1)

	for _, match := range matches {
		// Limpiar y formatear RUT
		rut := strings.ReplaceAll(match, ".", "")
		rut = strings.ReplaceAll(rut, " ", "")

		// Validar RUT
		if p.ValidateRUT(rut) {
			formattedRUT := p.FormatRUT(rut)
			return formattedRUT, 0.9, nil
		}
	}

	return "", 0, fmt.Errorf("no se encontró RUT válido")
}

// ValidateRUT valida un RUT chileno usando el algoritmo de verificación
func (p *ReceiptParser) ValidateRUT(rut string) bool {
	// Remover formato
	rut = strings.ReplaceAll(rut, ".", "")
	rut = strings.ReplaceAll(rut, "-", "")

	if len(rut) < 2 {
		return false
	}

	// Separar dígito verificador
	body := rut[:len(rut)-1]
	dv := strings.ToUpper(rut[len(rut)-1:])

	// Convertir cuerpo a número
	num, err := strconv.Atoi(body)
	if err != nil {
		return false
	}

	// Calcular dígito verificador
	sum := 0
	multiplier := 2

	for num > 0 {
		sum += (num % 10) * multiplier
		num /= 10
		multiplier++
		if multiplier > 7 {
			multiplier = 2
		}
	}

	expectedDV := 11 - (sum % 11)
	var expectedDVStr string

	switch expectedDV {
	case 11:
		expectedDVStr = "0"
	case 10:
		expectedDVStr = "K"
	default:
		expectedDVStr = strconv.Itoa(expectedDV)
	}

	return dv == expectedDVStr
}

// FormatRUT formatea un RUT al formato estándar XX.XXX.XXX-X
func (p *ReceiptParser) FormatRUT(rut string) string {
	// Remover formato existente
	rut = strings.ReplaceAll(rut, ".", "")
	rut = strings.ReplaceAll(rut, "-", "")

	if len(rut) < 2 {
		return rut
	}

	body := rut[:len(rut)-1]
	dv := rut[len(rut)-1:]

	// Agregar puntos cada 3 dígitos desde el final
	var formatted strings.Builder
	for i, digit := range body {
		if i > 0 && (len(body)-i)%3 == 0 {
			formatted.WriteString(".")
		}
		formatted.WriteRune(digit)
	}

	return formatted.String() + "-" + strings.ToUpper(dv)
}

// ExtractMerchantName extrae el nombre del comercio (usualmente en las primeras líneas)
func (p *ReceiptParser) ExtractMerchantName(text string) (string, float64) {
	lines := strings.Split(text, "\n")

	// Buscar en las primeras 5 líneas
	for i := 0; i < len(lines) && i < 5; i++ {
		line := strings.TrimSpace(lines[i])

		// Saltar líneas muy cortas o que parecen fechas/números
		if len(line) < 3 || regexp.MustCompile(`^\d+$`).MatchString(line) {
			continue
		}

		// Buscar nombres de comercio (texto en mayúsculas o capitalizado)
		if match := p.merchantRegex.FindString(line); match != "" {
			confidence := 0.7
			if i == 0 {
				confidence = 0.9 // Mayor confianza si está en la primera línea
			}
			return strings.TrimSpace(match), confidence
		}
	}

	return "", 0
}

// DetectDocumentType detecta el tipo de documento fiscal
func (p *ReceiptParser) DetectDocumentType(text string) (string, float64) {
	textLower := strings.ToLower(text)

	documentTypes := map[string]float64{
		"factura":     0.95,
		"boleta":      0.95,
		"ticket":      0.80,
		"comprobante": 0.70,
	}

	for docType, confidence := range documentTypes {
		if strings.Contains(textLower, docType) {
			return docType, confidence
		}
	}

	return "", 0
}

// CalculateOverallConfidence calcula la confianza general del parsing
func (p *ReceiptParser) CalculateOverallConfidence(receipt *ParsedReceipt) float64 {
	if receipt.ExtractedData == nil {
		return 0.0
	}

	// Pesos para cada campo
	weights := map[string]float64{
		"amount_confidence":         0.35,
		"date_confidence":           0.20,
		"rut_confidence":            0.20,
		"merchant_confidence":       0.15,
		"document_type_confidence":  0.10,
	}

	totalWeight := 0.0
	weightedSum := 0.0

	for field, weight := range weights {
		if confidence, exists := receipt.ExtractedData[field]; exists && confidence > 0 {
			weightedSum += confidence * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0.0
	}

	return weightedSum / totalWeight
}
