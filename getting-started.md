---
layout: page
title: "Getting Started"
permalink: /getting-started/
---

# 🏁 Getting Started

Welcome to Cloudflare DDNS Updater! This guide will help you get up and running in just a few minutes.

## What is Dynamic DNS?

Dynamic DNS (DDNS) automatically updates your domain's DNS records when your IP address changes. This is essential for:

- **Home servers** with changing ISP-assigned IP addresses
- **Self-hosted services** that need consistent domain access
- **Remote access** to your network and devices
- **Development environments** with dynamic hosting

## Prerequisites

Before you begin, you'll need:

1. **A Cloudflare account** with your domain configured
2. **A Cloudflare API token** with appropriate permissions
3. **A supported operating system** (Linux, Windows, macOS)

## Step 1: Get Your Cloudflare API Token

<div class="alert alert-info">
<strong>💡 Tip:</strong> API tokens are more secure than global API keys and provide granular permissions.
</div>

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/profile/api-tokens)
2. Click "Create Token"
3. Use the "Custom token" template
4. Configure permissions:
   - **Zone Resources**: Include specific zones or all zones
   - **Permissions**: 
     - Zone: Zone:Read
     - Zone: DNS:Edit
5. Copy the generated token (you'll need it for configuration)

## Step 2: Quick Installation

### Linux (Recommended)

The fastest way to get started on Linux:

```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash
```

This script will:
- Download the latest binary
- Install it to `/usr/local/bin/`
- Create a systemd service
- Set up configuration directory at `/etc/cf-ddns/`
- Configure proper permissions

### Manual Installation

If you prefer manual installation or are using Windows:

1. [Download the latest release](https://github.com/jlbyh2o/cf-ddns-updater/releases/latest)
2. Choose the appropriate binary for your system:
   - Linux x64: `cf-ddns-updater-linux-amd64`
   - Linux ARM: `cf-ddns-updater-linux-arm`  
   - Linux ARM64: `cf-ddns-updater-linux-arm64`
   - Windows x64: `cf-ddns-updater-windows-amd64.exe`
   - Windows ARM64: `cf-ddns-updater-windows-arm64.exe`
3. Make it executable (Linux/macOS): `chmod +x cf-ddns-updater-*`
4. Move to a directory in your PATH

## Step 3: Basic Configuration

Create a configuration file. The location depends on your installation method:

- **Linux (script install)**: `/etc/cf-ddns/cf-ddns.conf`
- **Manual install**: `cf-ddns.conf` in the same directory as the binary

### Minimal Configuration

```toml
[cloudflare]
api_token = "your_cloudflare_api_token_here"

[[domains]]
name = "example.com"
record_types = "both"  # "A", "AAAA", or "both"
ttl = 300
proxied = false
```

### Configuration Explained

- **`api_token`**: Your Cloudflare API token from Step 1
- **`name`**: Your domain or subdomain (e.g., "example.com", "home.example.com")
- **`record_types`**: Which DNS records to update:
  - `"A"`: IPv4 only
  - `"AAAA"`: IPv6 only  
  - `"both"`: Both IPv4 and IPv6 (recommended)
- **`ttl`**: Time-to-live in seconds (300 = 5 minutes, good for dynamic IPs)
- **`proxied`**: Whether to proxy traffic through Cloudflare (false for direct DNS)

## Step 4: Test Your Configuration

Before setting up automatic updates, test your configuration:

```bash
# Linux (script install)
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once

# Manual install
./cf-ddns-updater -config cf-ddns.conf -verbose -once
```

You should see output like:

```
Starting DNS update process...
Current IPv4 address: 203.0.113.1
Processing domain: example.com
Updating A record for example.com: 198.51.100.1 to 203.0.113.1
Successfully updated A record for example.com
```

<div class="alert alert-success">
<strong>✅ Success!</strong> If you see "Successfully updated" messages, your configuration is working correctly.
</div>

<div class="alert alert-danger">
<strong>❌ Errors?</strong> Check our <a href="/cf-ddns-updater/troubleshooting/">troubleshooting guide</a> for common issues and solutions.
</div>

## Step 5: Set Up Automatic Updates

### Linux (systemd service)

If you used the installation script, the systemd service is already configured:

```bash
# Enable and start the service
sudo systemctl enable cf-ddns-updater
sudo systemctl start cf-ddns-updater

# Check status
sudo systemctl status cf-ddns-updater

# View logs
sudo journalctl -u cf-ddns-updater -f
```

### Manual Setup (Linux)

Create a systemd service file at `/etc/systemd/system/cf-ddns-updater.service`:

```ini
[Unit]
Description=Cloudflare DDNS Updater
After=network.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/path/to/cf-ddns-updater -config /path/to/cf-ddns.conf
Restart=always
RestartSec=30
User=nobody
Group=nobody

[Install]
WantedBy=multi-user.target
```

Then enable and start it:

```bash
sudo systemctl daemon-reload
sudo systemctl enable cf-ddns-updater
sudo systemctl start cf-ddns-updater
```

### Windows

For Windows, you can use Task Scheduler or install as a service using [NSSM](https://nssm.cc/):

```cmd
# Download NSSM and add it to your PATH
nssm install "Cloudflare DDNS Updater" "C:\path\to\cf-ddns-updater.exe"
nssm set "Cloudflare DDNS Updater" Arguments "-config C:\path\to\cf-ddns.conf"
nssm start "Cloudflare DDNS Updater"
```

## Step 6: Configure Update Interval

By default, the updater runs once and exits. For continuous monitoring, add an interval to your configuration:

```toml
[cloudflare]
api_token = "your_token_here"

[[domains]]
name = "example.com"
record_types = "both"
ttl = 300
proxied = false

# Update every 5 minutes (300 seconds)
interval = 300
verbose = true  # Enable detailed logging
```

## Next Steps

🎉 **Congratulations!** Your DDNS updater is now running. Here's what to do next:

1. **Monitor the logs** to ensure it's working correctly
2. **Configure multiple domains** if needed
3. **Set up monitoring** to get notified of issues
4. **Review security settings** in our [security guide](/cf-ddns-updater/security/)

## Quick Reference

### Common Commands

```bash
# Check status (Linux systemd)
sudo systemctl status cf-ddns-updater

# View logs (Linux systemd)  
sudo journalctl -u cf-ddns-updater -f

# Test configuration
cf-ddns-updater -config /path/to/config -verbose -once

# Run manually with interval
cf-ddns-updater -config /path/to/config -verbose
```

### Configuration File Locations

- **Linux (script install)**: `/etc/cf-ddns/cf-ddns.conf`
- **Linux (manual)**: `./cf-ddns.conf` or `/usr/local/etc/cf-ddns.conf`
- **Windows**: `cf-ddns.conf` in the same directory as the executable

### Log Locations

- **Linux (systemd)**: `journalctl -u cf-ddns-updater`
- **Manual/Windows**: Configure with `-log /path/to/logfile.log`

---

## Need Help?

- 📖 [Full Configuration Guide](/cf-ddns-updater/configuration/)
- 🔧 [Troubleshooting](/cf-ddns-updater/troubleshooting/)
- 🔒 [Security Best Practices](/cf-ddns-updater/security/)
- 💬 [Community Support](https://github.com/jlbyh2o/cf-ddns-updater/discussions)
- 🐛 [Report Issues](https://github.com/jlbyh2o/cf-ddns-updater/issues)