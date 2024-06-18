package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateEndpointE2E(t *testing.T) {
	server := httptest.NewServer(buildRouter())
	defer server.Close()

	tests := []struct {
		description    string
		input          Input
		expectedOutput Output
		expectedStatus int
	}{
		{
			description:    "Valid input",
			input:          Input{A: 5, B: 3},
			expectedOutput: Output{A: 120, B: 6},
			expectedStatus: http.StatusOK,
		},
		{
			description:    "Negative input A",
			input:          Input{A: -1, B: 3},
			expectedOutput: Output{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			description:    "Negative input B",
			input:          Input{A: 5, B: -1},
			expectedOutput: Output{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			body, _ := json.Marshal(test.input)
			resp, err := http.Post(fmt.Sprintf("%s/calculate", server.URL), "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, resp.StatusCode)
			}

			if test.expectedStatus == http.StatusOK {
				var output Output
				body, _ := io.ReadAll(resp.Body)
				json.Unmarshal(body, &output)

				if output != test.expectedOutput {
					t.Errorf("Expected output %+v, got %+v", test.expectedOutput, output)
				}
			}
		})
	}
}
