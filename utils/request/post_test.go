package request

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func Test_Post(t *testing.T) {
	t.Run("Testing Post", func(t *testing.T) {
		t.Run("Should return error if invalid URL", func(t *testing.T) {
			// Arrange
			req := NewPost("invalid-url", nil, nil)
			_, err := req.Do()
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})

		t.Run("Should return response if valid URL", func(t *testing.T) {
			// Arrange
			req := NewPost("https://jsonplaceholder.typicode.com/posts", nil, nil)
			_, err := req.Do()
			if err != nil {
				t.Errorf("Expected nil, got %v", err)
			}
		})
		t.Run("Should send body successfully", func(t *testing.T) {
			// Prepare body
			bodyContent := []byte(`{"title": "foo", "body": "bar", "userId": 1}`)
			body := io.NopCloser(bytes.NewBuffer(bodyContent))

			// Set headers
			headers := http.Header{
				"Content-Type": {"application/json"},
			}

			// Arrange and Act
			req := NewPost("https://jsonplaceholder.typicode.com/posts", headers, body)
			resp, err := req.Do()

			// Assert
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("Expected status 201, got %d", resp.StatusCode)
			}
		})
	})
}
