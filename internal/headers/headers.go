package headers

import (
	"bytes"
	"fmt"
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
		return 0, true, nil
	}
	headerLineText := string(data[:idx])
	key, value, err := kvPairFromString(headerLineText)
	if err != nil {
		return 0, false, err
	}
	h[key] = value
	return idx + 2, false, nil
}

func kvPairFromString(headerString string) (key, value string, err error) {
	parts := strings.Fields(headerString)
	if len(parts) != 2 || !strings.HasSuffix(parts[0], ":") {
		return "", "", fmt.Errorf("error: invalid header string")
	}
	key = strings.Trim(parts[0], ":")
	value = parts[1]
	return key, value, nil
}
