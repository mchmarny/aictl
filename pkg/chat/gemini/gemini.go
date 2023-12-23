package gemini

import (
	"bufio"
	"context"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/google/generative-ai-go/genai"
	"github.com/mchmarny/aictl/pkg/content/file"
	"github.com/mchmarny/aictl/pkg/content/url"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	modelType = "gemini-pro"

	apiKeyEnvVar = "API_KEY"

	apiKeyFlag   = "api-key"
	tempFlag     = "temperature"
	maxTokenFlag = "tokens"
	topKFlag     = "top-k"
	topPFlag     = "top-p"

	filePrefix = "FILE:"
	urlPrefix  = "URL:"

	maxTokensDefault = 100 // 40-60 works (4 chars per token)
	tempDefault      = 0.2
	topKDefault      = 40
	topPDefault      = 0.95

	modelContentResponse = "Thank you for the context. What would you like to know?"
)

var (
	errStyle = color.New(color.FgRed, color.Bold)
	aiStyle  = color.New(color.FgGreen, color.Bold)
)

type Chat struct {
	client *genai.Client
	model  *genai.GenerativeModel

	apiKey      string
	temperature float32
	maxTokens   int32
	topK        int32
	topP        float32
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

	if flag.Lookup(topKFlag) == nil {
		flag.Func(topKFlag, "", func(flagValue string) error {
			for _, v := range strings.Fields(flagValue) {
				vv, err := strconv.ParseInt(v, 10, 32)
				if err != nil {
					return errors.Wrapf(err, "invalid configuration value for '%s'", topKFlag)
				}
				c.topK = int32(vv)
			}
			return nil
		})
	}

	if flag.Lookup(topPFlag) == nil {
		flag.Func(topPFlag, "", func(flagValue string) error {
			for _, v := range strings.Fields(flagValue) {
				vv, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return errors.Wrapf(err, "invalid configuration value for '%s'", topPFlag)
				}
				c.topP = float32(vv)
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

	if c.topK == 0 {
		c.topK = topKDefault
	}

	if c.topP == 0 {
		c.topP = topPDefault
	}

	return nil
}

func (c *Chat) Start(ctx context.Context, scanner *bufio.Scanner) error {
	// validation
	if err := c.validate(); err != nil {
		return err
	}

	if scanner == nil {
		return errors.New("missing scanner parameter")
	}

	// model
	if err := c.setup(ctx); err != nil {
		return err
	}

	// chat
	cs := c.model.StartChat()

	// send
	send := func(msg string) {
		// results
		iter := cs.SendMessageStream(ctx, genai.Text(msg))
		for {
			res, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				errStyle.Printf("error processing your prompt: %s\n", errors.Unwrap(err).Error())
				break
			}
			for _, c := range res.Candidates {
				if c.Content != nil {
					for _, p := range c.Content.Parts {
						aiStyle.Print(p)
					}
				}
			}
		}
		aiStyle.Println()
	}

	// load history
	load := func(msg string) {
		h := []*genai.Content{
			{
				Parts: []genai.Part{genai.Text(msg)},
				Role:  "user",
			},
			{
				Parts: []genai.Part{genai.Text(modelContentResponse)},
				Role:  "model",
			},
		}
		cs.History = append(cs.History, h...)
		aiStyle.Printf("%s:\n", modelContentResponse)
	}

	// files
	readFile := func(f string) error {
		aiStyle.Printf("Describe content of %s:\n", f)
		scanner.Scan()
		txt, err := file.GetContent(scanner.Text(), f)
		if err != nil {
			return errors.Wrapf(err, "error reading file: %s", f)
		}
		load(txt)
		return nil
	}

	// urls
	readURL := func(u string) error {
		aiStyle.Printf("Describe content of %s:\n", u)
		scanner.Scan()
		txt, err := url.GetContent(scanner.Text(), u)
		if err != nil {
			return errors.Wrapf(err, "error reading URL: %s", u)
		}
		load(txt)
		return nil
	}

	// prompt
	aiStyle.Println("How can I help?")
	for {
		scanner.Scan()
		text := scanner.Text()
		if len(text) == 0 {
			break
		}

		if strings.HasPrefix(text, filePrefix) {
			if err := readFile(text[len(filePrefix):]); err != nil {
				errStyle.Println(err.Error())
			}
			continue
		}

		if strings.HasPrefix(text, urlPrefix) {
			if err := readURL(text[len(urlPrefix):]); err != nil {
				errStyle.Println(err.Error())
			}
			continue
		}

		send(text)
		aiStyle.Println()
	}

	// error
	if scanner.Err() != nil {
		return errors.Wrapf(scanner.Err(), "error scanning input: %s", scanner.Err().Error())
	}

	return nil
}

func (c *Chat) setup(ctx context.Context) error {
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
	model.SetTopK(c.topK)
	model.SetTopP(c.topP)
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
	}

	c.model = model

	return nil
}
