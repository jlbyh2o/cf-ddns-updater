# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-09-04

### Added
- **Stable Release**: First stable release of Cloudflare DDNS Updater
- Dependabot configuration for automated dependency updates
- Comprehensive GitHub labels for better issue and PR management
- Enhanced project badges (Go version, releases, downloads)

### Changed
- **Code Quality**: Refactored complex functions to reduce cyclomatic complexity
- **Maintainability**: Split large functions into focused, single-responsibility functions
- **Documentation**: Updated README with improved badges and fixed GitHub Actions badge
- **Go Version**: Updated Go requirement from 1.19+ to 1.23+ for better compatibility
- **GitHub Actions**: Updated setup-go from v5 to v6 with Go 1.23

### Fixed
- GitHub Actions build badge URL now correctly points to workflow file
- Go version compatibility issues between go.mod and GitHub Actions
- Security documentation updated to reflect 1.0.x supported versions
- Issue templates updated with current version placeholders

### Security
- Updated SECURITY.md to reflect supported versions (1.0.x and above)
- Enhanced security documentation with comprehensive threat model

## [0.3.0] - 2025-01-27

### Added
- Comprehensive GitHub repository setup with issue templates, PR templates, and workflows
- Professional documentation including CODE_OF_CONDUCT.md, CONTRIBUTING.md, and SECURITY.md
- Automated build and release process via GitHub Actions
- Support for multiple architectures (Linux: amd64, arm, arm64; Windows: amd64, arm64)
- Enhanced installation and update scripts for Linux systems
- Service file management and automatic updates
- Improved logging and configuration handling

### Changed
- Updated project structure for better maintainability
- Enhanced error handling and user feedback
- Improved documentation and examples

### Fixed
- Configuration file handling and validation
- Service management in installation scripts

## [0.2.0] - Previous Release

### Added
- Core DDNS functionality for Cloudflare
- IPv4 and IPv6 support
- Configuration file support
- Basic installation scripts

### Features
- Automatic IP detection
- Multiple domain support
- Configurable update intervals
- TTL and proxy settings

## [0.1.0] - Initial Release

### Added
- Basic Cloudflare DDNS updater functionality
- Command-line interface
- Configuration file support

[Unreleased]: https://github.com/jlbyh2o/cf-ddns-updater/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/jlbyh2o/cf-ddns-updater/compare/v0.3.0...v1.0.0
[0.3.0]: https://github.com/jlbyh2o/cf-ddns-updater/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/jlbyh2o/cf-ddns-updater/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/jlbyh2o/cf-ddns-updater/releases/tag/v0.1.0