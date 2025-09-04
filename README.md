# Cloudflare DDNS Updater

🚀 **A reliable and lightweight Dynamic DNS updater for Cloudflare** - Keep your domains pointing to your dynamic IP address automatically!

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/D1D51K3UOB)
[![Build Status](https://github.com/jlbyh2o/cf-ddns-updater/actions/workflows/build.yml/badge.svg)](https://github.com/jlbyh2o/cf-ddns-updater/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/jlbyh2o/cf-ddns-updater)](https://github.com/jlbyh2o/cf-ddns-updater/blob/main/go.mod)
[![Release](https://img.shields.io/github/v/release/jlbyh2o/cf-ddns-updater)](https://github.com/jlbyh2o/cf-ddns-updater/releases)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/jlbyh2o/cf-ddns-updater)](https://goreportcard.com/report/github.com/jlbyh2o/cf-ddns-updater)
[![Downloads](https://img.shields.io/github/downloads/jlbyh2o/cf-ddns-updater/total)](https://github.com/jlbyh2o/cf-ddns-updater/releases)

## ⚡ Quick Install (Linux)

**Install with one command:**
```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash
```

**Update with one command:**
```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/update-linux.sh | sudo bash
```

**Uninstall with one command:**
```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/uninstall-linux.sh | sudo bash
```

> 💡 **That's it!** The installer automatically:
> - Downloads and installs the latest binary
> - Creates a systemd service for automatic startup
> - Sets up proper user permissions and security
> - Creates configuration directory at `/etc/cf-ddns/`

## ✨ Features

- 🔄 **Automatic Updates**: Keeps your DNS records in sync with your dynamic IP
- 🛡️ **Reliable**: Uses multiple IP detection services for redundancy
- ⚙️ **Flexible**: Supports IPv4 (A records), IPv6 (AAAA records), or both
- 🔐 **Secure**: API token authentication with minimal required permissions
- 🐧 **Linux Ready**: Complete systemd integration with security hardening
- 📊 **Logging**: Comprehensive logging with systemd journal integration
- 🚀 **Lightweight**: Single binary with no dependencies

## 🌐 IP Detection with fetch-ip.com

This project uses **[fetch-ip.com](https://fetch-ip.com)** as the primary IP detection service, with additional fallback services for maximum reliability.

### Why fetch-ip.com?

- 🚀 **Fast & Reliable**: Optimized for speed and uptime
- 🔒 **Privacy-Focused**: No logging, no tracking, no data collection
- 🌍 **Global Infrastructure**: Multiple endpoints for worldwide accessibility
- ⚡ **Lightweight**: Minimal response overhead for quick detection
- 🛡️ **Redundant**: Multiple fallback services ensure continuous operation

### IP Detection Flow

1. **Primary**: `https://v4.fetch-ip.com` (IPv4) / `https://v6.fetch-ip.com` (IPv6)
2. **Fallback Services**: icanhazip.com, ipify.org, ident.me
3. **Validation**: Each detected IP is validated before use
4. **Error Handling**: Automatic failover if any service is unavailable

> 💡 **Note**: fetch-ip.com is maintained by the same team behind this DDNS updater, ensuring optimal compatibility and performance.
> 
> 📖 **Learn more**: Check out the [fetch-ip.com documentation](https://fetch-ip.com/docs) for detailed API information and usage examples.

## 🚀 Getting Started

### 1. Install (Choose Your Method)

#### Option A: One-Line Install (Recommended)
```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash
```

#### Option B: Manual Installation
1. Download the latest release from [GitHub Releases](https://github.com/jlbyh2o/cf-ddns-updater/releases)
2. Extract and run the binary manually

### 2. Configure Your Cloudflare API

1. Get your Cloudflare API token:
   - Go to [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens)
   - Create a custom token with permissions: `Zone:DNS:Edit`, `Zone:Zone:Read`

2. Edit the configuration:
   ```bash
   sudo nano /etc/cf-ddns/cf-ddns.conf
   ```

3. Add your configuration:
   ```toml
   [cloudflare]
   api_token = "your_cloudflare_api_token_here"
   
   [[domains]]
   name = "example.com"
   record_types = "both"  # "A", "AAAA", or "both"
   ttl = 300
   proxied = false
   ```

### 3. Start the Service

```bash
# Enable and start the service
sudo systemctl enable cf-ddns-updater
sudo systemctl start cf-ddns-updater

# Check status
sudo systemctl status cf-ddns-updater

# View logs
sudo journalctl -u cf-ddns-updater -f
```

## 📋 Configuration Reference

### Configuration File Location
After installation, edit: `/etc/cf-ddns/cf-ddns.conf`

### Basic Configuration
```toml
[cloudflare]
api_token = "your_cloudflare_api_token_here"

[[domains]]
name = "example.com"
record_types = "both"  # "A", "AAAA", or "both"
ttl = 300
proxied = false

# Optional: Run continuously (in seconds)
interval = 300  # Update every 5 minutes
verbose = true  # Enable detailed logging
```

### Multiple Domains Example
```toml
[cloudflare]
api_token = "your_token_here"

[[domains]]
name = "example.com"
record_types = "both"
ttl = 300
proxied = false

[[domains]]
name = "www.example.com"
record_types = "A"
ttl = 300
proxied = true

[[domains]]
name = "ipv6.example.com"
record_types = "AAAA"
ttl = 300
proxied = false

interval = 300
verbose = true
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|----------|
| `api_token` | Cloudflare API token (recommended) | Required |
| `api_key` + `email` | Legacy authentication method | Alternative |
| `name` | Domain or subdomain name | Required |
| `record_types` | "A", "AAAA", or "both" | "both" |
| `ttl` | DNS record TTL in seconds | 300 |
| `proxied` | Proxy through Cloudflare | false |
| `interval` | Update interval in seconds (0 = run once) | 0 |
| `verbose` | Enable verbose logging | false |

## 🔧 Management Commands

### Service Management
```bash
# Check service status
sudo systemctl status cf-ddns-updater

# View live logs
sudo journalctl -u cf-ddns-updater -f

# Restart service
sudo systemctl restart cf-ddns-updater

# Stop service
sudo systemctl stop cf-ddns-updater

# Disable service
sudo systemctl disable cf-ddns-updater
```

### Manual Execution
```bash
# Run once manually
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -once -verbose

# Test configuration
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -once
```

### Command Line Options
| Option | Description |
|--------|-------------|
| `-config` | Path to configuration file |
| `-verbose` | Enable verbose logging |
| `-log` | Log file path (optional) |
| `-once` | Run once and exit |

## 🛠️ Building from Source

### Quick Build
```bash
git clone https://github.com/jlbyh2o/cf-ddns-updater.git
cd cf-ddns-updater

# Linux/macOS
./build.sh

# Windows
build.bat

# Or use Make (Linux)
make build-all
```

**Prerequisites:** Go 1.23+, Git

**Output:** Binaries in `bin/` directory for all supported platforms

## 🔒 Security Features

The systemd service includes comprehensive security hardening:
- ✅ Dedicated `cf-ddns` user with minimal privileges
- ✅ Private temporary directories
- ✅ Protected system directories
- ✅ Restricted system calls
- ✅ Network access limited to IPv4/IPv6 only
- ✅ Memory execution protection
- ✅ No new privileges allowed

## 📁 Installation Paths

| Component | Path |
|-----------|------|
| Binary | `/usr/local/bin/cf-ddns-updater` |
| Configuration | `/etc/cf-ddns/cf-ddns.conf` |
| Service File | `/etc/systemd/system/cf-ddns-updater.service` |
| Logs | `/var/log/cf-ddns-updater/` |
| System User | `cf-ddns` (auto-created) |

## 🐛 Troubleshooting

### Quick Diagnostics
```bash
# Check service status
sudo systemctl status cf-ddns-updater

# View recent logs
sudo journalctl -u cf-ddns-updater --since "1 hour ago"

# Test configuration manually
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once
```

### Common Issues

| Issue | Solution |
|-------|----------|
| **Authentication Errors** | Verify API token permissions: `Zone:DNS:Edit`, `Zone:Zone:Read` |
| **IP Detection Failures** | Check internet connectivity; tool uses multiple services for redundancy |
| **DNS Update Failures** | Verify domain ownership in Cloudflare and exact record name match |
| **Service Won't Start** | Check configuration file syntax and permissions |

### Debug Commands
```bash
# Test with verbose output
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once

# Check configuration file
sudo cat /etc/cf-ddns/cf-ddns.conf

# Verify service file
sudo systemctl cat cf-ddns-updater
```

## 🔐 Security Best Practices

- ✅ Use API tokens (not global API keys)
- ✅ Limit token permissions to minimum required
- ✅ Secure configuration file permissions: `sudo chmod 600 /etc/cf-ddns/cf-ddns.conf`
- ✅ Regularly rotate API credentials
- ✅ Monitor logs for suspicious activity

## 🐳 Alternative Deployments

### Docker
```dockerfile
FROM scratch
COPY cf-ddns-updater-linux-amd64 /cf-ddns-updater
COPY cf-ddns.conf /cf-ddns.conf
ENTRYPOINT ["/cf-ddns-updater", "-config", "/cf-ddns.conf"]
```

### Windows Service
Use [NSSM](https://nssm.cc/) for Windows service installation:
```cmd
nssm install "Cloudflare DDNS Updater" "C:\path\to\cf-ddns-updater.exe"
nssm set "Cloudflare DDNS Updater" Arguments "-config C:\path\to\cf-ddns.conf"
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

- 🐛 **Report bugs** - Open an issue with detailed information
- 💡 **Suggest features** - Share your ideas for improvements
- 🔧 **Submit PRs** - Fix bugs or add new features
- 📖 **Improve docs** - Help make the documentation better

## 📄 License

This project is licensed under the **GNU General Public License v3.0 (GPL-3.0)**.

See the [LICENSE](LICENSE) file for full details.

## 🆘 Support

**Need help?** Here's how to get support:

1. 📖 Check the [troubleshooting section](#-troubleshooting) above
2. 🔍 Search [existing issues](https://github.com/jlbyh2o/cf-ddns-updater/issues)
3. 🆕 [Open a new issue](https://github.com/jlbyh2o/cf-ddns-updater/issues/new) with:
   - Your configuration (remove sensitive data)
   - Log output with `-verbose` flag
   - System information (OS, architecture)
   - Steps to reproduce the issue

## ☕ Support Development

This project is developed and maintained in my free time. If Cloudflare DDNS Updater has helped you or saved you time, consider supporting continued development:

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/D1D51K3UOB)

**Your support helps with:**
- 🚀 **New Features** - Adding requested functionality and improvements
- 🔧 **Bug Fixes** - Faster response to issues and maintenance  
- 📖 **Documentation** - Keeping guides comprehensive and up-to-date
- 🛡️ **Security** - Regular security reviews and updates
- 🌍 **Infrastructure** - Testing across different platforms and environments

---

⭐ **Found this helpful?** Give us a star on GitHub!

🔄 **Keep your DNS records in sync automatically with Cloudflare DDNS Updater!**