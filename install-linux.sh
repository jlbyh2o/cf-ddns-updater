#!/bin/bash

# Cloudflare DDNS Updater - Linux AMD64 Install Script
# Usage: curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | bash

set -e

# Configuration
REPO="jlbyh2o/cf-ddns-updater"
BINARY_NAME="cf-ddns-updater"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/cf-ddns-updater"
SERVICE_DIR="/etc/systemd/system"
VERSION="latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Detect architecture
detect_arch() {
    local arch=$(uname -m)
    case $arch in
        x86_64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7l)
            echo "arm"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        log_error "Neither curl nor wget is available"
        exit 1
    fi
}

# Download binary
download_binary() {
    local version=$1
    local arch=$2
    local binary_name="cf-ddns-updater-linux-${arch}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${binary_name}"
    local temp_file="/tmp/${binary_name}"
    
    log_info "Downloading ${binary_name} from ${download_url}"
    
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_file" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_file" "$download_url"
    else
        log_error "Neither curl nor wget is available"
        exit 1
    fi
    
    if [[ ! -f "$temp_file" ]]; then
        log_error "Failed to download binary"
        exit 1
    fi
    
    echo "$temp_file"
}

# Install binary
install_binary() {
    local temp_file=$1
    
    log_info "Installing binary to ${INSTALL_DIR}/${BINARY_NAME}"
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    # Copy and set permissions
    cp "$temp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Clean up temp file
    rm -f "$temp_file"
    
    log_success "Binary installed successfully"
}

# Create config directory
create_config_dir() {
    log_info "Creating configuration directory at ${CONFIG_DIR}"
    mkdir -p "$CONFIG_DIR"
    
    # Download example config if available
    local config_url="https://raw.githubusercontent.com/${REPO}/main/config.example.json"
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "${CONFIG_DIR}/config.example.json" "$config_url" 2>/dev/null || true
    elif command -v wget >/dev/null 2>&1; then
        wget -O "${CONFIG_DIR}/config.example.json" "$config_url" 2>/dev/null || true
    fi
    
    log_success "Configuration directory created"
}

# Install systemd service
install_service() {
    log_info "Installing systemd service"
    
    local service_url="https://raw.githubusercontent.com/${REPO}/main/cf-ddns-updater.service"
    local temp_service="/tmp/cf-ddns-updater.service"
    
    # Download service file
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_service" "$service_url" 2>/dev/null || {
            log_warning "Could not download service file, creating basic one"
            create_basic_service "$temp_service"
        }
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_service" "$service_url" 2>/dev/null || {
            log_warning "Could not download service file, creating basic one"
            create_basic_service "$temp_service"
        }
    else
        create_basic_service "$temp_service"
    fi
    
    # Install service file
    cp "$temp_service" "${SERVICE_DIR}/cf-ddns-updater.service"
    rm -f "$temp_service"
    
    # Reload systemd and enable service
    systemctl daemon-reload
    systemctl enable cf-ddns-updater.service
    
    log_success "Systemd service installed and enabled"
}

# Create basic service file if download fails
create_basic_service() {
    local service_file=$1
    
    cat > "$service_file" << EOF
[Unit]
Description=Cloudflare DDNS Updater
After=network.target

[Service]
Type=simple
User=root
ExecStart=${INSTALL_DIR}/${BINARY_NAME} -config ${CONFIG_DIR}/config.json
Restart=always
RestartSec=30

[Install]
WantedBy=multi-user.target
EOF
}

# Main installation function
main() {
    log_info "Starting Cloudflare DDNS Updater installation"
    
    # Check prerequisites
    check_root
    
    # Detect system architecture
    local arch=$(detect_arch)
    log_info "Detected architecture: $arch"
    
    # Get latest version
    local version=$(get_latest_version)
    if [[ -z "$version" ]]; then
        log_error "Could not determine latest version"
        exit 1
    fi
    log_info "Latest version: $version"
    
    # Download binary
    local temp_file=$(download_binary "$version" "$arch")
    
    # Install binary
    install_binary "$temp_file"
    
    # Create config directory
    create_config_dir
    
    # Install systemd service
    install_service
    
    # Final success message
    log_success "Cloudflare DDNS Updater installed successfully!"
    echo
    log_info "Next steps:"
    echo "  1. Copy your configuration: cp ${CONFIG_DIR}/config.example.json ${CONFIG_DIR}/config.json"
    echo "  2. Edit the configuration: nano ${CONFIG_DIR}/config.json"
    echo "  3. Start the service: systemctl start cf-ddns-updater"
    echo "  4. Check service status: systemctl status cf-ddns-updater"
    echo "  5. View logs: journalctl -u cf-ddns-updater -f"
    echo
    log_info "Manual usage: ${INSTALL_DIR}/${BINARY_NAME} -config ${CONFIG_DIR}/config.json"
}

# Run main function
main "$@"