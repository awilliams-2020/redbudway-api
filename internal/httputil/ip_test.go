package httputil

import (
	"net/http"
	"testing"
)

func TestClientIP_nil(t *testing.T) {
	if got := ClientIP(nil); got != "" {
		t.Fatalf("ClientIP(nil) = %q, want empty", got)
	}
}

func TestClientIP_xForwardedFor(t *testing.T) {
	r := &http.Request{
		Header: http.Header{"X-Forwarded-For": {"203.0.113.1, 10.0.0.1"}},
		RemoteAddr: "192.168.1.1:12345",
	}
	if got := ClientIP(r); got != "203.0.113.1" {
		t.Fatalf("got %q", got)
	}
}

func TestClientIP_remoteAddrIPv4(t *testing.T) {
	r := &http.Request{RemoteAddr: "198.51.100.2:5000"}
	if got := ClientIP(r); got != "198.51.100.2" {
		t.Fatalf("got %q", got)
	}
}

func TestClientIP_remoteAddrNoPort(t *testing.T) {
	r := &http.Request{RemoteAddr: "198.51.100.3"}
	if got := ClientIP(r); got != "198.51.100.3" {
		t.Fatalf("got %q", got)
	}
}
