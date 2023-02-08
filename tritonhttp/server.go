package tritonhttp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	// Addr specifies the TCP address for the server to listen on,
	// in the form "host:port". It shall be passed to net.Listen()
	// during ListenAndServe().
	Addr string // e.g. ":0"

	// VirtualHosts contains a mapping from host name to the docRoot path
	// (i.e. the path to the directory to serve static files from) for
	// all virtual hosts that this server supports
	VirtualHosts map[string]string
}

// ListenAndServe listens on the TCP network address s.Addr and then
// handles requests on incoming connections.

func (s *Server) ListenAndServe() error {

	// Hint: Validate all docRoots
	for _, docRoot := range s.VirtualHosts {
		fil, err := os.Stat(docRoot)
		if os.IsNotExist(err) {
			return err
		}
		if !fil.IsDir() {
			return fmt.Errorf("doc root %q isn't a directory", docRoot)
		}
	}

	// Hint: create your listen socket and spawn off goroutines per incoming client
	li, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	log.Printf("server start listening on %q", li.Addr())
	defer li.Close()
	for {
		conn, err := li.Accept()
		if err != nil {
			return err
		}
		log.Printf("connect successfully")
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	br := bufio.NewReader(conn)
	for {
		err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Println("can't set timeout", err)
			conn.Close()
			return
		}
		//read request
		req, receivedSth, err := readRequest(br)
		log.Println("req,receivedsomething,err:", req, receivedSth, err)
		if errors.Is(err, io.EOF) {
			log.Println("connection closes :", err)
			conn.Close()
			return
		}
		if err, ok := err.(net.Error); ok && err.Timeout() {
			if !receivedSth {
				log.Println("timed out:", err)
				conn.Close()
				return
			}
			res := &Response{}
			res.handleBadRequest()
			_ = res.writeResponse(conn)
			_ = conn.Close()
			return
		}
		if err != nil {
			log.Println("bad request:", err)
			res := &Response{}
			res.handleBadRequest()
			res.writeResponse(conn)
			conn.Close()
			return
		}
		//nice request
		res := s.handleNiceRequest(req)
		err = res.writeResponse(conn)
		if err != nil {
			log.Println(err)
		}
		if req.Close {
			conn.Close()
			return
		}
	}

}
func (s *Server) handleNiceRequest(req *Request) (res *Response) {
	log.Println("it's a good request, now to determine wheather the file exists")
	res = &Response{}
	res.init(req)
	docRoot := s.VirtualHosts[req.Host]
	absPath := filepath.Join(docRoot, req.URL) //joins and cleans
	log.Println("absPath:", absPath)
	if absPath[:len(docRoot)] != docRoot {
		res.handleNotFound(req)
	} else if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
		res.handleNotFound(req)
	} else {
		res.handleGood(req, absPath)
	}

	return res
}
