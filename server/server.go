package server

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	StatusOK                  = "200 OK"
	StatusCreated             = "201 Created"
	StatusNotFound            = "404 Not Found"
	StatusInternalServerError = "500 Internal Server Error"
	ContentTypeText           = "text/plain"
	ContentTypeOctetStream    = "application/octet-stream"
	CRLF                      = "\r\n"
	CRLF2                     = "\r\n\r\n"
)

type Request struct {
	Method  string
	Path    string
	Body    string
	Headers map[string]string
}

type Response struct {
	StatusCode string
	Headers    map[string]string
	Body       string
}

func StartServer(port string, directory string) error {
	address := ":" + port // This will listen on the given port on all interfaces

	// bind to port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Failed to bind to port: ", port)
		return err
	}

	// accept multiple connections
	defer listener.Close()
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			conn.Close()
		}
		// magic of a goroutine
		// allows to concurrently handle multiple connections
		go handleConnection(conn, directory)
	}
}

func handleConnection(conn net.Conn, directory string) {
	// defer ensures connection is closed when the function returns
	defer conn.Close()

	// read data from the connection
	buf := make([]byte, 1024)
	requestSize, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return
	}
	fmt.Println("Read ", requestSize, " bytes from connection")

	// convert the request to a string
	requestString := string(buf[:requestSize])

	// handle the request and write the response
	request := createRequest(requestString)
	response := handleRequest(request, directory)
	writeResponse(conn, response)
}

func handleRequest(request Request, directory string) Response {
	switch {
	case request.Method == "GET" && strings.HasPrefix(request.Path, "/echo"):
		return handleEcho(request.Path)
	case request.Method == "GET" && strings.HasPrefix(request.Path, "/user-agent"):
		return handleUserAgent(request.Headers)
	case request.Method == "GET" && strings.HasPrefix(request.Path, "/files"):
		return handleGetFiles(request.Path, directory)
	case request.Method == "POST" && strings.HasPrefix(request.Path, "/files"):
		return handlePostFiles(request, directory)
	case request.Path == "/":
		return Response{StatusCode: StatusOK}
	default:
		return Response{StatusCode: StatusNotFound}
	}
}

func handleEcho(path string) Response {
	nextSegment := strings.TrimPrefix(path, "/echo/") // /echo/abc -> abc
	return Response{StatusCode: StatusOK, Body: nextSegment}
}

func handleUserAgent(headers map[string]string) Response {
	userAgent := headers["User-Agent"]
	return Response{StatusCode: StatusOK, Body: userAgent}
}

func handleGetFiles(path string, directory string) Response {
	filename := strings.TrimPrefix(path, "/files/") // /files/abc -> abc
	fileContents, err := os.ReadFile(directory + "/" + filename)
	if err != nil {
		return handleError(err, "Error reading file: ", StatusNotFound)
	} else {
		return Response{
			StatusCode: StatusOK,
			Headers:    map[string]string{"Content-Type": "application/octet-stream"},
			Body:       string(fileContents)}
	}
}

func handlePostFiles(request Request, directory string) Response {
	filename := strings.TrimPrefix(request.Path, "/files/") // /files/abc -> abc
	err := os.WriteFile(directory+"/"+filename, []byte(request.Body), 0644)
	if err != nil {
		return handleError(err, "Error writing file: ", StatusInternalServerError)
	} else {
		return Response{StatusCode: StatusCreated}
	}
}

func handleError(err error, message string, statusCode string) Response {
	fmt.Println(message, err)
	return Response{StatusCode: statusCode}
}

func createRequest(requestString string) Request {
	lines := strings.Split(requestString, CRLF)
	requestLine := strings.Split(lines[0], " ") // GET /user-agent HTTP/1.1
	method := requestLine[0]
	path := requestLine[1]

	headers := make(map[string]string)
	for _, line := range lines[1:] {
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ": ", 2)
		headers[parts[0]] = parts[1]
	}

	bodyStart := strings.Index(requestString, CRLF2) + len(CRLF2)
	body := requestString[bodyStart:]

	return Request{
		Method:  method,
		Path:    path,
		Body:    body,
		Headers: headers,
	}
}

func writeResponse(conn net.Conn, response Response) {
	// Ensure Headers is initialized
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	// Default Content-Type if not set
	if _, exists := response.Headers["Content-Type"]; !exists {
		response.Headers["Content-Type"] = "text/plain"
	}
	// Calculate Content-Length and set it
	response.Headers["Content-Length"] = fmt.Sprintf("%d", len(response.Body))

	responseString := fmt.Sprintf("HTTP/1.1 %s%s", response.StatusCode, CRLF)
	for header, value := range response.Headers {
		responseString += fmt.Sprintf("%s: %s%s", header, value, CRLF)
	}
	responseString += CRLF
	responseString += response.Body

	fmt.Println("Writing response: ", responseString)
	conn.Write([]byte(responseString))
}
