package responses

import (
	"bytes"
	"fmt"
)

type HttpStatusCode struct {
	Code   int
	Reason string
}

var statusCodes = map[int]HttpStatusCode{
	200: {Code: 200, Reason: "OK"},
	201: {Code: 201, Reason: "Created"},
	400: {Code: 400, Reason: "Bad Request"},
	404: {Code: 404, Reason: "Not Found"},
	500: {Code: 500, Reason: "Internal Server Error"},
}

const (
	JSON   = "application/json"
	TEXT   = "text/plain"
	BINARY = "application/octet-stream"

	CRLF              = "\r\n"
	PROTOCOL_HTTP_1_1 = "HTTP/1.1"
)

func (statusCode *HttpStatusCode) String() string {
	return fmt.Sprintf("%d %s", statusCode.Code, statusCode.Reason)
}

type Response struct {
	StatusCode  HttpStatusCode
	ContentType string
	Headers     map[string]string
	Body        []byte
}

func New(status int, contentType string, body []byte) *Response {
	return &Response{
		StatusCode:  statusCodes[status],
		ContentType: contentType,
		Headers:     make(map[string]string),
		Body:        body,
	}
}

func (resp *Response) OutputString() []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%s %s%s", PROTOCOL_HTTP_1_1, resp.StatusCode.String(), CRLF))

	// TODO: Headers here
	out.WriteString("Content-Type:" + resp.ContentType + CRLF)
	out.WriteString(fmt.Sprintf("Content-Length: %d%s", len(resp.Body), CRLF))

	out.WriteString(CRLF)
	out.Write(resp.Body)

	return out.Bytes()
}
