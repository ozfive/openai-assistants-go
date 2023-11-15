# openai-assistants-go

[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=ozfive_openai-assistants-go&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=ozfive_openai-assistants-go) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ozfive_openai-assistants-go&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=ozfive_openai-assistants-go) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ozfive_openai-assistants-go&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ozfive_openai-assistants-go) [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ozfive_openai-assistants-go&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=ozfive_openai-assistants-go) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ozfive_openai-assistants-go&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ozfive_openai-assistants-go) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ozfive_openai-assistants-go&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ozfive_openai-assistants-go) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=ozfive_openai-assistants-go&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=ozfive_openai-assistants-go)

## NOTE: THIS REPOSITORY IS STILL IN DEVELOPMENT AND WILL NOT BE PRODUCTION READY UNTIL TESTS ARE INCORPORATED

openai-assistants-go is a Go package providing a convenient and robust interface for interacting with the OpenAI Assistants API. Simplify the integration of OpenAI's powerful language models into your Go applications with this well-structured and easy-to-use package.

## Installation

```shell
go get github.com/ozfive/openai-assistants-go
```

## Getting Started

To get started with the OpenAI Assistants Go Library, you need to initialize the client by providing your OpenAI API key. Follow these steps:

### Import the library in your Go code

```go
import (
    "github.com/ozfive/openai-assistants-go"
)
```

### Initialize the client with your OpenAI API key

```go
client := assistants.NewClient("YOUR_OPENAI_API_KEY")
```