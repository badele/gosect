# Development Guide

This guide covers everything you need to know to contribute to gosect.

## Table of Contents

- [Building from Source](#building-from-source)
- [Development Environment](#development-environment)
- [Testing](#testing)
- [Pre-commit Hooks](#pre-commit-hooks)
- [Commit Message Conventions](#commit-message-conventions)
- [CI/CD Pipeline](#cicd-pipeline)
- [Release Process](#release-process)

## Building from Source

### Prerequisites

- Go 1.25.1 or later
- (Optional) Nix for reproducible development environment

### Standard Build

```bash
# Clone the repository
git clone https://github.com/badele/gosect.git
cd gosect

# Install dependencies
go mod download

# Build
go build -o gosect .

# Run
./gosect
```

### Build with Just

```bash
# Initialize Go module
just go-init

# Build
just go-build

# Run tests
just go-test

# Serve
just serve
```

### Build with Nix

```bash
# Enter development environment
nix develop

# Build the package
nix build

# Run
./result/bin/gosect
```

## Development Environment

### Using Nix (Recommended)

The project includes a Nix flake for a reproducible development environment:

```bash
# Enter the development shell
nix develop

# All tools are now available:
# - Go 1.25.1
# - gopls (LSP)
# - gotools (goimports, etc.)
# - just
# - pre-commit
# - hadolint
```

### Manual Setup

Install the following tools:

- Go 1.25.1+
- pre-commit
- hadolint (for Dockerfile linting)
- just (task runner)

## Testing

### Running Tests

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt

# Run specific test
go test -v -run TestMdToHTML

# Run tests with race detection (pre-commit does this)
go test -v -race
```

### Writing Tests

All tests are in `main_test.go`. When adding new features:

1. Write tests first (TDD)
2. Ensure tests pass locally
3. Pre-commit hooks will run tests automatically

## Docker

### Build Image

```bash
# Standard build
docker build -t badele/gosect:latest .

# Or with just
just docker-build
```

### Multi-arch Build

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t badele/gosect:latest \
  --push .
```

### Test Docker Image

```bash
# Run locally
just docker-run

# Custom port
just docker-run PORT=3000
```

## Pre-commit Hooks

This project uses pre-commit hooks to ensure code quality and commit message
standards.

### Installation

```bash
nix develop
just precommit-install
```

### What gets checked?

#### On every commit:

- **Go formatting** (`go fmt`)
- **Go vet** (code correctness)
- **Go imports** (organize imports)
- **Go tests** (run all tests with race detection)
- **Go build** (ensure project compiles)
- **YAML/JSON syntax** validation
- **Merge conflicts** detection
- **File endings** (EOF, trailing whitespace)
- **Large files** detection
- **Dockerfile** linting

#### On commit message:

- **Conventional Commits** format validation

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/)
specification:

```
<type>: <description>

[optional body]

[optional footer(s)]
```

#### Types:

- `feat`: New feature (triggers minor version bump)
- `fix`: Bug fix (triggers patch version bump)
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes
- `build`: Build system changes

### Manual Testing

You can run the pre-commit hooks manually:

```bash
# Run on all files
pre-commit run --all-files

# Run specific hook
pre-commit run go-test --all-files

# Skip hooks for a commit (not recommended)
git commit --no-verify -m "fix: emergency hotfix"
```

### Bypassing Hooks

If you need to bypass hooks (not recommended):

```bash
# Skip pre-commit hooks
git commit --no-verify

# Skip commit-msg hooks
SKIP=conventional-pre-commit git commit -m "your message"
```

### Troubleshooting

#### Hook installation fails

```bash
# Clean and reinstall
pre-commit clean
pre-commit install --install-hooks
```

#### Tests fail during commit

Fix the tests before committing. You can run tests manually:

```bash
go test -v
```

#### Commit message rejected

Ensure your commit message follows the Conventional Commits format. See examples
above.

## Development Workflow

### 1. Make changes

```bash
# Edit files
vim main.go
```

### 2. Test locally

```bash
# Run tests
go test -v

# Build
go build
```

### 3. Commit with conventional format

```bash
git add .
git commit -m "feat: add new feature"
```

The pre-commit hooks will:

1. Format your code
2. Run tests
3. Check code quality
4. Validate commit message

### 4. Push and create PR

```bash
git push origin main
```

## CI/CD Pipeline

Every push triggers:

- Unit tests
- Build verification
- Docker image build (multi-arch)
- Docker push to registry (on main branch)

## Release Process

We use [release-please](https://github.com/googleapis/release-please) for
automated releases.

### How It Works

1. **Make commits** with conventional format (e.g., `feat:`, `fix:`)
2. **Push to main** branch
3. **release-please** automatically creates/updates a release PR with:
   - Generated changelog
   - Version bump in `.release-please-manifest.json`
4. **Review and merge** the release PR
5. **Automatic release**:
   - Git tag created (e.g., `v1.2.3`)
   - GitHub release published
   - CI/CD builds and pushes Docker images with version tags

### Version Bumping Rules

Based on Conventional Commits:

- `feat:` → Minor version (1.0.0 → 1.1.0)
- `fix:` → Patch version (1.0.0 → 1.0.1)
- `feat!:` or `BREAKING CHANGE:` → Major version (1.0.0 → 2.0.0)
- `docs:`, `chore:`, etc. → No version bump

### Example Release Workflow

```bash
# 1. Make changes and commit
git commit -m "feat: add custom port support"
git commit -m "fix: correct CSS rendering"
git push origin main

# 2. Wait for release-please to create PR
# PR will be titled: "chore(main): release 1.1.0"

# 3. Review the generated CHANGELOG.md in the PR

# 4. Merge the PR
# → Tag v1.1.0 created
# → GitHub release published
# → Docker images pushed:
#    - badele/gosect:1.1.0
#    - badele/gosect:1.1
#    - badele/gosect:1
#    - badele/gosect:latest
```

## Configuration Files

Overview of important configuration files:

- `.pre-commit-config.yaml` - Pre-commit hooks configuration
- `.commitlintrc.json` - Commit message validation rules
- `release-please-config.json` - Release automation config
- `.release-please-manifest.json` - Current version tracking
- `Dockerfile` - Multi-stage Docker build
- `flake.nix` - Nix development environment
- `justfile` - Task runner commands
- `.github/workflows/` - CI/CD pipelines

## Troubleshooting

### Tests Fail

```bash
# Run tests to see detailed output
go test -v

# Check for race conditions
go test -v -race
```

### Docker Build Fails

```bash
# Check Dockerfile syntax
hadolint Dockerfile

# Build with verbose output
docker build --progress=plain -t gosect:test .
```

### Pre-commit Hooks Not Running

```bash
# Reinstall hooks
pre-commit uninstall
pre-commit install
pre-commit install --hook-type commit-msg

# Test manually
pre-commit run --all-files
```

## Getting Help

- Check existing [GitHub Issues](https://github.com/badele/gosect/issues)
- Create a new issue with detailed description
- Join discussions in Pull Requests
