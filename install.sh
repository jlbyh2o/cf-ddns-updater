#!/bin/bash

# Cloudflare DDNS Updater Installation Script
# This script installs the cf-ddns-updater on Linux systems

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="cf-ddns-updater"
VERSION="1.0"
BINARY_NAME="cf-ddns-updater"
SERVICE_NAME="cf-ddns-updater"
USER_NAME="cf-ddns"
GROUP_NAME="cf-ddns"

# Installation paths
BIN_DIR="/usr/local/bin"
CONF_DIR="/etc/cf-ddns"
SYSTEMD_DIR="/etc/systemd/system"
LOG_DIR="/var/log"

# Functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

check_dependencies() {
    print_info "Checking dependencies..."
    
    # Check if systemd is available
    if ! command -v systemctl &> /dev/null; then
        print_error "systemctl not found. This script requires systemd."
        exit 1
    fi
    
    # Check if make is available
    if ! command -v make &> /dev/null; then
        print_warning "make not found. You may need to install build-essential or similar package."
    fi
    
    print_success "Dependencies check passed"
}

check_existing_installation() {
    if [[ -f "$BIN_DIR/$BINARY_NAME" ]]; then
        print_warning "Existing installation found at $BIN_DIR/$BINARY_NAME"
        read -p "Do you want to continue and overwrite? [y/N]: " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled"
            exit 0
        fi
        
        # Stop service if running
        if systemctl is-active --quiet "$SERVICE_NAME"; then
            print_info "Stopping existing service..."
            systemctl stop "$SERVICE_NAME"
        fi
    fi
}

create_user() {
    print_info "Creating system user and group..."
    
    # Create group if it doesn't exist
    if ! getent group "$GROUP_NAME" >/dev/null 2>&1; then
        groupadd --system "$GROUP_NAME"
        print_success "Created group: $GROUP_NAME"
    else
        print_info "Group $GROUP_NAME already exists"
    fi
    
    # Create user if it doesn't exist
    if ! getent passwd "$USER_NAME" >/dev/null 2>&1; then
        useradd --system --gid "$GROUP_NAME" --home-dir "/var/lib/$APP_NAME" --shell /bin/false "$USER_NAME"
        print_success "Created user: $USER_NAME"
    else
        print_info "User $USER_NAME already exists"
    fi
}

install_binary() {
    print_info "Installing binary..."
    
    if [[ ! -f "$BINARY_NAME" ]]; then
        print_error "Binary $BINARY_NAME not found. Please build it first with 'make build'"
        exit 1
    fi
    
    install -D -m 755 "$BINARY_NAME" "$BIN_DIR/$BINARY_NAME"
    print_success "Binary installed to $BIN_DIR/$BINARY_NAME"
}

setup_config() {
    print_info "Setting up configuration..."
    
    # Create config directory
    install -d -m 755 "$CONF_DIR"
    
    # Install example configuration
    if [[ -f "cf-ddns.conf.example" ]]; then
        install -D -m 644 cf-ddns.conf.example "$CONF_DIR/cf-ddns.conf.example"
        print_success "Example configuration installed to $CONF_DIR/cf-ddns.conf.example"
    else
        print_warning "cf-ddns.conf.example not found, skipping example config installation"
    fi
    
    # Set ownership
    chown -R "$USER_NAME:$GROUP_NAME" "$CONF_DIR"
    print_success "Configuration directory ownership set"
}

install_service() {
    print_info "Installing systemd service..."
    
    if [[ ! -f "$SERVICE_NAME.service" ]]; then
        print_error "Service file $SERVICE_NAME.service not found"
        exit 1
    fi
    
    install -D -m 644 "$SERVICE_NAME.service" "$SYSTEMD_DIR/$SERVICE_NAME.service"
    systemctl daemon-reload
    print_success "Systemd service installed"
}

post_install() {
    print_info "Post-installation setup..."
    
    # Create log directory if it doesn't exist
    mkdir -p "$LOG_DIR"
    
    print_success "Installation completed successfully!"
    echo
    print_info "Next steps:"
    echo "1. Copy the example configuration:"
    echo "   sudo cp $CONF_DIR/cf-ddns.conf.example $CONF_DIR/cf-ddns.conf"
    echo "   # Edit the configuration file:"
    echo "   sudo nano $CONF_DIR/cf-ddns.conf"
    echo
    echo "3. Enable and start the service:"
    echo "   sudo systemctl enable $SERVICE_NAME"
    echo "   sudo systemctl start $SERVICE_NAME"
    echo
    echo "4. Check service status:"
    echo "   sudo systemctl status $SERVICE_NAME"
    echo
    echo "5. View logs:"
    echo "   sudo journalctl -u $SERVICE_NAME -f"
}

uninstall() {
    print_info "Uninstalling $APP_NAME..."
    
    # Stop and disable service
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        systemctl stop "$SERVICE_NAME"
        print_info "Service stopped"
    fi
    
    if systemctl is-enabled --quiet "$SERVICE_NAME" 2>/dev/null; then
        systemctl disable "$SERVICE_NAME"
        print_info "Service disabled"
    fi
    
    # Remove files
    rm -f "$BIN_DIR/$BINARY_NAME"
    rm -f "$SYSTEMD_DIR/$SERVICE_NAME.service"
    systemctl daemon-reload
    
    # Ask about config directory
    if [[ -d "$CONF_DIR" ]]; then
        read -p "Remove configuration directory $CONF_DIR? [y/N]: " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$CONF_DIR"
            print_info "Configuration directory removed"
        fi
    fi
    
    print_success "Uninstallation completed"
}

show_help() {
    echo "Cloudflare DDNS Updater Installation Script"
    echo
    echo "Usage: $0 [OPTION]"
    echo
    echo "Options:"
    echo "  install     Install the application (default)"
    echo "  uninstall   Remove the application"
    echo "  help        Show this help message"
    echo
    echo "Prerequisites:"
    echo "  - Run as root (use sudo)"
    echo "  - Build the binary first: make build"
    echo "  - Ensure cf-ddns.conf.example exists"
}

# Main script
case "${1:-install}" in
    install)
        print_info "Starting installation of $APP_NAME v$VERSION"
        check_root
        check_dependencies
        check_existing_installation
        create_user
        install_binary
        setup_config
        install_service
        post_install
        ;;
    uninstall)
        check_root
        uninstall
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown option: $1"
        show_help
        exit 1
        ;;
esac