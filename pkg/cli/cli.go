package cli

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mchmarny/aictl/pkg/chat"
	"github.com/mchmarny/aictl/pkg/chat/gemini"
)

var (
	chatter chat.Chat = &gemini.Chat{}

	// set at build time
	version = "v0.0.1-default"
	commit  = "not-set"
	date    = "not-set"
)

func Start() {
	ctx := context.Background()
	defer chatter.Close(ctx)

	// flags
	info := flag.Bool("info", false, "Show version info.")
	if err := chatter.Init(ctx); err != nil {
		fmt.Printf("error initializing chat: %s", err.Error())
		return
	}
	flag.Parse()

	// info
	if *info {
		fmt.Printf("aictl (version: %s, commit: %s, built: %s)\n", version, commit, date)
		return
	}

	// interruptions (e.g. ctrl+c)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// prompt
	go func() {
		if err := chatter.Start(ctx, bufio.NewScanner(os.Stdin)); err != nil {
			fmt.Printf("error starting chat: %s", err.Error())
			return
		}
	}()

	<-done
	fmt.Println()
}
