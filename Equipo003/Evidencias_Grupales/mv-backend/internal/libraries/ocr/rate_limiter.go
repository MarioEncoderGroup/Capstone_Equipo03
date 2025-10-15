package ocr

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RateLimiter implementa rate limiting para Google Vision API
// Google Vision tiene límites: 1800 requests/min, 10 requests/seg
type RateLimiter struct {
	requestsPerMinute int
	requestsPerSecond int
	
	minuteWindow []time.Time
	secondWindow []time.Time
	mu           sync.Mutex
}

// NewRateLimiter crea un nuevo rate limiter
func NewRateLimiter(requestsPerMinute, requestsPerSecond int) *RateLimiter {
	if requestsPerMinute == 0 {
		requestsPerMinute = 1800 // Límite de Google Vision
	}
	if requestsPerSecond == 0 {
		requestsPerSecond = 10 // Límite de Google Vision
	}

	return &RateLimiter{
		requestsPerMinute: requestsPerMinute,
		requestsPerSecond: requestsPerSecond,
		minuteWindow:      make([]time.Time, 0),
		secondWindow:      make([]time.Time, 0),
	}
}

// Wait espera si es necesario para cumplir con rate limits
func (rl *RateLimiter) Wait(ctx context.Context) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Limpiar ventanas antiguas
	rl.cleanOldRequests(now)

	// Verificar límite por segundo
	if len(rl.secondWindow) >= rl.requestsPerSecond {
		// Esperar hasta que pase 1 segundo desde la primera request
		waitUntil := rl.secondWindow[0].Add(1 * time.Second)
		waitDuration := time.Until(waitUntil)
		
		if waitDuration > 0 {
			rl.mu.Unlock()
			select {
			case <-time.After(waitDuration):
				rl.mu.Lock()
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		
		// Limpiar nuevamente después de esperar
		rl.cleanOldRequests(time.Now())
	}

	// Verificar límite por minuto
	if len(rl.minuteWindow) >= rl.requestsPerMinute {
		// Esperar hasta que pase 1 minuto desde la primera request
		waitUntil := rl.minuteWindow[0].Add(1 * time.Minute)
		waitDuration := time.Until(waitUntil)
		
		if waitDuration > 0 {
			rl.mu.Unlock()
			select {
			case <-time.After(waitDuration):
				rl.mu.Lock()
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		
		// Limpiar nuevamente después de esperar
		rl.cleanOldRequests(time.Now())
	}

	// Registrar nueva request
	now = time.Now()
	rl.secondWindow = append(rl.secondWindow, now)
	rl.minuteWindow = append(rl.minuteWindow, now)

	return nil
}

// cleanOldRequests limpia requests antiguas de las ventanas
func (rl *RateLimiter) cleanOldRequests(now time.Time) {
	// Limpiar ventana de segundo (requests > 1 segundo)
	cutoffSecond := now.Add(-1 * time.Second)
	newSecondWindow := make([]time.Time, 0)
	for _, t := range rl.secondWindow {
		if t.After(cutoffSecond) {
			newSecondWindow = append(newSecondWindow, t)
		}
	}
	rl.secondWindow = newSecondWindow

	// Limpiar ventana de minuto (requests > 1 minuto)
	cutoffMinute := now.Add(-1 * time.Minute)
	newMinuteWindow := make([]time.Time, 0)
	for _, t := range rl.minuteWindow {
		if t.After(cutoffMinute) {
			newMinuteWindow = append(newMinuteWindow, t)
		}
	}
	rl.minuteWindow = newMinuteWindow
}

// GetStats retorna estadísticas de uso
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rl.cleanOldRequests(now)

	return map[string]interface{}{
		"requests_last_second": len(rl.secondWindow),
		"requests_last_minute": len(rl.minuteWindow),
		"limit_per_second":     rl.requestsPerSecond,
		"limit_per_minute":     rl.requestsPerMinute,
		"capacity_second":      fmt.Sprintf("%d%%", (len(rl.secondWindow)*100)/rl.requestsPerSecond),
		"capacity_minute":      fmt.Sprintf("%d%%", (len(rl.minuteWindow)*100)/rl.requestsPerMinute),
	}
}

// RateLimitedVisionClient envuelve GoogleVisionClient con rate limiting
type RateLimitedVisionClient struct {
	*GoogleVisionClient
	limiter *RateLimiter
}

// NewRateLimitedVisionClient crea un cliente con rate limiting
func NewRateLimitedVisionClient(ctx context.Context, cfg GoogleVisionConfig) (*RateLimitedVisionClient, error) {
	client, err := NewGoogleVisionClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &RateLimitedVisionClient{
		GoogleVisionClient: client,
		limiter:            NewRateLimiter(1800, 10), // Límites de Google Vision
	}, nil
}

// DetectText con rate limiting
func (c *RateLimitedVisionClient) DetectText(ctx context.Context, imageData []byte) (string, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limit wait cancelled: %w", err)
	}
	return c.GoogleVisionClient.DetectText(ctx, imageData)
}

// AnalyzeReceipt con rate limiting
func (c *RateLimitedVisionClient) AnalyzeReceipt(ctx context.Context, imageData []byte) (*OCRResult, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait cancelled: %w", err)
	}
	return c.GoogleVisionClient.AnalyzeReceipt(ctx, imageData)
}

// GetStats retorna estadísticas del rate limiter
func (c *RateLimitedVisionClient) GetStats() map[string]interface{} {
	return c.limiter.GetStats()
}
