# dotenv

A Go package for loading `.env` files into the process environment.

## Installation

```bash
go get github.com/igopkg/dotenv
```

## Usage

```go
import "github.com/igopkg/dotenv"

// Load ".env" (default) — does not overwrite existing env vars
err := dotenv.Load()

// Load specific files
err = dotenv.Load(".env", ".env.local")

// Overwrite existing env vars
err = dotenv.Overload(".env.production")

// Parse without modifying the environment
vars, err := dotenv.Read(".env")
fmt.Println(vars["DB_HOST"])
```

## .env file format

```dotenv
# Comments are ignored
DB_HOST=localhost
DB_PORT=5432

# Quoted values preserve interior whitespace
DB_PASS="secret with spaces"
APP_NAME='my app'

# export prefix is optional
export APP_ENV=production

# Empty values are valid
OPTIONAL_KEY=
```

## API

| Function | Description |
|---|---|
| `Load(files ...string) error` | Parse files and set env vars. Skips keys already set. Defaults to `.env`. |
| `Overload(files ...string) error` | Like `Load` but overwrites existing env vars. |
| `Read(files ...string) (map[string]string, error)` | Parse files and return a map. Does not modify the environment. |
| `Unmarshal(v any) error` | Populate a struct from environment variables using `env` struct tags. |

### Unmarshal

`Unmarshal` reads environment variables into a struct. Use the `env` tag to map fields, and `req:"1"` to mark required variables.

```go
type Config struct {
    Addr    string        `env:"SERVER_ADDR" req:"1"`
    DbHost  string        `env:"DB_HOST"     req:"1"`
    Timeout time.Duration `env:"TIMEOUT"`
    Debug   bool          `env:"DEBUG"`
}

cfg := Config{}
if err := dotenv.Unmarshal(&cfg); err != nil {
    log.Fatal(err)
}
```

Supported field types: `string`, `bool`, `time.Duration`. Nested structs are recursed into. Fields without an `env` tag are ignored.

The `req` tag accepts any value parseable by `strconv.ParseBool` (`"1"`, `"true"`, `"0"`, `"false"`, etc.). An invalid `req` value returns an error at startup.

## Filename rules

Files must have `.env` somewhere in the base name: `.env`, `.env.local`, `.env.production`, `config.env`, etc.

## Error handling

```go
var pe *dotenv.ParseError
if errors.As(err, &pe) {
    fmt.Printf("parse failed at %s line %d: %s\n", pe.File, pe.Line, pe.Msg)
}

if errors.Is(err, dotenv.ErrFileName) { /* invalid filename */ }
if errors.Is(err, dotenv.ErrFileNotFound) { /* file missing */ }
```
