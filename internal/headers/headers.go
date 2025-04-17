package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const CRLF = "\r\n"
const CHARS = "!#$%&'*+-.^_`|~"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(CRLF))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	key = strings.ToLower(strings.TrimSpace(key))

	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}

	value := strings.TrimSpace(string(parts[1]))
	h.Set(key, value)

	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	if string(value[len(value) - 1]) == ";" {
		value = value[:len(value)-1]
	}

	if _, ok := h[key]; ok {
		h[key] = fmt.Sprintf("%s, %s", h[key], value)
		return
	}

	h[key] = value
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

// validTokens checks if the data contains only valid tokens
// or characters that are allowed in a token
func validTokens(data []byte) bool {
	for _, c := range data {
		if c >= 'A' && c <= 'Z' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			continue
		}
		if c >= '0' && c <= '9' {
			continue
		}
		if bytes.ContainsAny(tokenChars, string(c)) {
			continue
		}
		return false
	}
	return true
}
