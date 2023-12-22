# aictl

When interacting with AI models like [gemini-pro](https://ai.google.dev/models/gemini#model_variations) there are times when you need to add additional context to to support specific prompts based on data of which the model is not aware. This often includes coping and pasting content into the AI chatbot terminal (e.g. [bard](https://bard.google.com)). This terminal app allows for easier interaction with external data using prompts that add either local files (e.g. `+file:path`) or external resources (e.g. `+url:url`).

## Install 

You can install `aictl` CLI using one of the following ways:

* [Go](#go)
* [Binary](#binary)

> See the [release section](https://github.com/mchmarny/aictl/releases/latest) for `aictl` checksums and SBOMs.

### Binary 

You can also download the [latest release](https://github.com/mchmarny/aictl/releases/latest) version of `aictl` for your operating system/architecture from [here](https://github.com/mchmarny/aictl/releases/latest). Put the binary somewhere in your $PATH, and make sure it has that executable bit.

> The official `aictl` releases include SBOMs

### Go

If you have Go 1.17 or newer, you can install latest `aictl` using:

```shell
go install github.com/mchmarny/aictl@latest
```

## API Token

Create API key: https://makersuite.google.com/app/apikey

Either export the `API_KEY` environment variable, or pass the key as a flag (see below).

```shell
export API_KEY="your-key-goes-here"
```

## Run

To start the AI chat using default values:

```shell
aictl
```

If you haven't defined the `API_KEY` environment variable, you will have to pass it as a flag:

```shell
aictl --api-key your-key-goes-here
```

Additional parameters that can be passed as flags: 

* `temperature` (float `0.0` to `1.0`, default: `0.9`) the lower the number the more predictable the answers, higher numbers result in more creative responses.
* `tokens` (int `1` to `2048`, default: `100`) the maximum number of output tokens that will be returned from each prompt.

## Context

You can add your own context into the chat by inserting file content using `+file:` or remote content using `+url:` references. For example, at the chat prompt:

```shell
+file:content/monthly-gas-price.csv
```

The chat will ask you first for description of the file to understand its content:

```shell
chat: Describe content of content/annual-us-gdp.csv:
you: Annual US Gross Domestic Productivity
```

So then in chat you can combine that data with the content chat already knows: 

```shell
chat: How can I help?
you: What was the average gas price in US between 2010 and 2015?
chat: The average gas price in the US between 2010 and 2015 was $3.618 per gallon.
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

