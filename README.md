# aictl

WIP: Simple command line interface to AI models.

Currently supported models: 

* [gemini-pro](https://ai.google.dev/models/gemini#model_variations)


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

* `temperature` (float 0.0 to 1.0, default: 0.9) the lower the number the more predictable the answers, higher numbers result in more creative responses.
* `tokens` (int 1 to 2048, default: 100) the maximum number of output tokens that will be returned from each prompt.

You can also add your own context into the chat by providing one or more data source files: 

```shell
aictl --file content/annual-us-gdp.csv --file content/monthly-gas-price.csv
```

For each one of the files you will be asked to provide a description to help chat understand the context:

```shell
chat: Describe content of content/annual-us-gdp.csv:
you: Annual US Gross Domestic Productivity
chat: Describe content of content/monthly-gas-price.csv:
you: Averaged annual gas prices in US
```

So then in chat you can combine that data with the content chat already knows: 

```shell
chat: How can I help?
you: What was the average gas price in US between 2010 and 2015?
chat: The average gas price in the US between 2010 and 2015 was **$3.618** per gallon.
Here is a breakdown of the average gas prices by year:
* 2010: $3.430
* 2011: $3.569
* 2012: $3.622
* 2013: $3.638
```

Similarly, you can add context to your prompts using remote content by passing one or more URLs:

> Note, that URL must be publicly accessible. 

```shell
aictl --url https://ai.google.dev/docs/safety_guidance
```


## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

