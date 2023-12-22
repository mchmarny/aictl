package chat

import (
	"bufio"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type Config struct {
	Description  string
	DefaultValue string
}

func (c *Config) String() string {
	return fmt.Sprintf("%s (%s)", c.Description, c.DefaultValue)
}

func (c *Config) Validate() error {
	if c.Description == "" {
		return errors.Errorf("config description is empty")
	}

	if c.DefaultValue == "" {
		return errors.Errorf("config default value is empty")
	}

	return nil
}

type Chat interface {
	Init(ctx context.Context) error
	Start(ctx context.Context, scanner *bufio.Scanner) error
	Close(ctx context.Context) error
}
