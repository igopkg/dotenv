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
