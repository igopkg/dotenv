package dotenv

import (
	"path/filepath"
	"strings"
)

func isValidEnvFile(name string) bool {
	return strings.Contains(filepath.Base(name), ".env")
}

func parseLines(lines []string, filename string) (map[string]string, error) {
	result := make(map[string]string)
	for i, raw := range lines {
		if i == 0 {
			raw = strings.TrimPrefix(raw, "\uFEFF") // UTF-8 BOM from some editors
		}
		key, value, ok, err := parseLine(raw, i+1, filename)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		result[key] = value
	}
	return result, nil
}

func parseLine(raw string, lineNum int, filename string) (key, value string, ok bool, err error) {
	line := strings.TrimSpace(raw)
	if line == "" || line[0] == '#' {
		return "", "", false, nil
	}
	if strings.HasPrefix(line, "export ") {
		line = strings.TrimSpace(line[len("export "):])
	}
	before, after, ok0 := strings.Cut(line, "=")
	if !ok0 {
		return "", "", false, &ParseError{File: filename, Line: lineNum, Msg: "missing '=' separator"}
	}
	key = strings.TrimSpace(before)
	if key == "" {
		return "", "", false, &ParseError{File: filename, Line: lineNum, Msg: "empty key"}
	}
	value, ok = parseValue(after)
	if !ok {
		return "", "", false, &ParseError{File: filename, Line: lineNum, Msg: "unterminated quoted value"}
	}
	return key, value, true, nil
}

// parseValue strips surrounding quotes and inline comments from a raw value.
// A '#' starts a comment only after a closing quote or when preceded by
// whitespace, so unquoted values like COLOR=#fff are kept intact.
// It reports false if the value opens a quote that is never closed.
func parseValue(raw string) (string, bool) {
	s := strings.TrimSpace(raw)
	if s != "" && (s[0] == '"' || s[0] == '\'') {
		end := strings.IndexByte(s[1:], s[0])
		if end < 0 {
			return "", false
		}
		rest := strings.TrimSpace(s[end+2:])
		if rest == "" || rest[0] == '#' {
			return s[1 : end+1], true
		}
	}
	for i := 1; i < len(raw); i++ {
		if raw[i] == '#' && (raw[i-1] == ' ' || raw[i-1] == '\t') {
			return strings.TrimSpace(raw[:i]), true
		}
	}
	return s, true
}
