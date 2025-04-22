package response

import (
	"fmt"
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

func GetStatusLine(code StatusCode) []byte {
	desc := ""

	switch code {
	case STATUS_CODE_OK:
		desc = "OK"

	case STATUS_CODE_BAD_REQUEST:
		desc = "Bad Request"

	case STATUS_CODE_INTERNAL_ERROR:
		desc = "Internal Server Error"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode[code], desc, CRLF))
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	res := headers.NewHeaders()

	res.Set("Content-Length", strconv.Itoa(contentLen))
	res.Set("Connection", "close")
	res.Set("Content-Type", "text/plain")

	return res
}
