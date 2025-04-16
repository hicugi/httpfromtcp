package request

import (
	"bytes"
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

const CRLF = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(CRLF))
	if idx == -1 {
		return nil, fmt.Errorf("Could not find CRLF in headline")
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}

	return requestLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	METHODS := map[string]bool{
		"GET": true,
		"POST": true,
		"PUT": true,
		"DELETE": true,
		"OPTIONS": true,
	}

	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Poorly formatted headline: %s", str)
	}

	method := parts[0]
	if _, ok := METHODS[method]; !ok {
		return nil, fmt.Errorf("Invalid method: %s", method)
	}

	reqPath := parts[1]
	if s := string(reqPath[0]); s != "/" {
		return nil, fmt.Errorf("Invalid path: %s", reqPath)
	}

	httpVer := parts[2]
	if httpVer != "HTTP/1.1" {
		return nil, fmt.Errorf("The app support only 'HTTP/1.1'. Invalid HTTP version: %s", httpVer)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: reqPath,
		HttpVersion:   "1.1",
	}, nil
}
