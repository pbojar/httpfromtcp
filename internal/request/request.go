package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(r), "\r\n")
	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: *requestLine}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, fmt.Errorf("error: request line has %d parts, expected 3", len(parts))
	}
	httpVersionParts := strings.Split(parts[2], "/")
	if len(httpVersionParts) != 2 {
		return nil, fmt.Errorf("error: %s is an invalid HTTP-version", parts[2])
	}
	if httpVersionParts[0] != "HTTP" {
		return nil, fmt.Errorf("error: %s is an invalid HTTP-version", httpVersionParts[0])
	}
	requestLine := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpVersionParts[1],
	}
	// Validate method
	if !isUpper(requestLine.Method) {
		return nil, fmt.Errorf("error: %s is an invalid request method", requestLine.Method)
	}
	// Validate HTTP-version
	if requestLine.HttpVersion != "1.1" {
		return nil, fmt.Errorf("error: got %s as HTTP-version, expected 1.1", requestLine.HttpVersion)
	}
	return &requestLine, nil
}

func isUpper(s string) bool {
	for _, c := range s {
		if c > 'Z' || c < 'A' {
			return false
		}
	}
	return true
}
