#!/bin/bash

# Build script for drim
# Builds binaries for multiple platforms

set -e

VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

echo "Building drim version ${VERSION}..."
echo "Build time: ${BUILD_TIME}"
echo ""

# Create bin directory
mkdir -p bin

# Build for Linux
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/drim_Linux_x86_64 .

echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o bin/drim_Linux_arm64 .

# Build for macOS
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/drim_Darwin_x86_64 .

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o bin/drim_Darwin_arm64 .

# Build for Windows
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/drim_Windows_x86_64.exe .

echo ""
echo "Build complete! Binaries are in the bin/ directory:"
ls -lh bin/

echo ""
echo "To create a release, run:"
echo "  git tag v0.1.0"
echo "  git push origin v0.1.0"


