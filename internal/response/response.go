package response

import (
	"fmt"
	"io"
	"strconv"
	"myhttpfromtcp/internal/headers"
)

const CRLF = "\r\n"

type StatusCode int
const (
	STATUS_CODE_OK  StatusCode = iota
	STATUS_CODE_BAD_REQUEST
	STATUS_CODE_INTERNAL_ERROR
)
var statusCode = map[StatusCode]int{
	STATUS_CODE_OK:             200,
	STATUS_CODE_BAD_REQUEST:    400,
	STATUS_CODE_INTERNAL_ERROR: 500,
}

func WriteStatusLine(w io.Writer, code StatusCode) error {
	switch code {
	case STATUS_CODE_OK:
		w.Write([]byte("HTTP/1.1 200 OK"))

	case STATUS_CODE_BAD_REQUEST:
		w.Write([]byte("HTTP/1.1 400 Bad Request"))

	case STATUS_CODE_INTERNAL_ERROR:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error"))
	}

	w.Write([]byte(CRLF))
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	res := headers.NewHeaders()

	res.Set("Content-Length", strconv.Itoa(contentLen))
	res.Set("Connection", "close")
	res.Set("Content-Type", "text/plain")

	return res
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		line := fmt.Sprintf("%s: %s%s", key, val, CRLF)
		_, err := w.Write([]byte(line))

		if err != nil {
			return err
		}
	}

	w.Write([]byte(CRLF))
	return nil
}
