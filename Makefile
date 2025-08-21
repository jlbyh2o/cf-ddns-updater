# Cloudflare DDNS Updater Makefile

# Variables
APP_NAME = cf-ddns-updater
VERSION = 0.2.0
BINARY = $(APP_NAME)
SOURCES = *.go

# Installation paths
PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
CONFDIR = /etc/cf-ddns
SYSTEMD_DIR = /etc/systemd/system
USER = cf-ddns
GROUP = cf-ddns

# Build flags
LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION)"
GOFLAGS = -trimpath

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY) .

# Build for all supported architectures
.PHONY: build-all
build-all:
	@echo "Building $(APP_NAME) for all architectures..."
	@mkdir -p bin
	# Linux x86-64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY)-linux-amd64 .
	# Linux ARM
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY)-linux-arm .
	# Linux ARM64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY)-linux-arm64 .
	# Windows x86-64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY)-windows-amd64.exe .
	# Windows ARM64
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY)-windows-arm64.exe .
	@echo "All builds completed successfully!"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY)
	rm -rf bin/

# Install the application
.PHONY: install
install: build
	@echo "Installing $(APP_NAME)..."
	# Create user and group
	getent group $(GROUP) >/dev/null || groupadd --system $(GROUP)
	getent passwd $(USER) >/dev/null || useradd --system --gid $(GROUP) --home-dir /var/lib/$(APP_NAME) --shell /bin/false $(USER)
	# Install binary
	install -D -m 755 $(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	# Create configuration directory
	install -d -m 755 $(DESTDIR)$(CONFDIR)
	# Install example configuration
	install -D -m 644 cf-ddns.conf.example $(DESTDIR)$(CONFDIR)/cf-ddns.conf.example
	# Set ownership
	chown -R $(USER):$(GROUP) $(DESTDIR)$(CONFDIR)
	# Install systemd service
	install -D -m 644 $(APP_NAME).service $(DESTDIR)$(SYSTEMD_DIR)/$(APP_NAME).service
	# Reload systemd
	systemctl daemon-reload
	@echo "Installation complete!"
	@echo "1. Copy $(CONFDIR)/cf-ddns.conf.example to $(CONFDIR)/cf-ddns.conf"
	@echo "2. Edit $(CONFDIR)/cf-ddns.conf with your settings"
	@echo "3. Enable and start the service: sudo systemctl enable --now $(APP_NAME)"

# Uninstall the application
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	# Stop and disable service
	systemctl stop $(APP_NAME) || true
	systemctl disable $(APP_NAME) || true
	# Remove files
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY)
	rm -f $(DESTDIR)$(SYSTEMD_DIR)/$(APP_NAME).service
	# Remove configuration directory (ask for confirmation)
	@echo "Remove configuration directory $(CONFDIR)? [y/N]"
	@read -r response; if [ "$$response" = "y" ] || [ "$$response" = "Y" ]; then rm -rf $(DESTDIR)$(CONFDIR); fi
	# Reload systemd
	systemctl daemon-reload
	@echo "Uninstallation complete!"

# Install development version (no systemd service)
.PHONY: install-dev
install-dev: build
	@echo "Installing $(APP_NAME) for development..."
	install -D -m 755 $(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	install -d -m 755 $(DESTDIR)$(CONFDIR)
	install -D -m 644 cf-ddns.conf.example $(DESTDIR)$(CONFDIR)/cf-ddns.conf.example

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary for Linux x86-64"
	@echo "  build-all  - Build for all supported architectures"
	@echo "  install    - Install the application system-wide"
	@echo "  install-dev- Install without systemd service"
	@echo "  uninstall  - Remove the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Installation paths:"
	@echo "  Binary:        $(BINDIR)/$(BINARY)"
	@echo "  Configuration: $(CONFDIR)/"
	@echo "  Systemd:       $(SYSTEMD_DIR)/$(APP_NAME).service"