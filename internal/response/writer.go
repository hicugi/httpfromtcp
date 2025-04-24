package response

import (
	"fmt"
	"io"

	"myhttpfromtcp/internal/headers"
)

type WriterState int

const (
	WRITER_STATE_STATUS WriterState = iota
	WRITER_STATE_HEADERS
	WRITER_STATE_BODY
	WRITER_STATE_TRAILERS
)

type Writer struct {
	writerState WriterState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: WRITER_STATE_STATUS,
		writer:      w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != WRITER_STATE_STATUS {
		return fmt.Errorf("Cannot write status line in state %d", w.writerState)
	}

	defer func() { w.writerState = WRITER_STATE_HEADERS }()

	_, err := w.writer.Write(GetStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != WRITER_STATE_HEADERS {
		return fmt.Errorf("Cannot write headers in state %d", w.writerState)
	}

	defer func() { w.writerState = WRITER_STATE_BODY }()

	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s%s", k, v, CRLF)))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte(CRLF))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != WRITER_STATE_BODY {
		return 0, fmt.Errorf("Cannot write body in state %d", w.writerState)
	}

	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != WRITER_STATE_BODY {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != WRITER_STATE_BODY {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	endLine := fmt.Sprintf("0%s", CRLF)
	n, err := w.writer.Write([]byte(endLine))
	if err != nil {
		return n, err
	}
	w.writerState = WRITER_STATE_TRAILERS

	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != WRITER_STATE_TRAILERS {
		return fmt.Errorf("cannot write trailers in state %d", w.writerState)
	}

	defer func() { w.writerState = WRITER_STATE_BODY }()

	for k, v := range h {
		line := fmt.Sprintf("%s: %s%s", k, v, CRLF)
		_, err := w.writer.Write([]byte(line))

		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte(CRLF))
	return err
}
