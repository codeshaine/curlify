package request

import "testing"

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

	})
}
