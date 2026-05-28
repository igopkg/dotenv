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

func TestUnmarshal_Ints(t *testing.T) {
	t.Setenv("V_INT", "-1")
	t.Setenv("V_INT8", "-8")
	t.Setenv("V_INT16", "-16")
	t.Setenv("V_INT32", "-32")
	t.Setenv("V_INT64", "-64")

	var cfg struct {
		Int   int   `env:"V_INT"`
		Int8  int8  `env:"V_INT8"`
		Int16 int16 `env:"V_INT16"`
		Int32 int32 `env:"V_INT32"`
		Int64 int64 `env:"V_INT64"`
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Int != -1 || cfg.Int8 != -8 || cfg.Int16 != -16 || cfg.Int32 != -32 || cfg.Int64 != -64 {
		t.Errorf("unexpected values: %+v", cfg)
	}
}

func TestUnmarshal_Uints(t *testing.T) {
	t.Setenv("V_UINT", "1")
	t.Setenv("V_UINT8", "8")
	t.Setenv("V_UINT16", "16")
	t.Setenv("V_UINT32", "32")
	t.Setenv("V_UINT64", "64")
	t.Setenv("V_UINTPTR", "128")

	var cfg struct {
		Uint    uint    `env:"V_UINT"`
		Uint8   uint8   `env:"V_UINT8"`
		Uint16  uint16  `env:"V_UINT16"`
		Uint32  uint32  `env:"V_UINT32"`
		Uint64  uint64  `env:"V_UINT64"`
		Uintptr uintptr `env:"V_UINTPTR"`
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Uint != 1 || cfg.Uint8 != 8 || cfg.Uint16 != 16 || cfg.Uint32 != 32 || cfg.Uint64 != 64 || cfg.Uintptr != 128 {
		t.Errorf("unexpected values: %+v", cfg)
	}
}

func TestUnmarshal_Floats(t *testing.T) {
	t.Setenv("V_FLOAT32", "3.14")
	t.Setenv("V_FLOAT64", "2.718281828")

	var cfg struct {
		Float32 float32 `env:"V_FLOAT32"`
		Float64 float64 `env:"V_FLOAT64"`
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Float32 != 3.14 || cfg.Float64 != 2.718281828 {
		t.Errorf("unexpected values: %+v", cfg)
	}
}

func TestUnmarshal_Complex(t *testing.T) {
	t.Setenv("V_COMPLEX64", "1+2i")
	t.Setenv("V_COMPLEX128", "3+4i")

	var cfg struct {
		Complex64  complex64  `env:"V_COMPLEX64"`
		Complex128 complex128 `env:"V_COMPLEX128"`
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Complex64 != 1+2i || cfg.Complex128 != 3+4i {
		t.Errorf("unexpected values: %+v", cfg)
	}
}

func TestUnmarshal_InvalidInt(t *testing.T) {
	t.Setenv("V_INT", "notanint")

	var cfg struct {
		Int int `env:"V_INT"`
	}
	if err := Unmarshal(&cfg); err == nil {
		t.Error("expected error for invalid int, got nil")
	}
}

func TestUnmarshal_InvalidUint(t *testing.T) {
	t.Setenv("V_UINT", "-1")

	var cfg struct {
		Uint uint `env:"V_UINT"`
	}
	if err := Unmarshal(&cfg); err == nil {
		t.Error("expected error for negative uint, got nil")
	}
}

func TestUnmarshal_InvalidFloat(t *testing.T) {
	t.Setenv("V_FLOAT", "notafloat")

	var cfg struct {
		Float float64 `env:"V_FLOAT"`
	}
	if err := Unmarshal(&cfg); err == nil {
		t.Error("expected error for invalid float, got nil")
	}
}

func TestUnmarshal_InvalidComplex(t *testing.T) {
	t.Setenv("V_COMPLEX", "notacomplex")

	var cfg struct {
		Complex complex128 `env:"V_COMPLEX"`
	}
	if err := Unmarshal(&cfg); err == nil {
		t.Error("expected error for invalid complex, got nil")
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
