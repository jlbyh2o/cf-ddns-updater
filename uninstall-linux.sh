#!/bin/bash

# Cloudflare DDNS Updater - Linux Uninstall Script
# Usage: sudo ./uninstall-linux.sh

set -e

# Configuration
BINARY_NAME="cf-ddns-updater"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/cf-ddns-updater"
SERVICE_DIR="/etc/systemd/system"
SERVICE_NAME="cf-ddns-updater.service"
USER_NAME="cf-ddns"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" >&2
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" >&2
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" >&2
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Stop and disable systemd service
stop_service() {
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_info "Stopping $SERVICE_NAME service"
        systemctl stop "$SERVICE_NAME"
        log_success "Service stopped"
    else
        log_info "Service $SERVICE_NAME is not running"
    fi
    
    if systemctl is-enabled --quiet "$SERVICE_NAME" 2>/dev/null; then
        log_info "Disabling $SERVICE_NAME service"
        systemctl disable "$SERVICE_NAME"
        log_success "Service disabled"
    else
        log_info "Service $SERVICE_NAME is not enabled"
    fi
}

# Remove systemd service file
remove_service() {
    local service_file="${SERVICE_DIR}/${SERVICE_NAME}"
    if [[ -f "$service_file" ]]; then
        log_info "Removing systemd service file: $service_file"
        rm -f "$service_file"
        systemctl daemon-reload
        log_success "Service file removed"
    else
        log_info "Service file not found: $service_file"
    fi
}

# Remove binary
remove_binary() {
    local binary_path="${INSTALL_DIR}/${BINARY_NAME}"
    if [[ -f "$binary_path" ]]; then
        log_info "Removing binary: $binary_path"
        rm -f "$binary_path"
        log_success "Binary removed"
    else
        log_info "Binary not found: $binary_path"
    fi
}

# Remove configuration directory
remove_config() {
    if [[ -d "$CONFIG_DIR" ]]; then
        log_warning "Removing configuration directory: $CONFIG_DIR"
        log_warning "This will delete all configuration files!"
        read -p "Are you sure you want to continue? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$CONFIG_DIR"
            log_success "Configuration directory removed"
        else
            log_info "Configuration directory preserved"
        fi
    else
        log_info "Configuration directory not found: $CONFIG_DIR"
    fi
}

# Remove user account
remove_user() {
    if id "$USER_NAME" &>/dev/null; then
        log_info "Removing user account: $USER_NAME"
        userdel "$USER_NAME" 2>/dev/null || log_warning "Failed to remove user $USER_NAME"
        log_success "User account removed"
    else
        log_info "User account not found: $USER_NAME"
    fi
}

# Clean up log files
clean_logs() {
    log_info "Cleaning up systemd journal logs for $SERVICE_NAME"
    journalctl --vacuum-time=1s --unit="$SERVICE_NAME" >/dev/null 2>&1 || true
    log_success "Log cleanup completed"
}

# Main uninstall function
main() {
    log_info "Starting Cloudflare DDNS Updater uninstallation"
    
    # Check prerequisites
    check_root
    
    # Stop and disable service
    stop_service
    
    # Remove service file
    remove_service
    
    # Remove binary
    remove_binary
    
    # Remove configuration (with confirmation)
    remove_config
    
    # Remove user account
    remove_user
    
    # Clean up logs
    clean_logs
    
    # Final success message
    log_success "Cloudflare DDNS Updater uninstalled successfully!"
    echo
    log_info "The following items have been removed:"
    echo "  - Binary: ${INSTALL_DIR}/${BINARY_NAME}"
    echo "  - Service: ${SERVICE_DIR}/${SERVICE_NAME}"
    echo "  - User: $USER_NAME"
    echo "  - Logs: systemd journal entries"
    echo
    if [[ -d "$CONFIG_DIR" ]]; then
        log_info "Configuration directory preserved: $CONFIG_DIR"
    else
        log_info "Configuration directory removed: $CONFIG_DIR"
    fi
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Cloudflare DDNS Updater Uninstall Script"
        echo "Usage: sudo $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h    Show this help message"
        echo "  --force       Skip confirmation prompts"
        echo
        echo "This script will remove:"
        echo "  - Binary from $INSTALL_DIR"
        echo "  - Systemd service from $SERVICE_DIR"
        echo "  - User account: $USER_NAME"
        echo "  - Configuration directory: $CONFIG_DIR (with confirmation)"
        echo "  - Systemd journal logs"
        exit 0
        ;;
    --force)
        # Override remove_config function to skip confirmation
        remove_config() {
            if [[ -d "$CONFIG_DIR" ]]; then
                log_info "Removing configuration directory: $CONFIG_DIR"
                rm -rf "$CONFIG_DIR"
                log_success "Configuration directory removed"
            else
                log_info "Configuration directory not found: $CONFIG_DIR"
            fi
        }
        ;;
esac

# Run main function
main "$@"