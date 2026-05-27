# Installation

## Download binary

Download the latest binary for your platform from [GitHub Releases](https://github.com/vngcloud/greennode-cli/releases):

### macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-darwin-arm64

# Intel
curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-darwin-amd64

chmod +x grn
sudo mv grn /usr/local/bin/
```

### Linux

```bash
# x86_64
curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-linux-amd64

# ARM64
curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-linux-arm64

chmod +x grn
sudo mv grn /usr/local/bin/
```

### Windows

Download `grn-windows-amd64.exe` from [GitHub Releases](https://github.com/vngcloud/greennode-cli/releases) and add to your PATH.

## Build from source

Requires [Go 1.22+](https://go.dev/dl/):

```bash
git clone https://github.com/vngcloud/greennode-cli.git
cd greennode-cli/go
go build -o grn .
sudo mv grn /usr/local/bin/
```

## Verify installation

```bash
grn --version
# grn-cli/0.1.0 Go/1.22.2 darwin/arm64
```
