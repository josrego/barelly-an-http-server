package requests

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Method string

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"

	CRLF = "\r\n"
)

type Request struct {
	Method   Method
	Path     string
	Protocol string
	Headers  map[string]string
	Body     []byte
}

func (r *Request) LogString() string {
	var out bytes.Buffer

	method := fmt.Sprintf("%s", r.Method)
	out.WriteString("Method: " + method + "\t")
	out.WriteString("Path: " + r.Path + "\t")
	out.WriteString("Protocol: " + r.Protocol + "\t")

	return out.String()
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s %s", r.Method, r.Path, r.Protocol)
}

func New(requestData string) (*Request, error) {
	requestParts := strings.Split(requestData, CRLF)
	if len(requestParts) == 0 {
		log.Fatalln("request data has no lines")
		return nil, errors.New("request data has no lines")
	}

	iterator := 0
	// get verb, path and protocol from first line
	verb, path, protocol := processFirstLine(requestParts[iterator])

	// parse headers
	iterator++
	headers := make(map[string]string)
	for iterator < len(requestParts) {
		header := requestParts[iterator]
		if strings.Trim(header, "\r\n") == "" {
			fmt.Println("end of headers")
			iterator++
			break
		}

		hKey, hValue, found := strings.Cut(header, ": ")
		if !found {
			log.Fatalln("malformed header:", header)
			return nil, errors.New(fmt.Sprint("malformed header", header))
		}

		headers[strings.ToLower(hKey)] = hValue
		fmt.Printf("added header: %s: %s\n", hKey, hValue)
		iterator++
	}

	body := processStandardBody(requestParts[iterator:])

	return &Request{
		Method:   verb,
		Path:     path,
		Protocol: protocol,
		Headers:  headers,
		Body:     body,
	}, nil
}

func processFirstLine(firstLine string) (Method, string, string) {
	firsttLineParams := strings.Split(firstLine, " ")
	if len(firsttLineParams) > 3 {
		fmt.Println("Unexpected start line size. Got:", len(firsttLineParams),
			" - ", firsttLineParams)
	}

	verb := Method(firsttLineParams[0])
	path := firsttLineParams[1]
	protocol := firsttLineParams[2]

	return verb, path, protocol
}

func processStandardBody(requestLines []string) []byte {
	var bodyBuf bytes.Buffer

	for _, lines := range requestLines {
		bodyBuf.Write([]byte(lines))
	}
	return bodyBuf.Bytes()
}

func ParseMultiPartBody(requestLines []string) []byte {
	var bodyBuf bytes.Buffer

	return bodyBuf.Bytes()
}
