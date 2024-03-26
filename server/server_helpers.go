package server

import "strings"

func ParseResponse(rawResponse string) Response {
	// Split the response into parts by the double CRLF (which separates headers from body)
	parts := strings.SplitN(rawResponse, CRLF2, 2)
	headerPart := parts[0]
	body := ""
	if len(parts) > 1 {
		body = parts[1]
	}

	// Split the header part into lines
	lines := strings.Split(headerPart, CRLF)

	// The first line is the status line
	statusLine := lines[0]
	statusCode := strings.TrimSpace(strings.SplitN(statusLine, " ", 3)[1] + " " + strings.SplitN(statusLine, " ", 3)[2])

	// Parse headers (remaining lines)
	headers := make(map[string]string)
	for _, line := range lines[1:] {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}

	// Construct and return the Response struct
	return Response{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
	}
}
