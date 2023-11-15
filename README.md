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

### Step 1: Create an Assistant

Create an Assistant by defining its custom instructions and choosing a model. Enable tools like Code Interpreter, Retrieval, and Function calling if needed.

```go
createAssistantRequest := &assistants.CreateAssistantRequest{
    Instructions: "Your custom instructions here",
    Model:        "gpt-4-1106-preview",
    Tools: []map[string]string{
        {"type": "retrieval"},
    },
}

assistant, err := assistants.CreateAssistant(ctx, createAssistantRequest)
if err != nil {
    log.Fatal(err)
}
```

### Step 2: Create a Thread

Create a Thread when a user initiates a conversation. Pass any user-specific context and files in this thread by creating Messages.

```go
threadParams := assistants.CreateThreadParams{
    Messages: []assistants.MessageObject{
        // You can include messages in the thread if needed
        // Example: {"role": "system", "content": "System message"},
        //          {"role": "user", "content": "User message"},
    },
    Metadata: map[string]string{
        // Include any metadata relevant to the thread
        "user_id": "123",
        "context": "initial",
    },
}

thread, err := assistants.CreateThread(ctx, threadParams)
if err != nil {
    log.Fatal(err)
}
```

### Step 3: Add a Message to a Thread

Add a Message to the Thread containing the user's text. Optionally, include any files that the user uploads.

```go
messageParams := assistants.CreateMessageParams{
    Role:    "user",
    Content: "I need to solve the equation 3x + 11 = 14. Can you help me?",
    // Include file IDs if the user uploads any files
    FileIDs: []string{"file_id_1", "file_id_2"},
    // Include any metadata relevant to the message
    Metadata: map[string]string{
        "user_id": "123",
        "context": "equation_help",
    },
}

message, err := assistants.CreateMessage(ctx, thread.ID, messageParams)
if err != nil {
    log.Fatal(err)
}
```

### Step 4: Run the Assistant

Create a Run to make the Assistant read the Thread and decide whether to call tools or simply use the model to answer the user query.

```go
runParams := assistants.CreateRunParams{
    AssistantID:  assistant.ID,
    Model:        "gpt-4-1106-preview",  // Optional: Specify a different model
    Instructions: "Read and respond to the user's query",
    // Include tools if needed
    Tools: []assistants.ToolObject{
        {
            Type: "code",
            Function: &assistants.FunctionObject{
                Name:        "code_interpreter",
                Description: "Interpret and execute code",
                Parameters: map[string]interface{}{
                    "language": "python",
                },
            },
        },
        // Include other tools as needed
    },
    // Include any metadata relevant to the run
    Metadata: map[string]interface{}{
        "user_id": "123",
        "context": "equation_help",
    },
}

run, err := assistants.CreateRun(ctx, thread.ID, runParams)
if err != nil {
    log.Fatal(err)
}
```

### Step 5: Display the Assistant's Response

Check the status of the Run to see if it has moved to completed. Retrieve the Messages added by the Assistant to the Thread and display them to the user.

```go
runID := "your_run_id"     // Replace with the actual Run ID
threadID := "your_thread_id" // Replace with the actual Thread ID

for {
    run, err := client.RetrieveRun(context.Background(), threadID, runID)
    if err != nil {
        log.Fatal(err)
    }

    switch run.Status {
    case "completed":
        // Run completed successfully
        messages, err := client.ListMessages(context.Background(), threadID)
        if err != nil {
            log.Fatal(err)
        }

        // Display Assistant's response to the user
        for _, message := range messages {
            fmt.Printf("%s\t%s\n", message.Role, message.Content.Text.Value)
        }

        return
    case "failed":
        log.Fatal("Run failed:", run.LastError)
    case "cancelled":
        log.Fatal("Run cancelled.")
    default:
        // Wait for a few seconds before checking the status again
        time.Sleep(time.Second * 5)
    }
}
```

### Example

```go
package main

import (
"bufio"
"context"
"fmt"
"log"
"os"
"time"

assistants "github.com/ozfive/openai-assistants-go"
)

func main() {

    ctx := context.Background()

    client := assistants.NewClient("YOUR-API-KEY")

    // Create an Assistant
    assistant, err := client.CreateAssistant(ctx, &assistants.CreateAssistantRequest{
        Instructions: "You are an HR bot, and you have access to files to answer employee questions about company policies.",
        Tools:        []assistants.ToolObject{{Type: "retrieval"}},
        Model:        "gpt-4-1106-preview",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a Thread
    thread, err := client.CreateThread(ctx, assistants.CreateThreadParams{
        Messages: []assistants.MessageObject{},
        Metadata: map[string]string{"key": "value"},
    })
    if err != nil {
        log.Fatal(err)
    }

    // Start an infinite loop to continuously read user input and send messages
    scanner := bufio.NewScanner(os.Stdin)

    for {

        fmt.Print("You: ")
        scanner.Scan()
        userInput := scanner.Text()

        // Add a Message to the Thread
        message, err := client.CreateMessage(ctx, thread.ID, assistants.CreateMessageParams{
            Role:    "user",
            Content: userInput,
            FileIDs: []string{"user_file_id"},
        })
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(message.CreatedAt)

        // Run the Assistant
        tools := make([]assistants.ToolObject, len(assistant.Tools))
        for i, tool := range assistant.Tools {
            tools[i] = assistants.ToolObject{Type: tool.Type}
        }

        // Run the Assistant
        run, err := client.CreateRun(ctx, thread.ID, assistants.CreateRunParams{
            AssistantID:  assistant.ID,
            Model:        assistant.Model,
            Instructions: assistant.Instructions,
            Tools:        tools,
            Metadata:     map[string]interface{}{"user_id": "123"},
        })
        if err != nil {
            log.Fatal(err)
        }

        // Display the Assistant's Response with Typing Indicator
        runID := run.ID
        threadID := thread.ID

        for {
            run, err := client.RetrieveRun(ctx, threadID, runID)
            if err != nil {
                log.Fatal(err)
            }

            switch run.Status {
            case "completed":
                // Run completed successfully
                messages, err := client.ListMessages(ctx, threadID, assistants.ListMessagesParams{
                    Limit: "10",  // Specify your limit
                    Order: "asc", // Specify your order
                })
                if err != nil {
                    log.Fatal(err)
                }

                // Display a typing indicator
                fmt.Print("Assistant is typing")
                for i := 0; i < 3; i++ {
                    time.Sleep(time.Second)
                    fmt.Print(".")
                }
                fmt.Println()

                // Display Assistant's response to the user
                for _, msg := range messages {
                    displayAssistantResponse(&msg) // Pass the address of the message
                }

                break
            case "failed":
                log.Fatal("Run failed:", run.LastError)
            case "cancelled":
                log.Fatal("Run cancelled.")
            default:
                // Simulate typing indicator while waiting for the response
                fmt.Print("Assistant is typing")
                for i := 0; i < 3; i++ {
                    time.Sleep(time.Second)
                    fmt.Print(".")
                }
                fmt.Print("\r                    ") // Clear the typing indicator
                time.Sleep(time.Second * 2)         // Simulate processing time
            }
        }
    }
    }

    // displayAssistantResponse parses and displays the text content and entities from the assistant's message.
    func displayAssistantResponse(message *assistants.MessageObject) {
    fmt.Printf("%s\t%s\n", message.Role, message.Content[0].Text.Value)

    // Display entities if available
    if len(message.Content[0].Text.Annotations) > 0 {
        fmt.Println("Entities:")
        for _, entity := range message.Content[0].Text.Annotations {
            fmt.Printf("  Type: %s, Value: %s, Start: %d, End: %d\n", entity.Entity, entity.Value, entity.Start, entity.End)
        }
    }
}

```

## Contributing

Contributions from the community are very much encouraged. To contribute, please follow the contribution guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/ozfive/openai-assistants-go/blob/main/LICENSE) file for details.
