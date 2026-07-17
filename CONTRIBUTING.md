# Contributing to InterSpec

Thanks for your interest in contributing! This guide covers how to make changes to the specification, skills, and documentation.

## Setup

```bash
git clone https://github.com/DaviMGDev/InterSpec.git
cd InterSpec
```

No build tools or dependencies are required — all source files are Markdown (`.md`), InterSpec (`.is`), CSS, and JavaScript.

## Making Changes

This project contains:

- **`LANGUAGE.md`** — The full InterSpec language specification
- **`skills/`** — Agent skill definitions that teach AI agents how to work with `.is` files
- **`assets/`** — Demo assets and design token files
- **Root `.md` files** — README, CONTRIBUTING, LICENSE

When editing documentation:
- Ensure Markdown is valid and renders correctly
- Keep cross-references (file links) accurate
- Preserve the existing tone and structure
- Follow the language spec's own conventions (PascalCase for components, camelCase for variables)

## Commit Conventions

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(spec): add new built-in component
docs(skills): update interspec-consume viewport rules
fix: correct typo in LANGUAGE.md hints section
chore: update README with current project structure
```

Keep commits focused and atomic — one logical change per commit.

## Pull Requests

1. Fork the repository and create a branch from `main`.
2. Make your changes.
3. Review the diff to confirm only intended changes are included.
4. Push your branch and open a PR against `main`.
5. Describe what the PR does and why. Link any related issues.

## License

By contributing, you agree that your contributions will be licensed under the
MIT License (see [LICENSE](LICENSE)).
