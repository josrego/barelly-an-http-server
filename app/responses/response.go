package responses

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/encodings"
	"github.com/codecrafters-io/http-server-starter-go/app/requests"
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

func New(status int, contentType string, body []byte, req *requests.Request) *Response {
	headers := make(map[string]string)

	// Check for compression header
	encoder, err := getCompressionHeader(req)
	if err == nil {
		headers["Content-Encoding"] = encoder.EncodingType
		encodedBody, err := encoder.EncodingFun(&body)
		if err != nil {
			log.Fatalf("Error while compressing body: %s\n", err)
			return &Response{StatusCode: statusCodes[500]}
		} else {
			body = encodedBody
		}
	}

	return &Response{
		StatusCode:  statusCodes[status],
		ContentType: contentType,
		Headers:     headers,
		Body:        body,
	}
}

func getCompressionHeader(req *requests.Request) (*encodings.Encoder, error) {
	if encodingHeader, ok := req.Headers["accept-encoding"]; ok {
		// encoding can have multiple encoding values delimited by ','
		//  but just one correct one is accepted
		for _, encoding := range strings.Split(encodingHeader, ",") {
			// trim spaces first
			encoding = strings.TrimSpace(encoding)
			if encoder, found := encodings.AvailableEncoders[encoding]; found {
				fmt.Println("encoding found", encoding)
				return &encoder, nil
			}
		}
		fmt.Println("encoding not available in request: ", encodingHeader)
		return nil, errors.New("encoding not available in request")
	}

	return nil, errors.New("encoding not available")
}

func (resp *Response) OutputString() []byte {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%s %s%s", PROTOCOL_HTTP_1_1, resp.StatusCode.String(), CRLF))

	out.WriteString("Content-Type:" + resp.ContentType + CRLF)
	out.WriteString(fmt.Sprintf("Content-Length: %d%s", len(resp.Body), CRLF))

	for name, value := range resp.Headers {
		out.WriteString(fmt.Sprintf("%s: %s%s", name, value, CRLF))
	}

	out.WriteString(CRLF)
	out.Write(resp.Body)

	return out.Bytes()
}
