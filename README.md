# gopeed-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/gabrielramos02/gopeed-api-go.svg)](https://pkg.go.dev/github.com/gabrielramos02/gopeed-api-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabrielramos02/gopeed-api-go)](https://goreportcard.com/report/github.com/gabrielramos02/gopeed-api-go)

Go client for the [Gopeed](https://github.com/GopeedLab/gopeed) HTTP API.

## Features

- Resolve URLs into downloadable resources
- Create, list, get and delete tasks
- Configurable HTTP client, timeout and API token
- Context support for cancellation and timeouts
- Zero external dependencies beyond the Go standard library

## Installation

```bash
go get github.com/gabrielramos02/gopeed-api-go
```

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/gabrielramos02/gopeed-api-go"
)

func main() {
    client, err := gopeed.NewClient(
        "http://localhost:9999",
        gopeed.WithAPIToken("my-token"),
        gopeed.WithTimeout(30*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    taskID, err := client.CreateTaskFromURL(ctx, "https://example.com/file.zip", gopeed.GopeedOptions{
        Name: "file.zip",
        Path: "/downloads",
        Extra: &gopeed.GopeedExtraOptions{Connections: 16},
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Created task:", taskID)
}
```

## Usage

### Create a client

```go
client, err := gopeed.NewClient("http://localhost:9999")
```

If your Gopeed server requires authentication:

```go
client, err := gopeed.NewClient(
    "http://localhost:9999",
    gopeed.WithAPIToken("my-token"),
)
```

### Resolve then create a task

```go
resolved, err := client.Resolve(ctx, "https://example.com/file.zip", gopeed.GopeedOptions{})
if err != nil {
    log.Fatal(err)
}

taskID, err := client.CreateTask(ctx, resolved.ID, gopeed.GopeedOptions{
    Name: "file.zip",
})
```

### List, get and delete tasks

```go
tasks, err := client.GetTasks(ctx)

task, err := client.GetTask(ctx, taskID)

err = client.DeleteTask(ctx, taskID)
```

## Client options

| Option | Description |
|---|---|
| `WithAPIToken(token string)` | Sets the `X-Api-Token` header |
| `WithTimeout(d time.Duration)` | Sets a timeout on the default HTTP client |
| `WithHTTPClient(client *http.Client)` | Uses a custom `http.Client` |

## API coverage

| Method | Description |
|---|---|
| `GetInfo` | Server version and runtime info |
| `GetTasks` | List all tasks |
| `GetTask` | Get a single task |
| `Resolve` | Resolve a URL into a resource |
| `CreateTask` | Create a task from a resolved resource ID |
| `CreateTaskFromURL` | Resolve + create in one call |
| `DeleteTask` | Delete a task |

## Options

```go
type GopeedOptions struct {
    Name        string
    Path        string
    SelectFiles []int
    Extra       *GopeedExtraOptions
}
```

- `Name`: desired file name
- `Path`: download destination
- `SelectFiles`: indexes of files to download from multi-file resources
- `Extra.Connections`: number of concurrent connections

## Contributing

Contributions are welcome. Please open an issue or pull request.

## License

[MIT](LICENSE)
