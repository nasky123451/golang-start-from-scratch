package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/m/server"
)

// TestHandler tests the HTTP server handler
func TestHandler(t *testing.T) {
	// Create a request to pass to the handler
	tests := []struct {
		name               string
		method             string
		body               string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "POST request",
			method:             http.MethodPost,
			body:               "Test POST body",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "POST request received and logged.\n",
		},
		{
			name:               "OPTIONS request",
			method:             http.MethodOptions,
			expectedStatusCode: http.StatusOK,
			expectedBody:       "",
		},
		{
			name:               "Unsupported method",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "Method not allowed.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request with the specified method and body
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))

			// Use httptest to create a ResponseRecorder which captures the response
			rr := httptest.NewRecorder()

			// Call the handler with our request and ResponseRecorder
			server.Handler(rr, req)

			// Check if the status code is what we expect
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rr.Code)
			}

			// Read the response body
			responseBody, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			// Check if the response body is what we expect
			if string(responseBody) != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, string(responseBody))
			}
		})
	}
}
