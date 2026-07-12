# Contributing to InterSpec

Thanks for your interest in contributing! This guide covers how to set up the
project, make changes, and submit them.

## Setup

```bash
git clone https://github.com/DaviMGDev/InterSpec.git
cd InterSpec/isc
go build -o isc .
```

The only dependency is Go 1.24+. No third-party packages are used.

## Building

```bash
cd isc
go build -o isc .
```

## Testing

```bash
cd isc
go test ./...            # Run all tests
go test -v ./...         # Verbose output
go vet ./...             # Static analysis
```

All tests must pass before submitting a PR. Test coverage includes unit tests
for the `stripper` package and integration tests for the CLI entry point.

## Commit Conventions

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(isc): add new command
fix(stripper): correct edge case in block comment parsing
docs: update LANGUAGE.md
chore: bump Go version
```

Keep commits focused and atomic — one logical change per commit.

## Pull Requests

1. Fork the repository and create a branch from `main`.
2. Make your changes. Add tests if applicable.
3. Run `go test ./...` and `go vet ./...` to confirm nothing is broken.
4. Push your branch and open a PR against `main`.
5. Describe what the PR does and why. Link any related issues.

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`).
- Add package-level doc comments for new packages.
- Use descriptive variable names — the project favors clarity over brevity.
- Keep the `stripper` package free of third-party dependencies.

## License

By contributing, you agree that your contributions will be licensed under the
MIT License (see [LICENSE](LICENSE)).
