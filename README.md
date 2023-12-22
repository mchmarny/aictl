# aictl

WIP: Simple command line interface to AI models.

Currently supported models: 

* [gemini-pro](https://ai.google.dev/models/gemini#model_variations)


## Installation 

You can install `aictl` CLI using one of the following ways:

* [Go](#go)
* [Binary](#binary)

> See the [release section](https://github.com/mchmarny/aictl/releases/latest) for `aictl` checksums and SBOMs.

### Go

If you have Go 1.17 or newer, you can install latest `aictl` using:

```shell
go install github.com/mchmarny/aictl@latest
```

### Binary 

You can also download the [latest release](https://github.com/mchmarny/aictl/releases/latest) version of `aictl` for your operating system/architecture from [here](https://github.com/mchmarny/aictl/releases/latest). Put the binary somewhere in your $PATH, and make sure it has that executable bit.

> The official `aictl` releases include SBOMs

## API Token

Create API key: https://makersuite.google.com/app/apikey

Either export the `API_KEY` environment variable, or pass the key as a flag (see below).

```shell
export API_KEY="your-key-goes-here"
```

## Run




## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

