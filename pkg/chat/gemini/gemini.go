package gemini

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	modelType = "gemini-pro"

	apiKeyEnvVar = "API_KEY"
	apiKeyFlag   = "apiKey"
	tempFlag     = "temperature"
	maxTokenFlag = "tokens"

	maxTokensDefault = 100
	tempDefault      = 0.9
)

type Chat struct {
	client *genai.Client

	apiKey      string
	temperature float32
	maxTokens   int32
}

func (c *Chat) validate() error {
	makeErr := func(c string) error {
		return errors.Errorf("chat configuration is invalid: %s not set", c)
	}

	if c.apiKey == "" {
		return makeErr(apiKeyFlag)
	}

	if c.temperature == 0 {
		return makeErr(tempFlag)
	}

	if c.maxTokens == 0 {
		return makeErr(maxTokenFlag)
	}

	return nil
}

func (c *Chat) Close(_ context.Context) error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *Chat) Init(_ context.Context) error {
	if flag.Lookup(apiKeyFlag) == nil {
		flag.Func(apiKeyFlag, "", func(flagValue string) error {
			for _, v := range strings.Fields(flagValue) {
				c.apiKey = v
			}
			return nil
		})
	}

	if flag.Lookup(tempFlag) == nil {
		flag.Func(tempFlag, "", func(flagValue string) error {
			for _, v := range strings.Fields(flagValue) {
				vv, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return errors.Wrapf(err, "invalid configuration value for '%s'", tempFlag)
				}
				c.temperature = float32(vv)
			}
			return nil
		})
	}

	if flag.Lookup(maxTokenFlag) == nil {
		flag.Func(maxTokenFlag, "", func(flagValue string) error {
			for _, v := range strings.Fields(flagValue) {
				vv, err := strconv.ParseInt(v, 10, 32)
				if err != nil {
					return errors.Wrapf(err, "invalid configuration value for '%s'", maxTokenFlag)
				}
				c.maxTokens = int32(vv)
			}
			return nil
		})
	}

	// defaults
	if c.apiKey == "" {
		c.apiKey = os.Getenv(apiKeyEnvVar)
	}

	if c.maxTokens == 0 {
		c.maxTokens = maxTokensDefault
	}

	if c.temperature == 0 {
		c.temperature = tempDefault
	}

	return c.validate()
}

func (c *Chat) Start(ctx context.Context, scanner *bufio.Scanner) error {
	// validation
	if err := c.validate(); err != nil {
		return err
	}

	if scanner == nil {
		return errors.New("missing scanner parameter")
	}

	// client
	client, err := genai.NewClient(ctx, option.WithAPIKey(c.apiKey))
	if err != nil {
		return errors.Wrapf(err, "error creating GenAI client: %s", err.Error())
	}
	c.client = client

	// model
	model := client.GenerativeModel(modelType)
	model.SetTemperature(c.temperature)
	model.SetMaxOutputTokens(c.maxTokens)

	// chat
	cs := model.StartChat()
	cs.History = []*genai.Content{}

	// send
	send := func(msg string) {
		// input
		fmt.Printf("=== Me: %s\n=== Model:\n", msg)

		// results
		iter := cs.SendMessageStream(ctx, genai.Text(msg))
		for {
			res, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				fmt.Printf("error processing your prompt: %s\n", err.Error())
				break
			}
			for _, c := range res.Candidates {
				if c.Content != nil {
					for _, p := range c.Content.Parts {
						fmt.Print(p)
					}
				}
			}
		}
		fmt.Println()
	}

	// prompt
	fmt.Println("How can I help?")
	for {
		scanner.Scan()
		text := scanner.Text()
		if len(text) == 0 {
			break
		}

		send(text)
		fmt.Println("\nAnything else?")
	}

	// error
	if scanner.Err() != nil {
		fmt.Println("error scanning input: ", scanner.Err())
	}

	return nil
}