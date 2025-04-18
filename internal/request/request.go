package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"errors"

	"myhttpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers

	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int
const (
	requestStateInitialized requestState = iota
	requestStateDone
	requestStateParsingHeaders
)

const CRLF = "\r\n"
const BUFFER_SIZE = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFFER_SIZE)
	readToIndex := 0

	req := &Request{
		Headers: headers.NewHeaders(),
		state: requestStateInitialized,
	}

	for req.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.state, numBytesRead)
				}
				break
			}

			return nil, err
		}

		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(CRLF))
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

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		totalBytesParsed += n
		if n == 0 {
			break
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
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
		r.state = requestStateParsingHeaders

		return n, nil

	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			r.state = requestStateDone
		}

		return n, nil

	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")

	default:
		return 0, fmt.Errorf("unknown state")
	}
}
