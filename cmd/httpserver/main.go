package main

import (
	"log"
	"fmt"
	"io"
	"strings"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"crypto/sha256"
	"myhttpfromtcp/internal/request"
	"myhttpfromtcp/internal/response"
	"myhttpfromtcp/internal/server"
	"myhttpfromtcp/internal/headers"
)

const PORT = 42069

func main() {
	server, err := server.Serve(PORT, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()
	log.Println("Server started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget

	httpbinPrefix := "/httpbin/"
	if strings.HasPrefix(target, httpbinPrefix) {
		val := strings.TrimPrefix(target, httpbinPrefix)
		handlerHttpBin(w, req, val)
		return
	}

	if target == "/yourproblem" {
		handler400(w, req)
		return
	}
	if target == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
	return
}

func handlerHttpBin(w *response.Writer, req *request.Request, urlPath string) {
	w.WriteStatusLine(response.STATUS_CODE_OK)

	url := fmt.Sprintf("https://httpbin.org/%s", urlPath)
	fmt.Println("Proxying to", url)

	res, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer res.Body.Close()

	w.WriteStatusLine(response.STATUS_CODE_OK)

	h := response.GetDefaultHeaders(0)
	h.SetOverride("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-SHA256")
	h.Set("Trailer", "X-Content-Length")
	h.Remove("Content-Length")
	w.WriteHeaders(h)

	fullBody := make([]byte, 0)

	w.WriteHeaders(h)

	const maxChunkSize = 1024
	buffer := make([]byte, maxChunkSize)

	for {
		n, err := res.Body.Read(buffer)
		fmt.Println("Read", n, "bytes")

		if n > 0 {
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}

			fullBody = append(fullBody, buffer[:n]...)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body done:", err)
	}

	trailers := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	trailers.SetOverride("X-Content-SHA256", sha256)
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))

	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Println("Error writing trailers:", err)
	}

	fmt.Println("Wrote trailers")
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.STATUS_CODE_BAD_REQUEST)

	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.SetOverride("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.STATUS_CODE_INTERNAL_ERROR)

	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.SetOverride("Content-Type", "text/html")

	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.STATUS_CODE_OK)

	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.SetOverride("Content-Type", "text/html")

	w.WriteHeaders(h)
	w.WriteBody(body)
}
