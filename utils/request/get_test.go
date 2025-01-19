package request

import "testing"

func Test_Get(t *testing.T) {
	t.Run("Testing Get", func(t *testing.T) {
		t.Run("Should return error if invalid URL", func(t *testing.T) {
			// Arrange
			req := NewGet("invalid-url", nil)

			_, err := req.Do()

			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
		t.Run("Should return response if valid URL", func(t *testing.T) {
			// Arrange
			req := NewGet("https://jsonplaceholder.typicode.com/posts/1", nil)
			_, err := req.Do()
			if err != nil {
				t.Errorf("Expected nil, got %v", err)
			}
		})

	})
}
