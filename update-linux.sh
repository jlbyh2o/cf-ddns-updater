#!/bin/bash

# Cloudflare DDNS Updater - Linux Update Script
# Usage: sudo ./update-linux.sh

set -e

# Configuration
REPO="jlbyh2o/cf-ddns-updater"
BINARY_NAME="cf-ddns-updater"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/cf-ddns-updater"
LOG_DIR="/var/log"
SERVICE_NAME="cf-ddns-updater.service"
SERVICE_DIR="/etc/systemd/system"
BACKUP_DIR="/tmp/cf-ddns-backup-$(date +%Y%m%d-%H%M%S)"

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

# Check if cf-ddns-updater is installed
check_installation() {
    local binary_path="${INSTALL_DIR}/${BINARY_NAME}"
    if [[ ! -f "$binary_path" ]]; then
        log_error "Cloudflare DDNS Updater is not installed at $binary_path"
        log_info "Please run the install script first"
        exit 1
    fi
    
    if ! systemctl list-unit-files | grep -q "$SERVICE_NAME"; then
        log_error "Service $SERVICE_NAME is not installed"
        log_info "Please run the install script first"
        exit 1
    fi
}

# Get current version
get_current_version() {
    local binary_path="${INSTALL_DIR}/${BINARY_NAME}"
    if [[ -f "$binary_path" ]]; then
        "$binary_path" --version 2>/dev/null | head -n1 | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' || echo "unknown"
    else
        echo "not installed"
    fi
}

# Get latest version from GitHub
get_latest_version() {
    local latest_url="https://api.github.com/repos/${REPO}/releases/latest"
    
    if command -v curl >/dev/null 2>&1; then
        curl -s "$latest_url" | grep '"tag_name":' | sed -E 's/.*"tag_name":\s*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$latest_url" | grep '"tag_name":' | sed -E 's/.*"tag_name":\s*"([^"]+)".*/\1/'
    else
        log_error "Neither curl nor wget is available"
        exit 1
    fi
}

# Detect architecture
detect_architecture() {
    local arch
    arch=$(uname -m)
    case $arch in
        x86_64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        armv7l|armv6l) echo "arm" ;;
        *) 
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Download binary
download_binary() {
    local version="$1"
    local arch="$2"
    local binary_name="${BINARY_NAME}-linux-${arch}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${binary_name}"
    local temp_file
    temp_file=$(mktemp)
    
    log_info "Downloading $binary_name from $download_url"
    
    if command -v curl >/dev/null 2>&1; then
        if ! curl -L -o "$temp_file" "$download_url"; then
            log_error "Failed to download binary with curl"
            rm -f "$temp_file"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -O "$temp_file" "$download_url"; then
            log_error "Failed to download binary with wget"
            rm -f "$temp_file"
            exit 1
        fi
    else
        log_error "Neither curl nor wget is available"
        exit 1
    fi
    
    # Verify download
    if [[ ! -s "$temp_file" ]]; then
        log_error "Downloaded file is empty"
        rm -f "$temp_file"
        exit 1
    fi
    
    echo "$temp_file"
}

# Create backup
create_backup() {
    log_info "Creating backup at $BACKUP_DIR"
    mkdir -p "$BACKUP_DIR"
    mkdir -p "$BACKUP_DIR/config"
    
    # Backup binary
    if [[ -f "${INSTALL_DIR}/${BINARY_NAME}" ]]; then
        cp "${INSTALL_DIR}/${BINARY_NAME}" "$BACKUP_DIR/"
        log_info "Binary backed up"
    fi
    
    # Backup configuration files (not the entire directory)
    if [[ -d "$CONFIG_DIR" ]]; then
        cp "$CONFIG_DIR"/* "$BACKUP_DIR/config/" 2>/dev/null || true
        log_info "Configuration backed up"
    fi
    
    # Backup service file
    if [[ -f "/etc/systemd/system/$SERVICE_NAME" ]]; then
        cp "/etc/systemd/system/$SERVICE_NAME" "$BACKUP_DIR/"
        log_info "Service file backed up"
    fi
    
    log_success "Backup created at $BACKUP_DIR"
}

# Stop service
stop_service() {
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_info "Stopping $SERVICE_NAME service"
        systemctl stop "$SERVICE_NAME"
        log_success "Service stopped"
    else
        log_info "Service $SERVICE_NAME is not running"
    fi
}

# Start service
start_service() {
    log_info "Starting $SERVICE_NAME service"
    systemctl start "$SERVICE_NAME"
    
    # Wait a moment and check if service started successfully
    sleep 2
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log_success "Service started successfully"
    else
        log_error "Service failed to start"
        log_info "Check service status with: systemctl status $SERVICE_NAME"
        log_info "Check logs with: journalctl -u $SERVICE_NAME -f"
        exit 1
    fi
}

# Install binary
install_binary() {
    local temp_file="$1"
    local binary_path="${INSTALL_DIR}/${BINARY_NAME}"
    
    log_info "Installing binary to $binary_path"
    
    # Make executable and move to install directory
    chmod +x "$temp_file"
    mv "$temp_file" "$binary_path"
    
    # Set ownership
    chown root:root "$binary_path"
    
    log_success "Binary installed successfully"
}

# Create system user if it doesn't exist
ensure_user_exists() {
    if ! id "cf-ddns" &>/dev/null; then
        log_info "Creating cf-ddns system user (required for service)"
        useradd --system --no-create-home --shell /bin/false cf-ddns
        log_success "Created cf-ddns system user"
        
        # Set ownership of config directory
        if [[ -d "$CONFIG_DIR" ]]; then
            chown -R cf-ddns:cf-ddns "$CONFIG_DIR"
            chmod 750 "$CONFIG_DIR"
            log_info "Updated config directory ownership"
        fi
    fi
}

# Ensure log directory exists with proper permissions
ensure_log_directory() {
    local cf_ddns_log_dir="$LOG_DIR/cf-ddns-updater"
    
    if [[ ! -d "$cf_ddns_log_dir" ]]; then
        log_info "Creating log directory: $cf_ddns_log_dir"
        mkdir -p "$cf_ddns_log_dir"
        chown cf-ddns:cf-ddns "$cf_ddns_log_dir"
        chmod 755 "$cf_ddns_log_dir"
        log_success "Log directory created with proper permissions"
    else
        # Ensure proper ownership and permissions for existing directory
        chown cf-ddns:cf-ddns "$cf_ddns_log_dir"
        chmod 755 "$cf_ddns_log_dir"
        log_info "Log directory permissions updated"
    fi
}

# Update systemd service file
update_service() {
    log_info "Updating systemd service file"
    
    local service_url="https://raw.githubusercontent.com/${REPO}/main/cf-ddns-updater.service"
    local temp_service="/tmp/cf-ddns-updater.service"
    local service_path="${SERVICE_DIR}/${SERVICE_NAME}"
    
    # Download updated service file
    if command -v curl >/dev/null 2>&1; then
        if curl -L -o "$temp_service" "$service_url" 2>/dev/null; then
            log_info "Downloaded updated service file"
        else
            log_warning "Could not download updated service file, keeping existing one"
            return 0
        fi
    elif command -v wget >/dev/null 2>&1; then
        if wget -O "$temp_service" "$service_url" 2>/dev/null; then
            log_info "Downloaded updated service file"
        else
            log_warning "Could not download updated service file, keeping existing one"
            return 0
        fi
    else
        log_warning "Neither curl nor wget available, keeping existing service file"
        return 0
    fi
    
    # Backup existing service file if it exists
    if [[ -f "$service_path" ]]; then
        cp "$service_path" "$BACKUP_DIR/" 2>/dev/null || true
        log_info "Existing service file backed up"
    fi
    
    # Install updated service file
    cp "$temp_service" "$service_path"
    rm -f "$temp_service"
    
    # Reload systemd daemon
    systemctl daemon-reload
    
    log_success "Service file updated and systemd reloaded"
}

# Check for configuration migration needs
check_config_migration() {
    local old_json_config="${CONFIG_DIR}/config.json"
    local new_toml_config="${CONFIG_DIR}/cf-ddns.conf"
    
    # Check if we have old JSON config but no TOML config
    if [[ -f "$old_json_config" && ! -f "$new_toml_config" ]]; then
        log_warning "Found old JSON configuration file: $old_json_config"
        log_warning "Starting from v0.2.0, the application uses TOML configuration format"
        log_warning "Please manually convert your configuration to TOML format:"
        echo "  1. Copy $old_json_config to $new_toml_config"
        echo "  2. Convert JSON syntax to TOML syntax"
        echo "  3. Update systemd service if needed"
        echo "  4. Test the configuration"
        echo
        log_info "Example TOML configuration is available at: ${CONFIG_DIR}/cf-ddns.conf.example"
        return 1
    fi
    
    return 0
}

# Rollback function
rollback() {
    log_error "Update failed. Rolling back..."
    
    if [[ -d "$BACKUP_DIR" ]]; then
        # Stop service
        systemctl stop "$SERVICE_NAME" 2>/dev/null || true
        
        # Restore binary
        if [[ -f "$BACKUP_DIR/$BINARY_NAME" ]]; then
            cp "$BACKUP_DIR/$BINARY_NAME" "${INSTALL_DIR}/"
            log_info "Binary restored from backup"
        fi
        
        # Restore configuration files
        if [[ -d "$BACKUP_DIR/config" ]]; then
            mkdir -p "$CONFIG_DIR"
            cp "$BACKUP_DIR/config"/* "$CONFIG_DIR/" 2>/dev/null || true
            log_info "Configuration restored from backup"
        fi
        
        # Restore service file
        if [[ -f "$BACKUP_DIR/$SERVICE_NAME" ]]; then
            cp "$BACKUP_DIR/$SERVICE_NAME" "${SERVICE_DIR}/"
            systemctl daemon-reload
            log_info "Service file restored from backup"
        fi
        
        # Start service
        systemctl start "$SERVICE_NAME" 2>/dev/null || true
        
        log_info "Rollback completed. Backup preserved at: $BACKUP_DIR"
    fi
    
    exit 1
}

# Main update function
main() {
    log_info "Starting Cloudflare DDNS Updater update process"
    
    # Check prerequisites
    check_root
    check_installation
    
    # Get version information
    local current_version
    local latest_version
    current_version=$(get_current_version)
    latest_version=$(get_latest_version)
    
    if [[ -z "$latest_version" ]]; then
        log_error "Failed to get latest version information"
        exit 1
    fi
    
    log_info "Current version: $current_version"
    log_info "Latest version: $latest_version"
    
    # Check if update is needed
    if [[ "$current_version" == "$latest_version" ]]; then
        log_success "Already running the latest version ($current_version)"
        exit 0
    fi
    
    # Detect architecture
    local arch
    arch=$(detect_architecture)
    log_info "Detected architecture: $arch"
    
    # Create backup
    create_backup
    
    # Set up error handling for rollback
    trap rollback ERR
    
    # Stop service
    stop_service
    
    # Download new binary
    local temp_file
    temp_file=$(download_binary "$latest_version" "$arch")
    
    # Install new binary
    install_binary "$temp_file"
    
    # Ensure cf-ddns user exists (for older installations)
    ensure_user_exists
    
    # Ensure log directory exists with proper permissions
    ensure_log_directory
    
    # Update systemd service file
    update_service
    
    # Check for configuration migration needs
    if ! check_config_migration; then
        log_warning "Configuration migration required. Service not started."
        log_info "Please migrate your configuration and start the service manually:"
        echo "  sudo systemctl start $SERVICE_NAME"
        exit 0
    fi
    
    # Start service
    start_service
    
    # Clear error trap
    trap - ERR
    
    # Final success message
    log_success "Update completed successfully!"
    log_info "Updated from $current_version to $latest_version"
    log_info "Service is running and ready"
    echo
    log_info "Backup preserved at: $BACKUP_DIR"
    log_info "You can remove the backup once you've verified everything works correctly"
    echo
    log_info "Check service status: systemctl status $SERVICE_NAME"
    log_info "View logs: journalctl -u $SERVICE_NAME -f"
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Cloudflare DDNS Updater Update Script"
        echo "Usage: sudo $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --check        Check for updates without installing"
        echo "  --force        Force update even if versions match"
        echo
        echo "This script will:"
        echo "  - Check current and latest versions"
        echo "  - Create a backup of current installation"
        echo "  - Stop the service during update"
        echo "  - Download and install the latest binary"
        echo "  - Check for configuration migration needs"
        echo "  - Restart the service"
        echo "  - Provide rollback on failure"
        exit 0
        ;;
    --check)
        check_root
        check_installation
        current_version=$(get_current_version)
        latest_version=$(get_latest_version)
        echo "Current version: $current_version"
        echo "Latest version: $latest_version"
        if [[ "$current_version" == "$latest_version" ]]; then
            echo "Status: Up to date"
            exit 0
        else
            echo "Status: Update available"
            exit 1
        fi
        ;;
    --force)
        # Override version check
        get_current_version() {
            echo "force-update"
        }
        ;;
esac

# Run main function
main "$@"