### README: Simple Go HTTP Server

#### Overview

This project demonstrates a basic HTTP server implemented in Go. It supports handling multiple connections concurrently, serving static files from a specified directory, and several specific HTTP routes with unique behaviors.

It is based on this challenge and extended to make it modular and with test coverage:
https://app.codecrafters.io/courses/http-server/completed

#### Features

- **Concurrency**: Utilizes Go's goroutines and TCP networking to handle multiple connections simultaneously.
- **Dynamic Content**: Supports routes for echoing URL segments and displaying the request's User-Agent.
- **File Handling**: Allows creating and retrieving files within a specified directory, using both `GET` and `POST` methods.

#### Supported Routes

- `GET /`: Returns a simple 200 OK response.
- `GET /echo/:message`: Echos the message back in the response body.
- `GET /user-agent`: Returns the request's User-Agent in the response body.
- `POST /files/:filename`: Creates a new file with the provided body content.
- `GET /files/:filename`: Retrieves the contents of a specified file.

#### How to Run

1. Compile and run the server, specifying the directory to serve:
   ```sh
   go run main.go -directory /path/to/dir
   ```
2. The server listens on port `4221` by default.

#### Testing

A separate test file (`server_test.go`) is included, demonstrating how to test the server's functionality. Use the Go testing framework to run these tests.

#### Dependencies

- Go standard library: No external dependencies are required.

#### Project Structure

- `server.go`: Contains the main server logic, including route handling and connection management.
- `server_test.go`: Provides tests for verifying the server's behavior.
- `server_helpers.go`: Includes helper functions for parsing HTTP responses.
