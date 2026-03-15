# Development scripts

# Default command
_default:
    @just --list --unsorted

# Sync Go modules
tidy:
    go mod tidy
    @echo "All modules synced, Go workspace ready!"

# CLI local run wrapper
jiff *args:
    @go run . {{ args }}

# Run all BDD tests
test:
    @echo "Running unit tests!"
    go clean -testcache
    go test -cover ./...

# Build the binary
build:
    #!/usr/bin/env bash
    # Detect OS and architecture
    case "$(uname -s)" in
        Linux*) OS="linux" ;;
        Darwin*) OS="darwin" ;;
        *) echo "Error: Unsupported OS (${OS})"; exit 1 ;;
    esac
    case "$(uname -m)" in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        arm64) ARCH="arm64" ;;
        *) echo "Error: Unsupported architecture (${ENV_ARCH})"; exit 1 ;;
    esac

    echo "Building jiff for ${OS}/${ARCH}..."
    go mod download all
    CGO_ENABLED=0 GOOS="${OS}" GOARCH="${ARCH}" go build -o ./jiff .
    echo "Built binary for jiff successfully!"

# Update the project dependencies
update-deps:
    @echo "Updating project dependencies..."
    go get -u ./...
    go mod tidy
