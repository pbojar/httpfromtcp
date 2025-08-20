package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type status int

const (
	_ status = iota
	initialized
	done
)

type Request struct {
	RequestLine RequestLine

	state status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const clrf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		state: initialized,
	}
	buf := make([]byte, bufferSize)
	readToIndex := 0
	for req.state != done {
		// Double buffer length if full
		if readToIndex >= len(buf) {
			newBuf := make([]byte, 2*len(buf))
			copy(newBuf, buf)
			buf = newBuf
		}
		// Read to buffer
		nRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = done
				break
			}
			return nil, err
		}
		readToIndex += nRead
		// Parse portion of buffer read
		nParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[nParsed:])
		readToIndex -= nParsed
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something actually went wrong
			return 0, err
		}
		if n == 0 {
			// just need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = done
		return n, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(clrf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(requestString string) (*RequestLine, error) {
	parts := strings.Fields(requestString)
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
