# Contributing to Ferret

Thank you for your interest in contributing to Ferret! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

### Prerequisites

- Go 1.23 or later
- Make (optional, for running Makefile commands)

### Setting Up the Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/ferret.git
   cd ferret
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/MontFerret/ferret.git
   ```

### Building

```bash
go build ./...
```

### Running Tests

```bash
go test ./...
```

Or using Make:
```bash
make test
```

## How to Contribute

### Reporting Bugs

- Check if the bug has already been reported in [Issues](https://github.com/MontFerret/ferret/issues)
- If not, create a new issue with:
  - Clear title and description
  - Steps to reproduce
  - Expected vs actual behavior
  - Your environment (OS, Go version, Ferret version)

### Suggesting Features

- Open an issue describing the feature
- Explain the use case and benefits
- Be open to discussion about implementation approaches

### Pull Requests

1. Create a branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the coding guidelines below

3. Write or update tests as needed

4. Ensure all tests pass:
   ```bash
   go test ./...
   ```

5. Run the linter:
   ```bash
   golangci-lint run
   ```

6. Commit your changes with a clear message:
   ```bash
   git commit -m "Add feature X that does Y"
   ```

7. Push to your fork and create a Pull Request

### Coding Guidelines

- Follow standard Go conventions and formatting (`gofmt`)
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions focused and small
- Write unit tests for new functionality
- Update documentation when changing behavior

## Community

- [Discord](https://discord.gg/kzet32U)
- [Telegram](https://t.me/montferret_chat)

## License

By contributing to Ferret, you agree that your contributions will be licensed under the Apache 2.0 License.
