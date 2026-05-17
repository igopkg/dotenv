package dotenv

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
)

func Load(files ...string) error {
	if len(files) == 0 {
		files = []string{".env"}
	}
	vars, err := readFiles(files)
	if err != nil {
		return err
	}
	for k, v := range vars {
		if _, exists := os.LookupEnv(k); !exists {
			if err := os.Setenv(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func Overload(files ...string) error {
	if len(files) == 0 {
		files = []string{".env"}
	}
	vars, err := readFiles(files)
	if err != nil {
		return err
	}
	for k, v := range vars {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

func readFiles(files []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, path := range files {
		if !isValidEnvFile(path) {
			return nil, ErrFileName
		}
		f, err := os.Open(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
			}
			return nil, err
		}
		scanner := bufio.NewScanner(f)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		f.Close()
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		parsed, err := parseLines(lines, path)
		if err != nil {
			return nil, err
		}
		maps.Copy(result, parsed)
	}
	return result, nil
}
