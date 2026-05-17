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
	value = strings.TrimSpace(after)
	if key == "" {
		return "", "", false, &ParseError{File: filename, Line: lineNum, Msg: "empty key"}
	}
	value = stripQuotes(value)
	return key, value, true, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == s[len(s)-1] && (s[0] == '"' || s[0] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}
