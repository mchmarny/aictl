package chat

import (
	"bufio"
	"context"
)

type Config struct {
	Description  string
	DefaultValue string
}

type Chat interface {
	Init(ctx context.Context) error
	Start(ctx context.Context, scanner *bufio.Scanner) error
	Close(ctx context.Context) error
}
