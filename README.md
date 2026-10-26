# dotenv

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

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

# Quoted values preserve interior whitespace; surrounding quotes are stripped
DB_PASS="secret with spaces"
APP_NAME='my app'

# Inline comments start at a '#' preceded by whitespace
LOG_LEVEL=debug # trailing comment is dropped
GREETING="hello # not a comment" # but this is

# A '#' with no space before it is part of the value
COLOR=#fff
TAG=a#b

# export prefix is optional
export APP_ENV=production

# Empty values are valid
OPTIONAL_KEY=
```

A UTF-8 byte order mark at the start of a file is ignored, and Windows (`\r\n`) line endings are supported.

A line is rejected with a `ParseError` if it has no `=`, has an empty key, or opens a quote that is never closed (e.g. `KEY="unterminated` or `KEY="mismatched'`).

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
    DbPort  int           `env:"DB_PORT"`
    Timeout time.Duration `env:"TIMEOUT"`
    Debug   bool          `env:"DEBUG"`
}

cfg := Config{}
if err := dotenv.Unmarshal(&cfg); err != nil {
    log.Fatal(err)
}
```

Supported field types:

- `string`, `bool`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`
- `float32`, `float64`
- `complex64`, `complex128`
- `time.Duration` (parsed with `time.ParseDuration`, e.g. `"30s"`)

Integers accept base prefixes (`0x1F`, `0o17`, `0b101`) and are range-checked against the field's size. Any other field type returns an error. Nested structs are recursed into. Fields without an `env` tag are ignored, and an env var set to an empty string is treated as unset.

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

## License

MIT — see [LICENSE](LICENSE) for details.
