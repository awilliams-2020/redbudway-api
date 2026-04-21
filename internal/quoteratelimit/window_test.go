package quoteratelimit

import (
	"testing"
	"time"
)

func TestWindowLimiter_nilOrDisabled(t *testing.T) {
	var w *WindowLimiter
	if !w.Allow("k") {
		t.Fatal("nil limiter should allow")
	}
	w2 := NewWindowLimiter(0, time.Minute)
	if !w2.Allow("k") {
		t.Fatal("max<=0 should allow")
	}
	w3 := NewWindowLimiter(1, time.Minute)
	if !w3.Allow("") {
		t.Fatal("empty key should allow")
	}
}

func TestWindowLimiter_allowThenBlock(t *testing.T) {
	w := NewWindowLimiter(2, 50*time.Millisecond)
	if !w.Allow("a") {
		t.Fatal("first")
	}
	if !w.Allow("a") {
		t.Fatal("second")
	}
	if w.Allow("a") {
		t.Fatal("third should block")
	}
	time.Sleep(60 * time.Millisecond)
	if !w.Allow("a") {
		t.Fatal("after window should allow again")
	}
}

func TestWindowLimiter_independentKeys(t *testing.T) {
	w := NewWindowLimiter(1, time.Minute)
	if !w.Allow("x") {
		t.Fatal("x")
	}
	if !w.Allow("y") {
		t.Fatal("y different key")
	}
}

func TestEnvInt(t *testing.T) {
	t.Setenv("T_ENV_INT", "")
	if got := EnvInt("T_ENV_INT", 42); got != 42 {
		t.Fatalf("got %d", got)
	}
	t.Setenv("T_ENV_INT", "7")
	if got := EnvInt("T_ENV_INT", 42); got != 7 {
		t.Fatalf("got %d", got)
	}
	t.Setenv("T_ENV_INT", "nope")
	if got := EnvInt("T_ENV_INT", 3); got != 3 {
		t.Fatalf("got %d", got)
	}
	t.Setenv("T_ENV_INT", "-1")
	if got := EnvInt("T_ENV_INT", 9); got != 9 {
		t.Fatalf("negative invalid, got %d", got)
	}
}

func TestEnvDurationMinutes(t *testing.T) {
	t.Setenv("T_ENV_DUR", "")
	def := 10 * time.Minute
	if got := EnvDurationMinutes("T_ENV_DUR", def); got != def {
		t.Fatalf("got %v", got)
	}
	t.Setenv("T_ENV_DUR", "5")
	if got := EnvDurationMinutes("T_ENV_DUR", def); got != 5*time.Minute {
		t.Fatalf("got %v", got)
	}
	t.Setenv("T_ENV_DUR", "0")
	if got := EnvDurationMinutes("T_ENV_DUR", def); got != def {
		t.Fatalf("invalid should return default, got %v", got)
	}
}
