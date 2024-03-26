package server

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"
	"time"
)

// setupServerAndTest starts the server and returns a function for cleanup.
func setupServerAndTest(t *testing.T, port string) (func(), string) {

	t.Log("Setting up server for test")

	// Create a temporary directory to simulate the server's file storage
	tempDir, err := os.MkdirTemp("", "test_server_files")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Start the server in a separate goroutine.
	go func() {
		_ = StartServer(port, tempDir)
	}()

	// Wait to ensure the server is running.
	time.Sleep(time.Second)

	// Return a cleanup function.
	teardown := func() {
		t.Log("Tearing down after test")
		// Cleanup code: remove the temporary directory
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	}
	return teardown, tempDir
}

func makeRequest(port string, request string) (Response, error) {
	zeroResponse := Response{}

	// Connect to the server
	address := fmt.Sprintf("localhost:%s", port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return zeroResponse, fmt.Errorf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	// Send the request.
	_, err = fmt.Fprint(conn, request)

	if err != nil {
		return zeroResponse, fmt.Errorf("Failed to send request: %v", err)
	}

	// Read the response into a buffer.
	responseBuffer := make([]byte, 1024)
	bytesRead, err := conn.Read(responseBuffer)
	if err != nil {
		return zeroResponse, fmt.Errorf("Could not read response: %v", err)
	}
	rawResponse := string(responseBuffer[:bytesRead])

	actualResponse := ParseResponse(rawResponse)
	return actualResponse, nil
}

func TestHTTPServer(t *testing.T) {
	port := "8080"
	tests := []struct {
		name     string
		request  string
		response Response
	}{
		{
			name:    "Test / 200 OK",
			request: "GET / HTTP/1.1\r\n\r\n",
			response: Response{
				StatusCode: StatusOK,
				Headers:    map[string]string{"Content-Type": ContentTypeText, "Content-Length": "0"},
				Body:       "",
			},
		},
		{
			name:    "Test /not-found 404 Not Found",
			request: "GET /not-found HTTP/1.1\r\n\r\n",
			response: Response{
				StatusCode: StatusNotFound,
				Headers:    map[string]string{"Content-Type": ContentTypeText, "Content-Length": "0"},
				Body:       "",
			},
		},
		{
			name:    "Test /echo/abc 200 OK",
			request: "GET /echo/abc HTTP/1.1\r\n\r\n",
			response: Response{
				StatusCode: StatusOK,
				Headers:    map[string]string{"Content-Type": ContentTypeText, "Content-Length": "3"},
				Body:       "abc",
			},
		},
		{
			name:    "Test /user-agent/abc 200 OK",
			request: "GET /user-agent HTTP/1.1\r\nUser-Agent: curl/7.6\r\n\r\n",
			response: Response{
				StatusCode: StatusOK,
				Headers:    map[string]string{"Content-Type": ContentTypeText, "Content-Length": "8"},
				Body:       "curl/7.6",
			},
		},
		{
			name:    "Test /files/test.txt 201 OK",
			request: "POST /files/test.txt HTTP/1.1\r\nUser-Agent: curl/7.6\r\n\r\nHello",
			response: Response{
				StatusCode: StatusCreated,
				Headers:    map[string]string{"Content-Type": ContentTypeText, "Content-Length": "0"},
				Body:       "",
			},
		},
		{
			name:    "Test /files/test.txt 200 OK",
			request: "GET /files/test.txt HTTP/1.1\r\n\r\n",
			response: Response{
				StatusCode: StatusOK,
				Headers:    map[string]string{"Content-Type": ContentTypeOctetStream, "Content-Length": "5"},
				Body:       "Hello",
			},
		},
	}

	teardown, _ := setupServerAndTest(t, port)
	defer teardown()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualResponse, err := makeRequest(port, tt.request) // replace with your actual request function
			if err != nil {
				t.Fatalf("Failed to make GET request: %v", err)
			}

			if !reflect.DeepEqual(actualResponse, tt.response) {
				t.Errorf("Expected response:\n%v\ngot:\n%v", tt.response, actualResponse)
			}
		})
	}
}
