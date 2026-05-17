package dotenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Unmarshal reads environment variables into the exported fields of v, which
// must be a pointer to a struct. Fields are mapped via the `env` struct tag.
// Mark a field required with `req:"true"`; an error is returned if its env var
// is unset. Nested structs are traversed recursively. Supported field types:
// string, bool, and time.Duration.
func Unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dotenv: Unmarshal requires a pointer to a struct")
	}
	return populateFields(rv.Elem())
}

func populateFields(v reflect.Value) error {
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Struct {
			if err := populateFields(value); err != nil {
				return err
			}
			continue
		}

		tagName, ok := field.Tag.Lookup("env")
		if !ok || tagName == "" {
			continue
		}

		envVal := os.Getenv(tagName)
		if envVal == "" {
			if req, ok := field.Tag.Lookup("req"); ok {
				r, err := strconv.ParseBool(req)
				if err != nil {
					return fmt.Errorf("invalid req tag value %q for field %s", req, field.Name)
				}
				if r {
					return fmt.Errorf("required env var %s is not set", tagName)
				}
			}
			continue
		}

		switch value.Kind() {
		case reflect.String:
			value.SetString(envVal)
		case reflect.Int64:
			if field.Type == reflect.TypeFor[time.Duration]() {
				dur, err := time.ParseDuration(envVal)
				if err != nil {
					return fmt.Errorf("invalid duration for %s: %w", tagName, err)
				}
				value.SetInt(int64(dur))
			} else {
				return fmt.Errorf("unsupported int64 field %s", field.Name)
			}
		case reflect.Bool:
			parsed, err := strconv.ParseBool(envVal)
			if err != nil {
				return fmt.Errorf("invalid bool for %s: %w", tagName, err)
			}
			value.SetBool(parsed)
		default:
			return fmt.Errorf("unsupported config field type %s for %s", field.Type.String(), field.Name)
		}
	}
	return nil
}
