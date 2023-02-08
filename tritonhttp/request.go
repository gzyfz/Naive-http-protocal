package tritonhttp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	Method string // e.g. "GET"
	URL    string // e.g. "/path/to/a/file"
	Proto  string // e.g. "HTTP/1.1"

	// Headers stores the key-value HTTP headers
	Headers map[string]string

	Host  string // determine from the "Host" header
	Close bool   // determine from the "Connection" header
}

func ReadLine(br *bufio.Reader) (string, error) {
	var line string
	for {
		s, err := br.ReadString('\n')
		line += s
		if err != nil {
			if err == io.EOF {
				return "", nil
			} else {
				return line, err
			}

		}
		if strings.HasSuffix(line, "\r\n") {
			line = line[:len(line)-2]
			return line, nil
		}
	}
}

func parseRequest(line string) (Method string, URL string, Proto string, err error) {
	afterParse := strings.SplitN(line, " ", 3)
	if len(afterParse) != 3 {
		return "", "", "", fmt.Errorf("400")
	} else {
		log.Println("method:", afterParse[0], "url:", afterParse[1], "protocal:", afterParse[2])
		return afterParse[0], afterParse[1], afterParse[2], nil
	}

}

// we need to know whether the server received sth before error occurs
func readRequest(br *bufio.Reader) (req *Request, receivedSth bool, err error) {
	req = &Request{}
	req.Headers = make(map[string]string)
	line, err := ReadLine(br)
	if err != nil {
		return nil, false, err
	}
	req.Method, req.URL, req.Proto, err = parseRequest(line)
	if err != nil {
		return nil, true, err
	}
	if req.Method != "GET" || req.URL[0] != '/' {
		//bad request
		return nil, true, fmt.Errorf("400")
	}

	hostAppear := false
	req.Close = false
	for {
		line, err := ReadLine(br)
		if err != nil {
			return nil, true, err
		}
		if line == "" {
			//读完了
			break
		}
		pairs := strings.SplitN(line, ": ", 2)
		if len(pairs) != 2 {
			return nil, true, fmt.Errorf("400") //request有问题啊大哥
		}
		key := CanonicalHeaderKey(strings.TrimSpace(pairs[0]))
		value := strings.TrimSpace(pairs[1])

		if key == "Host" {
			req.Host = value
			hostAppear = true
		} else if key == "Connection" && value == "close" {
			req.Close = true
		} else {
			//nothing special
			req.Headers[key] = value
		}

	}
	if !hostAppear {
		return nil, true, fmt.Errorf("400")
	}
	//everything alright

	return req, true, nil

}
