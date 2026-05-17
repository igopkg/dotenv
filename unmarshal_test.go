package dotenv

import (
	"testing"
	"time"
)

func TestUnmarshal_String(t *testing.T) {
	t.Setenv("SERVER_ADDR", "localhost:8080")

	var cfg struct {
		Addr string `env:"SERVER_ADDR"`
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Addr != "localhost:8080" {
		t.Errorf("got %q, want %q", cfg.Addr, "localhost:8080")
	}
}

func TestUnmarshal_Bool(t *testing.T) {
	t.Setenv("DEBUG", "true")

	var cfg struct {
		Debug bool `env:"DEBUG"`
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if !cfg.Debug {
		t.Error("expected Debug to be true")
	}
}

func TestUnmarshal_Duration(t *testing.T) {
	t.Setenv("TIMEOUT", "5s")

	var cfg struct {
		Timeout time.Duration `env:"TIMEOUT"`
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("got %v, want %v", cfg.Timeout, 5*time.Second)
	}
}

func TestUnmarshal_NestedStruct(t *testing.T) {
	t.Setenv("DB_HOST", "db.local")
	t.Setenv("SERVER_ADDR", "0.0.0.0:9090")

	var cfg struct {
		Server struct {
			Addr string `env:"SERVER_ADDR"`
		}
		DB struct {
			Host string `env:"DB_HOST"`
		}
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Addr != "0.0.0.0:9090" {
		t.Errorf("Server.Addr: got %q, want %q", cfg.Server.Addr, "0.0.0.0:9090")
	}
	if cfg.DB.Host != "db.local" {
		t.Errorf("DB.Host: got %q, want %q", cfg.DB.Host, "db.local")
	}
}

func TestUnmarshal_RequiredMissing(t *testing.T) {
	var cfg struct {
		Addr string `env:"REQUIRED_ADDR" req:"1"`
	}
	if err := Unmarshal(&cfg); err == nil {
		t.Error("expected error for missing required var, got nil")
	}
}

func TestUnmarshal_RequiredPresent(t *testing.T) {
	t.Setenv("REQUIRED_ADDR", "127.0.0.1")

	var cfg struct {
		Addr string `env:"REQUIRED_ADDR" req:"1"`
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Addr != "127.0.0.1" {
		t.Errorf("got %q, want %q", cfg.Addr, "127.0.0.1")
	}
}

func TestUnmarshal_NoEnvTag_Skipped(t *testing.T) {
	var cfg struct {
		Ignored string
		Addr    string `env:"UNMARSHAL_ADDR"`
	}
	t.Setenv("UNMARSHAL_ADDR", "set")

	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Ignored != "" {
		t.Errorf("Ignored should be empty, got %q", cfg.Ignored)
	}
	if cfg.Addr != "set" {
		t.Errorf("Addr: got %q, want %q", cfg.Addr, "set")
	}
}

func TestUnmarshal_InvalidBool(t *testing.T) {
	t.Setenv("DEBUG", "notabool")

	var cfg struct {
		Debug bool `env:"DEBUG"`
	}
	if err := Unmarshal(&cfg); err == nil {
		t.Error("expected error for invalid bool, got nil")
	}
}

func TestUnmarshal_InvalidDuration(t *testing.T) {
	t.Setenv("TIMEOUT", "notaduration")

	var cfg struct {
		Timeout time.Duration `env:"TIMEOUT"`
	}
	if err := Unmarshal(&cfg); err == nil {
		t.Error("expected error for invalid duration, got nil")
	}
}

func TestUnmarshal_NotAPointer(t *testing.T) {
	type cfg struct{ Addr string }
	if err := Unmarshal(cfg{}); err == nil {
		t.Error("expected error for non-pointer, got nil")
	}
}

func TestUnmarshal_PointerToNonStruct(t *testing.T) {
	s := "hello"
	if err := Unmarshal(&s); err == nil {
		t.Error("expected error for pointer to non-struct, got nil")
	}
}
