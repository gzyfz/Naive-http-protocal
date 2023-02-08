This is a simple tool to test server's implementation. It constructs a HTTP request by reading from the input file and sends this request to the server listening at  `<host_name>:<port_number>`.  If the request is valid and the server finds the object being requested, it sends back a response, whose contents are stored in the output file.

Usage:
`go run main.go -req <path_to_input_file> -resp <path_to_output_file> <host_name>:<port_number>`

