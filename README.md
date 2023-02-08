
# TritonHTTP

## Spec Summary

Here provide a concise summary of the TritonHTTP spec.

### HTTP Messages

TritonHTTP follows the [general HTTP message format](https://developer.mozilla.org/en-US/docs/Web/HTTP/Messages). And it has some further specifications:

- HTTP version supported: `HTTP/1.1`
- Request method supported: `GET`
- Response status supported:
  - `200 OK`
  - `400 Bad Request`
  - `404 Not Found`
- Request headers:
  - `Host` (required)
  - `Connection` (optional, `Connection: close` has special meaning influencing server logic)
  - Other headers are allowed, but won't have any effect on the server logic
- Response headers:
  - `Date` (required)
  - `Last-Modified` (required for a `200` response)
  - `Content-Type` (required for a `200` response)
  - `Content-Length` (required for a `200` response)
  - `Connection: close` (required in response for a `Connection: close` request, or for a `400` response)
  - Response headers should be written in sorted order for the ease of testing
  - Response headers should be returned in 'canonical form', meaning that the first letter and any letter following a hyphen should be upper-case. All other letters in the header string should be lower-case.

### Server Logic

When to send a `200` response?
- When a valid request is received, and the requested file can be found.

When to send a `404` response?
- When a valid request is received, and the requested file cannot be found or is not under the doc root.

When to send a `400` response?
- When an invalid request is received.
- When timeout occurs and a partial request is received.

When to close the connection?
- When timeout occurs and no partial request is received.
- When EOF occurs.
- After sending a `400` response.
- After handling a valid request with a `Connection: close` header.

When to update the timeout?
- When trying to read a new request.

What is the timeout value?
- 5 seconds.

## some features

- The client can issuemore than one TritonHTTP request without necessarily waiting for full HTTP replies to be returned (HTTP pipelining).
 <img width="530" alt="iShot_2023-02-07_19 34 57" src="https://user-images.githubusercontent.com/114261503/217423619-ff491f20-9efb-4d87-a6eb-aee82df246f8.png">

- In some cases, it is desirable to host multiple web servers on a single physical machine. This allows all the hosted web servers to share the physical server’s resources such as memory and processing, and in particular, to share a single IP address. This project implements virtual hosting by allowing TritonHTTP to host multiple servers. Each of these servers has a unique host name and maps to a unique docroot directory on the physical server. Every request sent to TritonHTTP includes the “Host” header, which is used to determine the web server that each request is destined for. 

## Usage

The source code for tools needed to interact with TritonHTTP can be found in `cmd`. The following commands can be used to launch these tools:

1) `make fetch` - A tool that allows you to construct custom responses and send them to your web server. Please refer to the README in `fetch`'s directory for more information.

2) `make gohttpd` - Starts up Go's inbuilt web-server.

3) `make tritonhttpd`  - Starts up implementation of TritonHTTP
