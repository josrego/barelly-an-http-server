package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/requests"
	"github.com/codecrafters-io/http-server-starter-go/app/responses"
	"github.com/codecrafters-io/http-server-starter-go/app/routes"
)

const (
	OK_RESPONSE        = "HTTP/1.1 200 OK\r\n\r\n"
	NOT_FOUND_RESPONSE = "HTTP/1.1 404 NOT_FOUND\r\n\r\n"

	// settings
	BUFFER_SIZE = 4096
)

var currentRequests []requests.Request = make([]requests.Request, 0)

// handles new connection received
func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024) // Create a buffer to hold data

	n, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading:", err)
			return
		}
	}

	data := string(buffer[:n])
	req, err := requests.New(data)

	if err != nil {
		log.Fatalln("Error while creating request from data: ", data)
		return
	}

	currentRequests = append(currentRequests, *req)
	fmt.Println("Parsed new request. ", req.LogString())

	resp := getResp(req)

	output := resp.OutputString()
	conn.Write(output)
	fmt.Printf("Response: %q\n", output)
}

func getResp(req *requests.Request) *responses.Response {
	if req.Path == "/" {
		return responses.New(200, responses.TEXT, nil, req)
	} else {
		handler := routes.GetRouteHandler(req.Path)
		if handler == nil {
			return responses.New(404, responses.TEXT, nil, req)
		}

		return handler.HandleRequest(req)
	}
}

// make HTTP server initializations
func initialize() {
	currentRequests = make([]requests.Request, 0)

	// create routes
	directory := flag.String("directory", "", "directory for files API")
	flag.Parse()

	routes.AddRoute(routes.New("/echo", &routes.EchoHandler{}))
	routes.AddRoute(routes.New("/user-agent", &routes.UserAgendHandler{}))

	if *directory != "" {
		fmt.Println("Found directory argument. Setting up Files API on directory: ", *directory)
		routes.AddRoute(routes.New("/files", &routes.FilesHandler{FilesPath: *directory}))
	}
}

func main() {
	initialize()
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// initiates infinite loop to accept connections after port is binded
	for {
		con, err := l.Accept()
		if err == nil {
			go handleConnection(con)
			continue
		}

		fmt.Println("Error receiving connection: ", err.Error())
	}
}
