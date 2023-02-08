package tritonhttp

import (
	"time"
)

// From a web page from hostname:port, returning the response as an array of bytes
// along with the duration for how long it took to retreive the page
func Fetch(hostname string, port string, inp []byte) ([]byte, time.Duration, error) {
	// connect to the server

	// start timing this http session

	// send the input to the server

	// read the input back from the server

	// stop timing this http session

	panic("todo")
}
