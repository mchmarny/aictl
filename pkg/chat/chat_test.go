package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	c := &Config{
		Description:  "test1",
		DefaultValue: "test2",
	}

	assert.Equal(t, "test1 (test2)", c.String())

	if err := c.Validate(); err != nil {
		t.Errorf("config validation failed: %s", err)
	}
}
