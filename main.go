package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mchmarny/aictl/pkg/chat"
	"github.com/mchmarny/aictl/pkg/chat/gemini"
)

var (
	chatter chat.Chat = &gemini.Chat{}
)

func main() {
	ctx := context.Background()
	defer chatter.Close(ctx)

	// flags
	if err := chatter.Init(ctx); err != nil {
		fmt.Printf("Error initializing chat: %s", err.Error())
		return
	}
	flag.Parse()

	// prompt
	if err := chatter.Start(ctx, bufio.NewScanner(os.Stdin)); err != nil {
		fmt.Printf("Error starting chat: %s", err.Error())
		return
	}
}
