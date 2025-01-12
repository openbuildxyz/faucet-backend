# Faucet Backend

This is the backend service for the OpenBuild Faucet project, built using **Gin** and **go-ethereum** (Geth). **OpenBuild** is the open-source community that supports this project.

## Dependencies 

- **Go**: 1.22
- **Gin**: 1.10
- **go-ethereum**: 1.14

## Quick Start

### Install Dependencies

```bash
go mod tidy
```

### Configure

```bash
cp config.example.yaml config.yaml
```

The configuration file must be written according to your specific setup.

### Run

```bash
go run main.go
```
By default, the server will run on http://localhost:8080.


