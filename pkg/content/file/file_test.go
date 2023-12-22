package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileContent(t *testing.T) {
	t.Run("File not found", func(t *testing.T) {
		_, err := GetContent("test", "file-not-exists.txt")
		assert.Error(t, err)
	})

	t.Run("File found", func(t *testing.T) {
		_, err := GetContent("test", "../../../content/annual-us-gdp.csv")
		assert.NoError(t, err)
	})
}
