// Package quoteratelimit implements sliding-window rate limits for anonymous quote requests (US-2).
package quoteratelimit

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// WindowLimiter counts events per key within a rolling time window.
type WindowLimiter struct {
	mu     sync.Mutex
	byKey  map[string][]time.Time
	max    int
	window time.Duration
}

// NewWindowLimiter allows at most max events per key per window duration.
func NewWindowLimiter(max int, window time.Duration) *WindowLimiter {
	return &WindowLimiter{byKey: make(map[string][]time.Time), max: max, window: window}
}

// Allow returns false if the key has reached max events in the current window.
func (w *WindowLimiter) Allow(key string) bool {
	if w == nil || w.max <= 0 || key == "" {
		return true
	}
	now := time.Now()
	cutoff := now.Add(-w.window)

	w.mu.Lock()
	defer w.mu.Unlock()

	hits := w.byKey[key]
	out := hits[:0]
	for _, t := range hits {
		if t.After(cutoff) {
			out = append(out, t)
		}
	}
	if len(out) >= w.max {
		w.byKey[key] = out
		return false
	}
	out = append(out, now)
	w.byKey[key] = out
	return true
}

// EnvInt parses os.Getenv(key) as int, or returns def.
func EnvInt(key string, def int) int {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return def
	}
	return n
}

// EnvDurationMinutes parses minutes from env, or returns def.
func EnvDurationMinutes(key string, def time.Duration) time.Duration {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return def
	}
	return time.Duration(n) * time.Minute
}
