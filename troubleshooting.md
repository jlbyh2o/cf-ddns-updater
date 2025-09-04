---
layout: page
title: "Troubleshooting Guide"
permalink: /troubleshooting/
---

# 🔧 Troubleshooting Guide

Comprehensive troubleshooting guide for common issues and their solutions.

<div class="toc">
<h3>Quick Navigation</h3>
<ul>
  <li><a href="#quick-diagnostics">Quick Diagnostics</a></li>
  <li><a href="#common-errors">Common Errors</a></li>
  <li><a href="#authentication-issues">Authentication Issues</a></li>
  <li><a href="#network-problems">Network Problems</a></li>
  <li><a href="#service-issues">Service Issues</a></li>
  <li><a href="#logging-debugging">Logging & Debugging</a></li>
  <li><a href="#performance-issues">Performance Issues</a></li>
</ul>
</div>

## Quick Diagnostics

### Health Check Commands

Run these commands to quickly diagnose your installation:

```bash
# Check if binary exists and is executable
which cf-ddns-updater
cf-ddns-updater -version

# Test configuration syntax
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once

# Check service status (Linux)
sudo systemctl status cf-ddns-updater

# View recent logs (Linux)
sudo journalctl -u cf-ddns-updater --since "1 hour ago"
```

### Quick Fix Checklist

Before diving deep, try these common fixes:

- ✅ **Restart the service**: `sudo systemctl restart cf-ddns-updater`
- ✅ **Check internet connectivity**: `ping 8.8.8.8`
- ✅ **Verify configuration file permissions**: `ls -la /etc/cf-ddns/cf-ddns.conf`
- ✅ **Test with verbose logging**: Add `-verbose` flag
- ✅ **Run once manually**: Add `-once` flag to test

---

## Common Errors

### Configuration Errors

#### "Configuration file not found"

**Error:**
```
Error: Configuration file not found: cf-ddns.conf
```

**Solutions:**
1. **Check file path:**
   ```bash
   # Verify file exists
   ls -la /etc/cf-ddns/cf-ddns.conf
   
   # Use absolute path
   cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf
   ```

2. **Check permissions:**
   ```bash
   # Fix permissions
   sudo chmod 644 /etc/cf-ddns/cf-ddns.conf
   sudo chown cf-ddns:cf-ddns /etc/cf-ddns/cf-ddns.conf
   ```

#### "TOML parse error"

**Error:**
```
Error: TOML parse error at line 5: invalid character
```

**Solutions:**
1. **Check TOML syntax:**
   ```toml
   # Correct
   api_token = "your_token_here"
   
   # Incorrect - missing quotes
   api_token = your_token_here
   ```

2. **Validate brackets:**
   ```toml
   # Correct
   [[domains]]
   name = "example.com"
   
   # Incorrect - wrong bracket type
   [domains]
   name = "example.com"
   ```

3. **Check special characters:**
   ```toml
   # Escape backslashes in Windows paths
   log = "C:\\logs\\cf-ddns.log"
   ```

### DNS Update Errors

#### "Zone not found"

**Error:**
```
Error: Zone not found for domain: example.com
```

**Solutions:**
1. **Verify domain in Cloudflare:**
   - Log into Cloudflare dashboard
   - Ensure domain is added and active
   - Check nameservers are pointing to Cloudflare

2. **Check API token permissions:**
   - Token needs `Zone:Zone:Read` permission
   - Token must have access to the correct zones

#### "DNS record update failed"

**Error:**
```
Error: Failed to update A record for example.com: API error
```

**Solutions:**
1. **Check API permissions:**
   ```bash
   # Test API token
   curl -X GET "https://api.cloudflare.com/client/v4/user/tokens/verify" \
        -H "Authorization: Bearer your_token" \
        -H "Content-Type: application/json"
   ```

2. **Verify record exists:**
   - Check if DNS record exists in Cloudflare dashboard
   - If missing, create it manually first

3. **Check proxied status conflicts:**
   ```toml
   # Some record types can't be proxied
   [[domains]]
   name = "mx.example.com"
   record_types = "A"
   proxied = false  # MX records can't be proxied
   ```

---

## Authentication Issues

### API Token Problems

#### "Authentication failed"

**Error:**
```
Error: Authentication failed: Invalid API token
```

**Diagnosis:**
```bash
# Test token manually
curl -X GET "https://api.cloudflare.com/client/v4/user/tokens/verify" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json"
```

**Solutions:**
1. **Regenerate API token:**
   - Go to Cloudflare dashboard → API Tokens
   - Create new token with correct permissions
   - Update configuration file

2. **Check token permissions:**
   Required permissions:
   - `Zone:Zone:Read`
   - `Zone:DNS:Edit`

3. **Verify token scope:**
   - Ensure token has access to your zones
   - Check if token is expired

#### "Insufficient permissions"

**Error:**
```
Error: API error: Insufficient permissions to edit DNS records
```

**Solutions:**
1. **Update token permissions:**
   - Add `Zone:DNS:Edit` permission
   - Ensure zone scope includes your domains

2. **Check zone ownership:**
   - Verify you own the domain in Cloudflare
   - Check if domain is in a different Cloudflare account

### Legacy API Key Issues

#### "Invalid email or API key"

**Error:**
```
Error: Authentication failed: Invalid email or API key
```

**Solutions:**
1. **Switch to API tokens (recommended):**
   ```toml
   [cloudflare]
   api_token = "your_new_token"
   # Remove old api_key and email lines
   ```

2. **Verify email and key:**
   - Check email matches Cloudflare account
   - Regenerate global API key if needed

---

## Network Problems

### IP Detection Issues

#### "Failed to get current IP address"

**Error:**
```
Warning: Failed to get IPv4 address: connection timeout
```

**Solutions:**
1. **Check internet connectivity:**
   ```bash
   # Test basic connectivity
   ping 8.8.8.8
   
   # Test HTTPS connectivity
   curl -I https://v4.fetch-ip.com
   ```

2. **Firewall/proxy issues:**
   ```bash
   # Check if behind corporate firewall
   curl -v https://v4.fetch-ip.com
   
   # Try different IP services manually
   curl https://ipv4.icanhazip.com
   ```

3. **IPv6 connectivity problems:**
   ```bash
   # Test IPv6 connectivity
   ping6 2001:4860:4860::8888
   curl -6 https://v6.fetch-ip.com
   ```

#### "All IP detection services failed"

**Error:**
```
Error: All IP detection services failed
```

**Solutions:**
1. **Check DNS resolution:**
   ```bash
   # Test DNS
   nslookup v4.fetch-ip.com
   nslookup ipv4.icanhazip.com
   ```

2. **Corporate network restrictions:**
   - Contact network administrator
   - Request access to IP detection services
   - Consider running from different network

3. **Temporary service outages:**
   - Wait and retry
   - Check service status pages
   - Use verbose logging to see which services are failing

### Cloudflare API Connectivity

#### "Connection to Cloudflare API failed"

**Error:**
```
Error: Failed to connect to Cloudflare API: connection timeout
```

**Solutions:**
1. **Test API connectivity:**
   ```bash
   # Test API endpoint
   curl -I https://api.cloudflare.com/client/v4/user
   ```

2. **Check firewall rules:**
   - Ensure port 443 (HTTPS) is open
   - Check if corporate firewall blocks API access

3. **DNS resolution issues:**
   ```bash
   # Test DNS resolution
   nslookup api.cloudflare.com
   ```

---

## Service Issues

### Linux systemd Problems

#### "Service failed to start"

**Error:**
```
Job for cf-ddns-updater.service failed
```

**Diagnosis:**
```bash
# Check service status
sudo systemctl status cf-ddns-updater

# View detailed logs
sudo journalctl -u cf-ddns-updater --no-pager

# Check service file
sudo systemctl cat cf-ddns-updater
```

**Solutions:**
1. **Fix service file:**
   ```ini
   # Correct service file
   [Unit]
   Description=Cloudflare DDNS Updater
   After=network-online.target
   Wants=network-online.target
   
   [Service]
   Type=simple
   ExecStart=/usr/local/bin/cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf
   Restart=always
   RestartSec=30
   User=cf-ddns
   Group=cf-ddns
   
   [Install]
   WantedBy=multi-user.target
   ```

2. **Fix binary permissions:**
   ```bash
   sudo chmod +x /usr/local/bin/cf-ddns-updater
   sudo chown root:root /usr/local/bin/cf-ddns-updater
   ```

3. **Fix user issues:**
   ```bash
   # Create system user
   sudo useradd --system --no-create-home --shell /bin/false cf-ddns
   
   # Fix config permissions
   sudo chown cf-ddns:cf-ddns /etc/cf-ddns/cf-ddns.conf
   sudo chmod 600 /etc/cf-ddns/cf-ddns.conf
   ```

#### "Service keeps restarting"

**Symptoms:**
- Service status shows constant restarts
- Logs show repeated start/exit cycles

**Diagnosis:**
```bash
# Check restart count
sudo systemctl status cf-ddns-updater

# View recent logs
sudo journalctl -u cf-ddns-updater --since "10 minutes ago"
```

**Solutions:**
1. **Check configuration:**
   ```bash
   # Test config manually
   sudo -u cf-ddns cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once
   ```

2. **Fix interval setting:**
   ```toml
   # Ensure interval is set for continuous operation
   interval = 300  # Don't set to 0 for service mode
   ```

3. **Check dependencies:**
   - Ensure network is available before service starts
   - Add proper service dependencies

### Windows Service Issues

#### "Service won't start" (Windows)

**Solutions:**
1. **Check NSSM installation:**
   ```cmd
   # Verify service is installed
   nssm status "Cloudflare DDNS Updater"
   
   # Check service configuration
   nssm get "Cloudflare DDNS Updater" Application
   nssm get "Cloudflare DDNS Updater" Arguments
   ```

2. **Fix file paths:**
   ```cmd
   # Use absolute paths
   nssm set "Cloudflare DDNS Updater" Application "C:\cf-ddns\cf-ddns-updater.exe"
   nssm set "Cloudflare DDNS Updater" Arguments "-config C:\cf-ddns\cf-ddns.conf"
   ```

3. **Check Windows Event Viewer:**
   - Open Event Viewer
   - Check Windows Logs → Application
   - Look for cf-ddns-updater errors

---

## Logging & Debugging

### Enable Debug Logging

#### Temporary Debug Session

```bash
# Run with maximum verbosity
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once

# Log to file for analysis
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -log /tmp/debug.log -once
```

#### Permanent Debug Logging

```toml
# Add to configuration file
verbose = true
```

### Log Analysis

#### Common Log Patterns

**Normal Operation:**
```
Starting DNS update process...
Current IPv4 address: 203.0.113.1
Processing domain: example.com
Successfully processed domain: example.com
```

**Authentication Issues:**
```
Error: Authentication failed: Invalid API token
```

**Network Issues:**
```
Warning: Failed to get IPv4 address: connection timeout
Trying fallback IP service...
```

**DNS Update Issues:**
```
Error: Failed to update A record for example.com: record not found
```

### Log Locations

| Platform | Location | Command |
|----------|----------|---------|
| Linux (systemd) | systemd journal | `journalctl -u cf-ddns-updater` |
| Linux (manual) | Specify with `-log` | `tail -f /path/to/log.txt` |
| Windows (service) | Windows Event Log | Event Viewer |
| Windows (manual) | Specify with `-log` | `type C:\path\to\log.txt` |

### Advanced Debugging

#### Network Traffic Analysis

```bash
# Monitor network connections (Linux)
sudo netstat -tulpn | grep cf-ddns-updater

# Trace DNS queries
dig @8.8.8.8 example.com A
dig @8.8.8.8 example.com AAAA
```

#### API Debugging

```bash
# Test Cloudflare API manually
curl -X GET "https://api.cloudflare.com/client/v4/zones" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" | jq

# Get specific zone records
curl -X GET "https://api.cloudflare.com/client/v4/zones/ZONE_ID/dns_records" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" | jq
```

---

## Performance Issues

### High Resource Usage

#### CPU Usage

**Symptoms:**
- High CPU usage during updates
- System becomes unresponsive

**Solutions:**
1. **Increase update interval:**
   ```toml
   interval = 600  # Update every 10 minutes instead of 5
   ```

2. **Reduce number of domains:**
   - Consolidate similar domains
   - Remove unused domain entries

3. **Check for loops:**
   - Ensure configuration doesn't cause update loops
   - Monitor logs for excessive updates

#### Memory Usage

**Solutions:**
1. **Check for memory leaks:**
   ```bash
   # Monitor memory usage
   top -p $(pgrep cf-ddns-updater)
   ```

2. **Restart service periodically:**
   ```bash
   # Add to crontab for weekly restart
   0 2 * * 0 systemctl restart cf-ddns-updater
   ```

### Slow IP Detection

#### Multiple Service Timeouts

**Symptoms:**
- Long delays before IP detection
- Multiple timeout warnings in logs

**Solutions:**
1. **Check network latency:**
   ```bash
   # Test latency to IP services
   ping v4.fetch-ip.com
   curl -w "%{time_total}" https://v4.fetch-ip.com
   ```

2. **Firewall optimization:**
   - Whitelist IP detection services
   - Reduce firewall inspection overhead

### DNS Update Delays

#### Slow Cloudflare API

**Solutions:**
1. **Check API status:**
   - Visit [Cloudflare Status Page](https://www.cloudflarestatus.com/)
   - Monitor API response times

2. **Optimize requests:**
   - Reduce number of domains per update cycle
   - Use batch operations when possible (future feature)

---

## Getting Help

### Before Seeking Help

Gather this information:

1. **System Information:**
   - OS and version: `uname -a` (Linux) or `systeminfo` (Windows)
   - cf-ddns-updater version: `cf-ddns-updater -version`

2. **Configuration:**
   - Sanitized config file (remove API tokens)
   - Command line arguments used

3. **Logs:**
   - Recent logs with `-verbose` flag
   - Error messages and timestamps

4. **Network Information:**
   - Internet connectivity test results
   - Firewall/proxy configuration

### Support Channels

1. **📖 Documentation:**
   - [Configuration Guide](/cf-ddns-updater/configuration/)
   - [Security Guide](/cf-ddns-updater/security/)

2. **🤝 Community:**
   - [GitHub Discussions](https://github.com/jlbyh2o/cf-ddns-updater/discussions) - Q&A and help
   - [GitHub Issues](https://github.com/jlbyh2o/cf-ddns-updater/issues) - Bug reports

3. **🔒 Security Issues:**
   - [Security Advisory](https://github.com/jlbyh2o/cf-ddns-updater/security/advisories/new) - Private reporting

### Creating Effective Bug Reports

Include this information in bug reports:

```markdown
## System Information
- OS: Ubuntu 20.04 LTS
- Architecture: x86_64
- cf-ddns-updater version: 1.0.0

## Configuration
[Sanitized configuration file content]

## Expected Behavior
What should happen?

## Actual Behavior
What actually happens?

## Steps to Reproduce
1. Step one
2. Step two
3. Error occurs

## Logs
[Log output with -verbose flag]

## Additional Context
Any other relevant information
```

---

## Recovery Procedures

### Complete Reset

If everything is broken, start fresh:

```bash
# Stop service
sudo systemctl stop cf-ddns-updater

# Backup configuration
sudo cp /etc/cf-ddns/cf-ddns.conf /etc/cf-ddns/cf-ddns.conf.backup

# Reinstall (Linux)
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash

# Restore configuration
sudo cp /etc/cf-ddns/cf-ddns.conf.backup /etc/cf-ddns/cf-ddns.conf

# Test and restart
sudo systemctl start cf-ddns-updater
```

### Emergency DNS Update

If the service is completely broken but you need to update DNS:

```bash
# Manual DNS update using curl
curl -X PUT "https://api.cloudflare.com/client/v4/zones/ZONE_ID/dns_records/RECORD_ID" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     --data '{"type":"A","name":"example.com","content":"YOUR.IP.ADDRESS.HERE","ttl":300}'
```

---

Remember: Most issues are configuration-related. Always test with `-verbose -once` flags first!