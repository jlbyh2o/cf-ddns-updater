---
layout: page
title: "Installation Guide"
permalink: /installation/
---

# 📦 Installation Guide

This comprehensive guide covers all installation methods for Cloudflare DDNS Updater across different platforms.

## Quick Installation (Linux)

### One-Line Installation Script

The fastest way to get started on Linux systems:

```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash
```

**What this script does:**
- Downloads the latest release binary for your architecture
- Installs to `/usr/local/bin/cf-ddns-updater`
- Creates configuration directory at `/etc/cf-ddns/`
- Downloads example configuration file
- Creates systemd service with security hardening
- Sets up proper file permissions
- Creates dedicated `cf-ddns` user

### Update Script

To update to the latest version:

```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/update-linux.sh | sudo bash
```

### Uninstall Script

To completely remove the installation:

```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/uninstall-linux.sh | sudo bash
```

---

## Manual Installation

### Download Pre-built Binaries

1. Go to the [releases page](https://github.com/jlbyh2o/cf-ddns-updater/releases/latest)
2. Download the appropriate binary for your system:

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| Linux | x86_64 (AMD64) | `cf-ddns-updater-linux-amd64` |
| Linux | ARM | `cf-ddns-updater-linux-arm` |
| Linux | ARM64 | `cf-ddns-updater-linux-arm64` |
| Windows | x86_64 (AMD64) | `cf-ddns-updater-windows-amd64.exe` |
| Windows | ARM64 | `cf-ddns-updater-windows-arm64.exe` |

### Linux Manual Installation

```bash
# Download (replace with your architecture)
wget https://github.com/jlbyh2o/cf-ddns-updater/releases/latest/download/cf-ddns-updater-linux-amd64

# Make executable
chmod +x cf-ddns-updater-linux-amd64

# Move to system path (optional)
sudo mv cf-ddns-updater-linux-amd64 /usr/local/bin/cf-ddns-updater

# Create configuration directory
sudo mkdir -p /etc/cf-ddns

# Create systemd user (optional but recommended)
sudo useradd --system --no-create-home --shell /bin/false cf-ddns
```

### Windows Manual Installation

1. Download `cf-ddns-updater-windows-amd64.exe` or `cf-ddns-updater-windows-arm64.exe`
2. Place it in a directory (e.g., `C:\Program Files\cf-ddns-updater\`)
3. Optionally add the directory to your PATH environment variable

---

## Build from Source

### Prerequisites

- **Go 1.23+**: [Download from golang.org](https://golang.org/dl/)
- **Git**: For cloning the repository
- **Make**: (optional) For using the Makefile

### Build Steps

```bash
# Clone the repository
git clone https://github.com/jlbyh2o/cf-ddns-updater.git
cd cf-ddns-updater

# Build for your platform
go build -o cf-ddns-updater

# Or build for all supported platforms
make build-all
```

### Build Scripts

The repository includes build scripts for convenience:

**Linux/macOS:**
```bash
./build.sh
```

**Windows:**
```cmd
build.bat
```

**Make targets:**
```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make clean          # Clean build artifacts
make test           # Run tests
```

---

## Platform-Specific Setup

### Ubuntu/Debian

```bash
# Using the installation script (recommended)
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash

# Or install manually
sudo apt update
wget https://github.com/jlbyh2o/cf-ddns-updater/releases/latest/download/cf-ddns-updater-linux-amd64
sudo mv cf-ddns-updater-linux-amd64 /usr/local/bin/cf-ddns-updater
sudo chmod +x /usr/local/bin/cf-ddns-updater
```

### CentOS/RHEL/Fedora

```bash
# Using the installation script (recommended)
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash

# Or install manually
sudo dnf install wget  # or yum install wget
wget https://github.com/jlbyh2o/cf-ddns-updater/releases/latest/download/cf-ddns-updater-linux-amd64
sudo mv cf-ddns-updater-linux-amd64 /usr/local/bin/cf-ddns-updater
sudo chmod +x /usr/local/bin/cf-ddns-updater
```

### Arch Linux

```bash
# Using the installation script
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash

# Or build from AUR (if available)
# yay -S cf-ddns-updater  # Future AUR package
```

### macOS

```bash
# Download and install manually
wget https://github.com/jlbyh2o/cf-ddns-updater/releases/latest/download/cf-ddns-updater-linux-amd64
chmod +x cf-ddns-updater-linux-amd64
sudo mv cf-ddns-updater-linux-amd64 /usr/local/bin/cf-ddns-updater

# Or build from source
git clone https://github.com/jlbyh2o/cf-ddns-updater.git
cd cf-ddns-updater
go build -o cf-ddns-updater
```

### Windows

#### Manual Installation

1. Download the Windows binary
2. Create a directory: `C:\Program Files\cf-ddns-updater\`
3. Place the executable there
4. Create config directory: `C:\cf-ddns\`

#### Windows Service (using NSSM)

1. Download [NSSM (Non-Sucking Service Manager)](https://nssm.cc/download)
2. Install as service:

```cmd
# Install NSSM and add to PATH
nssm install "Cloudflare DDNS Updater" "C:\Program Files\cf-ddns-updater\cf-ddns-updater.exe"
nssm set "Cloudflare DDNS Updater" Arguments "-config C:\cf-ddns\cf-ddns.conf"
nssm set "Cloudflare DDNS Updater" AppDirectory "C:\Program Files\cf-ddns-updater"
nssm set "Cloudflare DDNS Updater" DisplayName "Cloudflare DDNS Updater"
nssm set "Cloudflare DDNS Updater" Description "Automatic DNS updater for Cloudflare"
nssm set "Cloudflare DDNS Updater" Start SERVICE_AUTO_START

# Start the service
nssm start "Cloudflare DDNS Updater"
```

#### Windows Task Scheduler

Alternative to service installation:

1. Open Task Scheduler
2. Create Basic Task
3. Configure to run at startup and repeat every 5 minutes
4. Set action to start your cf-ddns-updater.exe with appropriate arguments

---

## Docker Installation

### Pre-built Container (Future)

```bash
# Pull the official image
docker pull jlbyh2o/cf-ddns-updater:latest

# Run with configuration
docker run -d \
  --name cf-ddns-updater \
  --restart unless-stopped \
  -v /path/to/config:/config \
  jlbyh2o/cf-ddns-updater:latest
```

### Build Your Own Container

```dockerfile
FROM scratch
COPY cf-ddns-updater-linux-amd64 /cf-ddns-updater
COPY cf-ddns.conf /cf-ddns.conf
ENTRYPOINT ["/cf-ddns-updater", "-config", "/cf-ddns.conf"]
```

```bash
# Build
docker build -t cf-ddns-updater .

# Run
docker run -d --name cf-ddns-updater --restart unless-stopped cf-ddns-updater
```

---

## Post-Installation Setup

### 1. Create Configuration

After installation, create your configuration file:

**Linux (script install):**
```bash
sudo nano /etc/cf-ddns/cf-ddns.conf
```

**Manual install:**
```bash
nano cf-ddns.conf  # In the same directory as binary
```

**Windows:**
```cmd
notepad C:\cf-ddns\cf-ddns.conf
```

### 2. Basic Configuration Template

```toml
[cloudflare]
api_token = "your_cloudflare_api_token_here"

[[domains]]
name = "example.com"
record_types = "both"
ttl = 300
proxied = false

# Optional: continuous mode
interval = 300
verbose = true
```

### 3. Set Proper Permissions (Linux)

```bash
# If using script installation, permissions are set automatically
# For manual installation:
sudo chown cf-ddns:cf-ddns /etc/cf-ddns/cf-ddns.conf
sudo chmod 600 /etc/cf-ddns/cf-ddns.conf
```

### 4. Test Installation

```bash
# Test configuration
cf-ddns-updater -config /path/to/cf-ddns.conf -verbose -once

# Check version
cf-ddns-updater -version
```

---

## Service Management

### Linux (systemd)

```bash
# Enable on boot
sudo systemctl enable cf-ddns-updater

# Start service
sudo systemctl start cf-ddns-updater

# Check status
sudo systemctl status cf-ddns-updater

# View logs
sudo journalctl -u cf-ddns-updater -f

# Restart service
sudo systemctl restart cf-ddns-updater

# Stop service
sudo systemctl stop cf-ddns-updater
```

### Windows (NSSM)

```cmd
# Start service
nssm start "Cloudflare DDNS Updater"

# Stop service
nssm stop "Cloudflare DDNS Updater"

# Restart service
nssm restart "Cloudflare DDNS Updater"

# Remove service
nssm remove "Cloudflare DDNS Updater" confirm
```

---

## Installation Verification

### Check Installation

```bash
# Verify binary location
which cf-ddns-updater
# or
whereis cf-ddns-updater

# Check version
cf-ddns-updater -version

# Verify configuration file exists
ls -la /etc/cf-ddns/cf-ddns.conf  # Linux
dir C:\cf-ddns\cf-ddns.conf       # Windows
```

### Test Functionality

```bash
# Test with verbose output
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once
```

Expected output:
```
Cloudflare DDNS Updater v1.0.0
Starting DNS update process...
Current IPv4 address: 203.0.113.1
Processing domain: example.com
Successfully processed domain: example.com
```

---

## Troubleshooting Installation

### Common Issues

| Issue | Solution |
|-------|----------|
| `Permission denied` | Ensure binary is executable: `chmod +x cf-ddns-updater` |
| `Command not found` | Add binary location to PATH or use full path |
| `Config file not found` | Verify config file path and permissions |
| `Service failed to start` | Check systemd logs: `journalctl -u cf-ddns-updater` |

### Log Locations

- **Linux (systemd)**: `journalctl -u cf-ddns-updater`
- **Linux (manual)**: Specify with `-log /path/to/log.txt`
- **Windows**: Check Windows Event Viewer or specify log file

### Getting Help

- 📖 [Configuration Guide](/cf-ddns-updater/configuration/)
- 🔧 [Troubleshooting Guide](/cf-ddns-updater/troubleshooting/)
- 💬 [GitHub Discussions](https://github.com/jlbyh2o/cf-ddns-updater/discussions)
- 🐛 [Report Issues](https://github.com/jlbyh2o/cf-ddns-updater/issues)

---

## Next Steps

After successful installation:

1. ✅ **Configure your domains** - [Configuration Guide](/cf-ddns-updater/configuration/)
2. ✅ **Set up monitoring** - Monitor logs and set up alerts
3. ✅ **Review security** - [Security Guide](/cf-ddns-updater/security/)
4. ✅ **Test failover** - Ensure backup IP services work