package routes

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/requests"
	"github.com/codecrafters-io/http-server-starter-go/app/responses"
)

type RouteHandler interface {
	HandleRequest(r *requests.Request) *responses.Response
}

type Route struct {
	Path    string // TODO: Path should be handled better (with path and query params)
	Handler RouteHandler
}

var routes map[string]*Route = make(map[string]*Route)

func New(path string, handler RouteHandler) *Route {
	return &Route{
		Path:    path,
		Handler: handler,
	}
}

func AddRoute(r *Route) {
	routes[r.Path] = r
}

func GetRouteHandler(requestPath string) RouteHandler {
	if requestPath == "" {
		fmt.Println("no request path given")
		return nil
	}

	for path, route := range routes {
		// TODO: this should be handled better.
		if strings.HasPrefix(requestPath, path) {
			fmt.Println("route found for ", requestPath, ": ", path)
			return route.Handler
		}
	}

	fmt.Println("no route found for ", requestPath)
	return nil
}

// Echo
type EchoHandler struct{}

func (echo *EchoHandler) HandleRequest(r *requests.Request) *responses.Response {
	message, found := strings.CutPrefix(r.Path, "/echo/")
	if !found {
		fmt.Println("no echo message found in: ", r.Path)
		return responses.New(200, responses.TEXT, nil)
	}

	fmt.Println("responding echo message: ", message)
	return responses.New(200, responses.TEXT, []byte(message))
}

// User Agent
type UserAgendHandler struct{}

func (userAgent *UserAgendHandler) HandleRequest(r *requests.Request) *responses.Response {
	if userAgent, ok := r.Headers["user-agent"]; ok {
		fmt.Printf("founder 'user-agent': %q\n", userAgent)
		return responses.New(200, responses.TEXT, []byte(userAgent))
	}

	fmt.Println("header 'User-Agent' not present in request.")
	return responses.New(404, responses.TEXT, []byte("header 'User-Agent' not present in request."))
}

// Files API
type FilesHandler struct {
	FilesPath string
}

func (filesHandler *FilesHandler) HandleRequest(r *requests.Request) *responses.Response {
	if r.Method == requests.GET {
		return filesHandler.HandleGetFileRequest(r)
	} else if r.Method == requests.POST {
		return filesHandler.HandlePostFileRequest(r)
	}

	return responses.New(400, responses.TEXT, []byte("Could not handle file request"))
}

func (filesHandler *FilesHandler) HandleGetFileRequest(r *requests.Request) *responses.Response {
	fileName, found := strings.CutPrefix(r.Path, "/files/")
	if !found {
		log.Fatalf("no fileName found in path: %s\n", r.Path)
		return responses.New(404, responses.TEXT, []byte(fmt.Sprintf("no fileName found in path: %s\n", r.Path)))
	}

	fileData, err := os.ReadFile(filesHandler.FilesPath + "/" + fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return responses.New(404, responses.TEXT, nil)
		}

		log.Fatalln("Error while reading file", err)
		return responses.New(500, responses.TEXT, []byte(""))
	}

	return responses.New(200, responses.BINARY, fileData)
}

func (filesHandler *FilesHandler) HandlePostFileRequest(r *requests.Request) *responses.Response {
	fileName, found := strings.CutPrefix(r.Path, "/files/")
	if !found {
		log.Fatalf("no fileName found in path: %s\n", r.Path)
		return responses.New(404, responses.TEXT, []byte(fmt.Sprintf("no fileName found in path: %s\n", r.Path)))
	}

	err := os.WriteFile(filesHandler.FilesPath+"/"+fileName, r.Body, 0644)
	if err != nil {
		log.Fatalln("Error while storing uploaded file:", err)
	}

	return responses.New(201, responses.TEXT, nil)
}
