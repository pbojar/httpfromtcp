package response

import (
	"fmt"
	"io"

	"github.com/pbojar/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	_             StatusCode = iota
	OK                       = 200
	BadRequest               = 400
	InternalError            = 500
)

func getStatusLine(statusCode StatusCode) string {
	reasonMsg := ""
	switch statusCode {
	case OK:
		reasonMsg = "OK"
	case BadRequest:
		reasonMsg = "Bad Request"
	case InternalError:
		reasonMsg = "Internal Server Error"
	}
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonMsg)
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := fmt.Fprint(w, getStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := make(headers.Headers)
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", key, val)
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
