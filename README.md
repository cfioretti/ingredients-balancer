# Ingredients Balancer

This project implements a gRPC service for balancing ingredients in recipes based on pan sizes.

## Overview

The Ingredients Balancer service provides a `Balance` RPC method that takes a recipe and a set of pans as input, and returns a balanced recipe with ingredients adjusted according to the pan sizes.

## Project Structure

- `cmd/main.go`: Entry point for the application
- `pkg/application/ingredients_balancer_service.go`: Core business logic for balancing ingredients
- `pkg/infrastructure/grpc/server.go`: gRPC server implementation
- `pkg/infrastructure/grpc/proto/ingredients_balancer.proto`: Proto definition for the service
- `test/ingredients_balancer_integration_test.go`: Integration tests for the service

## Prerequisites

Before you can run the service, you need to install the Protocol Buffers compiler (protoc) and the Go plugins.

### 1. Install Protocol Buffers Compiler (protoc)

#### macOS (using Homebrew)
```bash
brew install protobuf
```

#### Linux (Ubuntu/Debian)
```bash
apt-get update
apt-get install -y protobuf-compiler
```

#### Windows (using Chocolatey)
```bash
choco install protoc
```

Alternatively, you can download the pre-compiled binaries from the [Protocol Buffers GitHub releases page](https://github.com/protocolbuffers/protobuf/releases).

### 2. Install Go Plugins

You need to install two Go plugins for protoc:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Make sure that your `$GOPATH/bin` directory is in your `$PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Generating the Proto Files

Once you have installed the prerequisites, you can generate the proto files by running:

```bash
make proto
```

This will execute the following command defined in the Makefile:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/infrastructure/grpc/proto/ingredients_balancer.proto
```

The generated files will be placed in the `pkg/infrastructure/grpc/proto/generated` directory:

1. `ingredients_balancer.pb.go`: Contains the Go structs for the messages defined in the proto file
2. `ingredients_balancer_grpc.pb.go`: Contains the gRPC client and server interfaces

## Running the Service

After generating the proto files, you can run the service with:

```bash
make run
```

The service will listen on port 50052 by default, or on the port specified by the `PORT` environment variable.

## Running Tests

To run the tests, use:

```bash
make test
```

This will run both unit tests and integration tests to ensure that the gRPC service is working correctly.

## Docker

The project includes a Dockerfile that can be used to build and run the service in a container:

```bash
docker build -t ingredients-balancer -f deployments/Dockerfile .
docker run -p 50052:50052 ingredients-balancer
```
