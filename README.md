# go-metric

A simple Go package for OpenTelemetry metrics with OTLP HTTP export.

## Installation

```bash
go get github.com/gokpm/go-metric
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    
    "github.com/gokpm/go-metric"
)

func setup() error {
    config := metric.Config{
        Ok:          true,
        Name:        "my-service",
        Environment: "production",
        URL:         "http://localhost:4318/v1/metrics",
    }
    
    ctx := context.Background()
    _, err := metric.Setup(ctx, config)
    return err
}

func main() {
    if err := setup(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    defer metric.Shutdown(5 * time.Second)
    
    // Use meter for metrics...
}
```

## Configuration

- `Ok`: Enable/disable metrics
- `Name`: Service name
- `Environment`: Deployment environment
- `URL`: OTLP HTTP endpoint URL (default: `http://localhost:4318/v1/metrics`)

## Features

- OTLP HTTP export with gzip compression
- Automatic resource detection (hostname, service info)
- Periodic metric collection
- Graceful shutdown with timeout