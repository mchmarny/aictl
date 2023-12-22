package gemini

import (
	"bufio"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	c := Chat{}
	t.Setenv(apiKeyEnvVar, "")

	t.Run("Init without API key", func(t *testing.T) {
		err := c.Init(context.TODO())
		assert.NoError(t, err)
	})

	t.Run("Start without API key", func(t *testing.T) {
		err := c.Start(context.TODO(), bufio.NewScanner(os.Stdin))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), apiKeyFlag)
	})

	t.Run("Init with API key", func(t *testing.T) {
		t.Setenv(apiKeyEnvVar, "test")
		err := c.Init(context.TODO())
		assert.NoError(t, err)
	})

	t.Run("Start with API Key", func(t *testing.T) {
		t.Setenv(apiKeyEnvVar, "test")
		err := c.Start(context.TODO(), bufio.NewScanner(os.Stdin))
		assert.NoError(t, err)
	})

	t.Run("Close", func(t *testing.T) {
		err := c.Close(context.TODO())
		assert.NoError(t, err)
	})
}

func TestGetFileContent(t *testing.T) {
	t.Run("File not found", func(t *testing.T) {
		_, err := getFileContent("test", "file-not-exists.txt")
		assert.Error(t, err)
	})

	t.Run("File found", func(t *testing.T) {
		_, err := getFileContent("test", "../../../content/annual-us-gdp.csv")
		assert.NoError(t, err)
	})
}
