// Package httputil holds small HTTP helpers shared by handlers.
package httputil

import (
	"net"
	"net/http"
	"strings"
)

// ClientIP returns the best-effort client IP (first X-Forwarded-For hop, else RemoteAddr host).
func ClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	h, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return h
}
