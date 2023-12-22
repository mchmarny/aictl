package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURLContent(t *testing.T) {
	t.Run("URL not provided", func(t *testing.T) {
		_, err := GetContent("test", "")
		assert.Error(t, err)
	})
	t.Run("URL not found", func(t *testing.T) {
		_, err := GetContent("test", "http://bad-url-not-found.com")
		assert.Error(t, err)
	})

	t.Run("Valid URL", func(t *testing.T) {
		content, err := GetContent("test", "https://ai.google.dev/docs/safety_guidance")
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Contains(t, content, "Understanding the safety risks of your application")
	})
}
