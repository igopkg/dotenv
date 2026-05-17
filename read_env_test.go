package dotenv

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// --- parser unit tests ---

func TestParseLine(t *testing.T) {
	cases := []struct {
		input   string
		wantKey string
		wantVal string
		wantOk  bool
		wantErr bool
	}{
		{"KEY=VALUE", "KEY", "VALUE", true, false},
		{" KEY = VALUE ", "KEY", "VALUE", true, false},
		{`KEY="hello world"`, "KEY", "hello world", true, false},
		{"KEY='hello world'", "KEY", "hello world", true, false},
		{"export KEY=VALUE", "KEY", "VALUE", true, false},
		{"# comment", "", "", false, false},
		{"", "", "", false, false},
		{"   ", "", "", false, false},
		{"KEY=", "KEY", "", true, false},
		{"NOEQUALS", "", "", false, true},
		{"=VALUE", "", "", false, true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			key, val, ok, err := parseLine(tc.input, 1, "test.env")
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok != tc.wantOk {
				t.Fatalf("ok: got %v, want %v", ok, tc.wantOk)
			}
			if !tc.wantOk {
				return
			}
			if key != tc.wantKey {
				t.Errorf("key: got %q, want %q", key, tc.wantKey)
			}
			if val != tc.wantVal {
				t.Errorf("val: got %q, want %q", val, tc.wantVal)
			}
		})
	}
}

func TestParseLines(t *testing.T) {
	lines := []string{
		"# This is a comment",
		"",
		"DB_HOST=localhost",
		"DB_PORT=5432",
		`DB_PASS="secret with spaces"`,
		"export APP_ENV=production",
	}
	got, err := parseLines(lines, "test.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"DB_PASS": "secret with spaces",
		"APP_ENV": "production",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("key %s: got %q, want %q", k, got[k], v)
		}
	}
	if len(got) != len(want) {
		t.Errorf("map length: got %d, want %d", len(got), len(want))
	}
}

func TestParseError_LineNumber(t *testing.T) {
	lines := []string{
		"GOOD=value",
		"BADLINE",
	}
	_, err := parseLines(lines, "test.env")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Line != 2 {
		t.Errorf("line: got %d, want 2", pe.Line)
	}
}

// --- isValidEnvFile tests ---

func TestIsValidEnvFile(t *testing.T) {
	valid := []string{".env", ".env.local", ".env.production", "config.env", ".env.test"}
	for _, name := range valid {
		if !isValidEnvFile(name) {
			t.Errorf("expected %q to be valid", name)
		}
	}
	invalid := []string{"envfile", "config.txt", "myfile"}
	for _, name := range invalid {
		if isValidEnvFile(name) {
			t.Errorf("expected %q to be invalid", name)
		}
	}
}

// --- integration tests ---

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad_SetsEnvVars(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "LOAD_TEST_KEY=loadval\n")
	t.Setenv("LOAD_TEST_KEY", "") // ensure clean state; t.Setenv restores on cleanup
	os.Unsetenv("LOAD_TEST_KEY")

	if err := Load(path); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("LOAD_TEST_KEY"); got != "loadval" {
		t.Errorf("got %q, want %q", got, "loadval")
	}
}

func TestLoad_NoOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "NO_OVERWRITE_KEY=fromfile\n")
	t.Setenv("NO_OVERWRITE_KEY", "original")

	if err := Load(path); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("NO_OVERWRITE_KEY"); got != "original" {
		t.Errorf("Load should not overwrite: got %q, want %q", got, "original")
	}
}

func TestOverload_Overwrites(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "OVERLOAD_KEY=newval\n")
	t.Setenv("OVERLOAD_KEY", "oldval")

	if err := Overload(path); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("OVERLOAD_KEY"); got != "newval" {
		t.Errorf("Overload should overwrite: got %q, want %q", got, "newval")
	}
}

func TestInvalidFilename(t *testing.T) {
	err := Load("config.txt")
	if !errors.Is(err, ErrFileName) {
		t.Errorf("expected ErrFileName, got %v", err)
	}
}

func TestFileNotFound(t *testing.T) {
	err := Load("/nonexistent/path/.env")
	if !errors.Is(err, ErrFileNotFound) {
		t.Errorf("expected ErrFileNotFound, got %v", err)
	}
}
