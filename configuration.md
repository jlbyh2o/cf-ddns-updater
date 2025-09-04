---
layout: page
title: "Configuration Reference"
permalink: /configuration/
---

# ⚙️ Configuration Reference

Complete configuration guide for Cloudflare DDNS Updater with examples and best practices.

## Configuration File Format

Cloudflare DDNS Updater uses TOML (Tom's Obvious Minimal Language) format for configuration. This format is human-readable and easy to edit.

## Basic Structure

```toml
# Cloudflare API configuration
[cloudflare]
api_token = "your_token_here"

# Domain configurations (can have multiple)
[[domains]]
name = "example.com"
record_types = "both"
ttl = 300
proxied = false

# Optional global settings
interval = 300
verbose = true
```

---

## Cloudflare Section

### API Authentication

#### API Token (Recommended)

```toml
[cloudflare]
api_token = "your_cloudflare_api_token_here"
```

**Required Permissions:**
- `Zone:Zone:Read` - To read zone information
- `Zone:DNS:Edit` - To update DNS records

#### Legacy API Key (Not Recommended)

```toml
[cloudflare]
api_key = "your_global_api_key"
email = "your_cloudflare_email@example.com"
```

<div class="alert alert-warning">
<strong>⚠️ Security Note:</strong> API tokens are more secure than global API keys. Use tokens when possible as they provide granular permissions and can be easily revoked.
</div>

---

## Domain Configuration

### Basic Domain Settings

```toml
[[domains]]
name = "example.com"           # Domain or subdomain name
record_types = "both"          # "A", "AAAA", or "both"
ttl = 300                      # Time-to-live in seconds
proxied = false                # Cloudflare proxy (orange cloud)
```

### Record Types

| Value | Description | Use Case |
|-------|-------------|----------|
| `"A"` | IPv4 only | IPv4-only networks or services |
| `"AAAA"` | IPv6 only | IPv6-only networks |
| `"both"` | Both IPv4 and IPv6 | Most common - dual-stack networks |

### TTL (Time-to-Live)

| Value | Description | Recommended For |
|-------|-------------|-----------------|
| `60` | 1 minute | Very dynamic IPs, testing |
| `300` | 5 minutes | **Recommended for home/dynamic IPs** |
| `1800` | 30 minutes | Semi-stable IPs |
| `3600` | 1 hour | Stable IPs |

### Proxied Setting

```toml
proxied = false    # Direct DNS - IP visible
proxied = true     # Proxied through Cloudflare - IP hidden
```

<div class="alert alert-info">
<strong>💡 Tip:</strong> Use <code>proxied = false</code> for services that need direct IP access (SSH, VPN, email). Use <code>proxied = true</code> for web services to benefit from Cloudflare's CDN and DDoS protection.
</div>

---

## Global Settings

### Update Interval

```toml
interval = 300    # Check for IP changes every 5 minutes
interval = 0      # Run once and exit (default)
```

### Logging

```toml
verbose = true    # Enable detailed logging
verbose = false   # Minimal logging (default)
```

---

## Complete Configuration Examples

### Single Domain (Basic)

```toml
[cloudflare]
api_token = "your_token_here"

[[domains]]
name = "home.example.com"
record_types = "both"
ttl = 300
proxied = false

interval = 300
verbose = true
```

### Multiple Domains

```toml
[cloudflare]
api_token = "your_token_here"

# Main domain - web services (proxied)
[[domains]]
name = "example.com"
record_types = "A"
ttl = 300
proxied = true

# WWW subdomain - web services (proxied)
[[domains]]
name = "www.example.com"
record_types = "A"
ttl = 300
proxied = true

# Home server - direct access (not proxied)
[[domains]]
name = "home.example.com"
record_types = "both"
ttl = 300
proxied = false

# VPN server - IPv4 only (not proxied)
[[domains]]
name = "vpn.example.com"
record_types = "A"
ttl = 300
proxied = false

# IPv6-only service
[[domains]]
name = "ipv6.example.com"
record_types = "AAAA"
ttl = 300
proxied = false

interval = 300
verbose = true
```

### Development Environment

```toml
[cloudflare]
api_token = "your_token_here"

[[domains]]
name = "dev.example.com"
record_types = "A"
ttl = 60          # Short TTL for frequent changes
proxied = false   # Direct access for development

[[domains]]
name = "staging.example.com"
record_types = "A"
ttl = 300
proxied = true    # Use Cloudflare features for staging

interval = 60     # Check every minute for dev environments
verbose = true
```

### Home Lab Setup

```toml
[cloudflare]
api_token = "your_token_here"

# Main home server
[[domains]]
name = "home.example.com"
record_types = "both"
ttl = 300
proxied = false

# Services
[[domains]]
name = "nas.example.com"
record_types = "both"
ttl = 300
proxied = false

[[domains]]
name = "plex.example.com"
record_types = "A"
ttl = 300
proxied = false

[[domains]]
name = "homeassistant.example.com"
record_types = "A"
ttl = 300
proxied = false

# Wildcard alternative (configure *.home.example.com in Cloudflare dashboard)
[[domains]]
name = "wildcard.example.com"
record_types = "both"
ttl = 300
proxied = false

interval = 300
verbose = false   # Less verbose for production
```

---

## Configuration File Locations

### Linux (Script Installation)

```bash
/etc/cf-ddns/cf-ddns.conf
```

**Ownership and Permissions:**
```bash
# Check current permissions
ls -la /etc/cf-ddns/cf-ddns.conf

# Set secure permissions (done automatically by install script)
sudo chown cf-ddns:cf-ddns /etc/cf-ddns/cf-ddns.conf
sudo chmod 600 /etc/cf-ddns/cf-ddns.conf
```

### Linux (Manual Installation)

Common locations:
- `./cf-ddns.conf` (same directory as binary)
- `/usr/local/etc/cf-ddns.conf`
- `/home/user/.config/cf-ddns/cf-ddns.conf`

### Windows

Recommended location:
```
C:\cf-ddns\cf-ddns.conf
```

### Docker

Mount configuration as volume:
```bash
docker run -v /path/to/cf-ddns.conf:/config/cf-ddns.conf cf-ddns-updater
```

---

## Command Line Options

Override configuration file settings with command line flags:

```bash
cf-ddns-updater [options]

Options:
  -config string    Path to configuration file (default "cf-ddns.conf")
  -verbose          Enable verbose logging
  -log string       Log file path (optional, logs to stdout if not specified)
  -once            Run once and exit (ignore interval setting)
  -version         Show version information and exit

Examples:
  # Use custom config file
  cf-ddns-updater -config /path/to/my-config.conf
  
  # Run once with verbose logging
  cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once
  
  # Log to file
  cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -log /var/log/cf-ddns.log
```

---

## Environment Variables

You can use environment variables for sensitive information:

```bash
# Set environment variable
export CF_API_TOKEN="your_token_here"
```

```toml
# Reference in configuration
[cloudflare]
api_token = "${CF_API_TOKEN}"
```

<div class="alert alert-info">
<strong>💡 Note:</strong> Environment variable substitution is not currently implemented but is planned for a future release. For now, use secure file permissions to protect your configuration.
</div>

---

## Configuration Validation

### Test Your Configuration

```bash
# Dry run - test configuration without making changes
cf-ddns-updater -config /path/to/cf-ddns.conf -verbose -once

# Check for syntax errors
cf-ddns-updater -config /path/to/cf-ddns.conf -version
```

### Common Validation Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `TOML parse error` | Invalid TOML syntax | Check brackets, quotes, and formatting |
| `Missing api_token` | No Cloudflare credentials | Add `api_token` or `api_key` + `email` |
| `Invalid domain name` | Malformed domain | Check domain name format |
| `Invalid record_types` | Wrong record type value | Use "A", "AAAA", or "both" |
| `Invalid TTL` | TTL out of range | Use value between 60-2147483647 |

---

## Security Best Practices

### Configuration File Security

```bash
# Secure file permissions (Linux)
sudo chmod 600 /etc/cf-ddns/cf-ddns.conf
sudo chown cf-ddns:cf-ddns /etc/cf-ddns/cf-ddns.conf

# Windows - Set file permissions to restrict access
icacls C:\cf-ddns\cf-ddns.conf /inheritance:d
icacls C:\cf-ddns\cf-ddns.conf /grant:r "Users:(R)"
```

### API Token Security

1. **Use API tokens instead of API keys**
2. **Limit token permissions** to only what's needed
3. **Set token expiration** if supported
4. **Regularly rotate tokens**
5. **Never commit tokens to version control**

### Network Security

```toml
# Use specific domains instead of wildcards when possible
[[domains]]
name = "specific.example.com"  # Good
# name = "*.example.com"       # Avoid if possible
```

---

## Advanced Configuration

### Multiple Cloudflare Accounts

Currently, each configuration file supports one Cloudflare account. For multiple accounts, use separate configuration files:

```bash
# Account 1
cf-ddns-updater -config /etc/cf-ddns/account1.conf &

# Account 2  
cf-ddns-updater -config /etc/cf-ddns/account2.conf &
```

### Zone-Specific Configurations

If you have domains across multiple zones:

```toml
[cloudflare]
api_token = "token_with_access_to_all_zones"

# Zone 1: example.com
[[domains]]
name = "home.example.com"
record_types = "both"
ttl = 300
proxied = false

# Zone 2: mydomain.net
[[domains]]
name = "server.mydomain.net"
record_types = "A"
ttl = 300
proxied = false
```

### Conditional Updates

Currently, the updater always checks for IP changes. Future versions may support:
- Time-based conditions
- IP change thresholds
- External trigger support

---

## Migration from Other DDNS Clients

### From ddclient

Common ddclient configuration elements and their cf-ddns-updater equivalents:

| ddclient | cf-ddns-updater |
|----------|-----------------|
| `protocol=cloudflare` | Built-in Cloudflare support |
| `login=email@example.com` | `email = "email@example.com"` |
| `password=api_key` | `api_key = "key"` or `api_token = "token"` |
| `zone=example.com` | Automatic zone detection |
| `example.com` | `name = "example.com"` |

### From Cloudflare's ddns.sh

If migrating from Cloudflare's official ddns.sh script:

1. Extract your API token from the script
2. Identify your domains and record types
3. Convert to TOML format using examples above

---

## Troubleshooting Configuration

### Debug Configuration Issues

1. **Enable verbose logging:**
   ```bash
   cf-ddns-updater -config /path/to/config -verbose -once
   ```

2. **Check TOML syntax:**
   Use an online TOML validator or:
   ```bash
   # Python method
   python3 -c "import tomli; print(tomli.load(open('cf-ddns.conf', 'rb')))"
   ```

3. **Validate API credentials:**
   Test with curl:
   ```bash
   curl -X GET "https://api.cloudflare.com/client/v4/user/tokens/verify" \
        -H "Authorization: Bearer your_token_here" \
        -H "Content-Type: application/json"
   ```

### Common Issues

- **File not found**: Check config file path and permissions
- **Permission denied**: Ensure readable permissions for the user running the service
- **Invalid token**: Verify token permissions and expiration
- **DNS update failed**: Check domain ownership and zone configuration

---

## Next Steps

After configuring your DDNS updater:

1. 📋 **Test thoroughly** in verbose mode
2. 🔒 **Review security settings** - [Security Guide](/cf-ddns-updater/security/)
3. 🔧 **Set up monitoring** - [Troubleshooting Guide](/cf-ddns-updater/troubleshooting/)
4. 🤝 **Join the community** - [GitHub Discussions](https://github.com/jlbyh2o/cf-ddns-updater/discussions)