# Contributing

Thanks for your interest in contributing to drim.

## Development Setup

Prerequisites:
- Go 1.23 or later
- Docker for testing

Clone and build:

```bash
git clone https://github.com/usekaneo/drim.git
cd drim
go mod download
go build -o drim .
```

## Project Structure

```
drim/
├── main.go              # Entry point
├── cmd/                 # CLI commands
├── pkg/
│   ├── banner/         # UI banner
│   ├── docker/         # Docker operations
│   ├── generator/      # File generation
│   └── ui/             # User prompts
```

## Making Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Format code: `go fmt ./...`
5. Test: `go test ./...`
6. Commit: `git commit -m "Add feature"`
7. Push: `git push origin feature-name`
8. Open a Pull Request

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and simple
- Add tests for new functionality

## Testing

Run tests:

```bash
go test ./...
```

Test locally:

```bash
go build -o drim .
./drim setup
```

## Pull Requests

- Provide a clear description
- Reference related issues
- Ensure all tests pass
- Keep changes focused

## Questions

Open an issue for:
- Bug reports
- Feature requests
- Questions about development

## License

By contributing, you agree your contributions will be licensed under the MIT License.
