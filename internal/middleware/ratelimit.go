package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/IrsanaAhmad/go-starter-kit/shared/response"
	"github.com/gofiber/fiber/v2"
)

type visitor struct {
	tokens    float64
	lastVisit time.Time
}

type rateLimiterStore struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     float64 // tokens yang ditambahkan per detik
	burst    int     // kapasitas maksimal token
}

func newRateLimiterStore(rate float64, burst int) *rateLimiterStore {
	store := &rateLimiterStore{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}

	// Goroutine untuk membersihkan entry yang sudah expired (setiap 3 menit)
	go store.cleanup()

	return store
}

func (s *rateLimiterStore) cleanup() {
	for {
		time.Sleep(3 * time.Minute)
		s.mu.Lock()
		for ip, v := range s.visitors {
			if time.Since(v.lastVisit) > 3*time.Minute {
				delete(s.visitors, ip)
			}
		}
		s.mu.Unlock()
	}
}

// allow mengecek apakah visitor dengan key tertentu masih diizinkan (token bucket algorithm)
func (s *rateLimiterStore) allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, exists := s.visitors[key]
	now := time.Now()

	if !exists {
		s.visitors[key] = &visitor{
			tokens:    float64(s.burst) - 1,
			lastVisit: now,
		}
		return true
	}

	// Tambahkan token berdasarkan waktu yang berlalu
	elapsed := now.Sub(v.lastVisit).Seconds()
	v.tokens += elapsed * s.rate
	if v.tokens > float64(s.burst) {
		v.tokens = float64(s.burst)
	}
	v.lastVisit = now

	if v.tokens >= 1 {
		v.tokens--
		return true
	}

	return false
}

// RateLimiter mengembalikan middleware yang membatasi request per IP menggunakan token bucket.
// Parameter:
//   - maxRequests: jumlah maksimal request yang diizinkan dalam periode (burst capacity)
//   - window: durasi periode waktu untuk pengisian ulang token
//
// Contoh: RateLimiter(60, 1*time.Minute) = 60 request per menit per IP
func RateLimiter(maxRequests int, window time.Duration) fiber.Handler {
	rate := float64(maxRequests) / window.Seconds()
	store := newRateLimiterStore(rate, maxRequests)

	return func(c *fiber.Ctx) error {
		ip := c.IP()

		if !store.allow(ip) {
			return response.Error(c, http.StatusTooManyRequests, "rate limit exceeded, please try again later", nil)
		}

		return c.Next()
	}
}
