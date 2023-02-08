package tritonhttp

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type Response struct {
	Proto      string // e.g. "HTTP/1.1"
	StatusCode int    // e.g. 200
	StatusText string // e.g. "OK"

	// Headers stores all headers to write to the response.
	Headers map[string]string

	// Request is the valid request that leads to this response.
	// It could be nil for responses not resulting from a valid request.
	// Hint: you might need this to handle the "Connection: Close" requirement
	Request *Request

	// FilePath is the local path to the file to serve.
	// It could be "", which means there is no file to serve.
	FilePath string
}

func (res *Response) init(req *Request) {
	res.Proto = "HTTP/1.1"
	res.Headers = make(map[string]string)
	res.Headers["Date"] = FormatTime(time.Now())
	res.Request = req
	if req != nil { //nil means a bad request
		if req.URL[len(req.URL)-1] == '/' {
			req.URL = req.URL + "index.html"
		}
		if req.Close {
			res.Headers["connection"] = "close"
		}
	}
}
func (res *Response) handleBadRequest() {
	res.init(nil)
	res.StatusCode = 400
	res.StatusText = "Bad Request"
	res.Request = nil
	res.Headers["connection"] = "close"
}

func (res *Response) handleNotFound(req *Request) {
	// res.init(nil)
	res.Request = req
	res.StatusCode = 404
	res.StatusText = "Not Found"
	//don't write any keyvalue pairs
}
func (res *Response) handleGood(req *Request, path string) {
	res.init(req)
	res.StatusCode = 200
	res.StatusText = "OK"
	res.FilePath = path
	stats, err := os.Stat(path)
	if err != nil {
		log.Println(err)
	}
	res.Headers["Last-Modified"] = FormatTime(stats.ModTime())
	res.Headers["Content-Type"] = MIMETypeByExtension(filepath.Ext(path))
	res.Headers["Content-Length"] = strconv.FormatInt(stats.Size(), 10)
}

func (res *Response) writeResponse(w io.Writer) error {
	err := res.writeStatus(w)
	if err != nil {
		return err
	}
	err = res.writeHeaders(w)
	if err != nil {
		return err
	}
	err = res.writeOptionalBody(w)
	return handleErr(err)
}

///////////////////////help function///////////////////////

func (res *Response) writeStatus(w io.Writer) error {
	var status string
	switch res.StatusCode {
	case 200:
		status = "200 OK"
	case 400:
		status = "400 Bad Request"
	case 404:
		status = "404 Not Found"
	}

	line := res.Proto + " " + status + "\r\n"
	_, err := w.Write([]byte(line))
	return handleErr(err)
}

func (res *Response) writeHeaders(w io.Writer) error {
	keys := make([]string, 0, len(res.Headers))
	for key, _ := range res.Headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		header := key + ": " + res.Headers[key] + "\r\n"
		_, err := w.Write([]byte(header))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return handleErr(err)
}

func (res *Response) writeOptionalBody(w io.Writer) error {
	var content []byte
	var err error
	if res.FilePath != "" {
		content, err = os.ReadFile(res.FilePath)
		if err != nil {
			return err
		}
	}
	_, err = w.Write(content)
	return handleErr(err)
}

func handleErr(err error) error {
	if err != nil {
		return err
	}
	return nil
}
