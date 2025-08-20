package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

const clrf = "\r\n"

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(clrf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}
	headerLineText := string(data[:idx])
	key, value, err := kvPairFromString(headerLineText)
	if err != nil {
		return 0, false, err
	}
	h.Set(key, value)
	return idx + 2, false, nil
}

func kvPairFromString(headerString string) (key, value string, err error) {
	parts := strings.Fields(headerString)
	if len(parts) != 2 || !strings.HasSuffix(parts[0], ":") {
		return "", "", fmt.Errorf("error: invalid header string")
	}
	key = strings.Trim(parts[0], ":")
	if !validTokens([]byte(key)) {
		return "", "", fmt.Errorf("error: invalid header token found: %s", key)
	}
	value = parts[1]
	return key, value, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

// validTokens checks if the data contains only valid tokens
// or characters that are allowed in a token
func validTokens(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}
