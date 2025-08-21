# Contributing to Cloudflare DDNS Updater

Thank you for your interest in contributing to Cloudflare DDNS Updater! We welcome contributions from the community and are grateful for your help in making this project better.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)
- [Getting Help](#getting-help)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- **Go**: Version 1.21 or later
- **Git**: For version control
- **Make**: For build automation (optional but recommended)
- **GitHub Account**: For submitting contributions

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/cf-ddns-updater.git
   cd cf-ddns-updater
   ```
3. Add the original repository as upstream:
   ```bash
   git remote add upstream https://github.com/jlbyh2o/cf-ddns-updater.git
   ```

## Development Setup

### Building from Source

```bash
# Install dependencies
go mod download

# Build the application
go build -o cf-ddns-updater

# Or use the Makefile
make build
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Development Configuration

Create a test configuration file:

```bash
cp config.example.json config.test.json
# Edit config.test.json with your test Cloudflare credentials
```

**⚠️ Important**: Never commit real API tokens or credentials!

## How to Contribute

### Types of Contributions

We welcome several types of contributions:

- **🐛 Bug Reports**: Help us identify and fix issues
- **✨ Feature Requests**: Suggest new functionality
- **📝 Documentation**: Improve or add documentation
- **🔧 Code Contributions**: Fix bugs or implement features
- **🧪 Testing**: Add or improve tests
- **🎨 UI/UX**: Improve user experience

### Before You Start

1. **Check existing issues**: Look for existing issues or discussions
2. **Create an issue**: For significant changes, create an issue first
3. **Discuss**: Get feedback on your proposed changes
4. **Assign yourself**: Comment on the issue to avoid duplicate work

### Working on Issues

1. **Choose an issue**: Look for issues labeled `good first issue` or `help wanted`
2. **Create a branch**: Use a descriptive branch name
   ```bash
   git checkout -b feature/add-ipv6-support
   git checkout -b fix/config-validation-bug
   git checkout -b docs/update-readme
   ```
3. **Make changes**: Follow our coding standards
4. **Test thoroughly**: Ensure your changes work as expected
5. **Commit**: Use clear, descriptive commit messages

## Coding Standards

### Go Style Guide

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use `golint` and `go vet` to check for issues
- Write clear, self-documenting code

### Code Organization

```
cf-ddns-updater/
├── main.go           # Application entry point
├── config.go         # Configuration handling
├── cloudflare.go     # Cloudflare API interactions
├── ip.go            # IP detection logic
├── updater.go       # Core update logic
└── *_test.go        # Test files
```

### Naming Conventions

- **Functions**: Use camelCase (`updateDNSRecord`)
- **Variables**: Use camelCase (`apiToken`)
- **Constants**: Use UPPER_SNAKE_CASE (`DEFAULT_TIMEOUT`)
- **Files**: Use snake_case (`dns_updater.go`)

### Documentation

- Add comments for exported functions and types
- Use godoc format for documentation comments
- Include examples in documentation when helpful

```go
// UpdateDNSRecord updates a DNS record in Cloudflare.
// It returns an error if the update fails.
//
// Example:
//   err := UpdateDNSRecord("example.com", "192.168.1.1")
//   if err != nil {
//       log.Fatal(err)
//   }
func UpdateDNSRecord(domain, ip string) error {
    // Implementation
}
```

### Error Handling

- Use descriptive error messages
- Wrap errors with context using `fmt.Errorf`
- Handle errors at the appropriate level

```go
if err != nil {
    return fmt.Errorf("failed to update DNS record for %s: %w", domain, err)
}
```

## Testing

### Test Requirements

- Write tests for new functionality
- Maintain or improve test coverage
- Include both unit tests and integration tests
- Test error conditions and edge cases

### Test Structure

```go
func TestUpdateDNSRecord(t *testing.T) {
    tests := []struct {
        name     string
        domain   string
        ip       string
        wantErr  bool
    }{
        {
            name:    "valid update",
            domain:  "example.com",
            ip:      "192.168.1.1",
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := UpdateDNSRecord(tt.domain, tt.ip)
            if (err != nil) != tt.wantErr {
                t.Errorf("UpdateDNSRecord() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Running Specific Tests

```bash
# Run tests for a specific package
go test ./config

# Run a specific test function
go test -run TestUpdateDNSRecord

# Run tests with race detection
go test -race ./...
```

## Submitting Changes

### Commit Messages

Use clear, descriptive commit messages following this format:

```
type(scope): brief description

Detailed explanation of the change (if needed)

Fixes #123
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(api): add support for IPv6 addresses

fix(config): validate domain names properly

docs(readme): update installation instructions
```

### Pull Request Process

1. **Update your branch**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push your changes**:
   ```bash
   git push origin your-branch-name
   ```

3. **Create a Pull Request**:
   - Use a clear, descriptive title
   - Fill out the PR template completely
   - Link related issues
   - Add screenshots for UI changes

4. **Address feedback**:
   - Respond to review comments
   - Make requested changes
   - Push updates to your branch

### Pull Request Checklist

- [ ] Code follows the project's coding standards
- [ ] Tests pass locally
- [ ] New functionality includes tests
- [ ] Documentation is updated (if applicable)
- [ ] Commit messages are clear and descriptive
- [ ] PR description explains the changes
- [ ] Related issues are linked

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Workflow

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create a git tag
4. GitHub Actions automatically builds and releases binaries

## Getting Help

### Communication Channels

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Pull Request Reviews**: For code-specific discussions

### Resources

- [Go Documentation](https://golang.org/doc/)
- [Cloudflare API Documentation](https://developers.cloudflare.com/api/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)

### Maintainer Response Times

- **Issues**: We aim to respond within 48 hours
- **Pull Requests**: Initial review within 7 days
- **Security Issues**: Within 24 hours

## Recognition

Contributors are recognized in:

- Release notes
- GitHub contributors list
- Special mentions for significant contributions

---

## Quick Start for Contributors

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/cf-ddns-updater.git
cd cf-ddns-updater

# 2. Set up development environment
go mod download
cp config.example.json config.test.json

# 3. Create a feature branch
git checkout -b feature/your-feature-name

# 4. Make changes and test
go test ./...
go build

# 5. Commit and push
git add .
git commit -m "feat: add your feature"
git push origin feature/your-feature-name

# 6. Create a Pull Request on GitHub
```

Thank you for contributing to Cloudflare DDNS Updater! 🎉