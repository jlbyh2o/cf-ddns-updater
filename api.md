---
layout: page
title: "API Reference"
permalink: /api/
---

# 🔌 API Reference

Technical reference for integrating with and extending Cloudflare DDNS Updater.

<div class="toc">
<h3>API Topics</h3>
<ul>
  <li><a href="#command-line-interface">Command Line Interface</a></li>
  <li><a href="#configuration-api">Configuration API</a></li>
  <li><a href="#exit-codes">Exit Codes</a></li>
  <li><a href="#integration-examples">Integration Examples</a></li>
  <li><a href="#monitoring-apis">Monitoring APIs</a></li>
  <li><a href="#extending-functionality">Extending Functionality</a></li>
</ul>
</div>

## Command Line Interface

### Basic Syntax

```bash
cf-ddns-updater [OPTIONS]
```

### Command Line Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `-config` | string | `"cf-ddns.conf"` | Path to configuration file |
| `-verbose` | flag | `false` | Enable verbose logging output |
| `-log` | string | `stdout` | Log file path (optional) |
| `-once` | flag | `false` | Run once and exit (ignore interval) |
| `-version` | flag | `false` | Show version and exit |

### Usage Examples

```bash
# Basic usage with default config
cf-ddns-updater

# Specify configuration file
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf

# Run once with verbose output
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -verbose -once

# Log to specific file
cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -log /var/log/cf-ddns.log

# Show version information
cf-ddns-updater -version
```

---

## Configuration API

### Configuration File Structure

The application uses TOML format for configuration:

```toml
[cloudflare]
api_token = "string"
api_key = "string"      # Alternative to api_token
email = "string"        # Required with api_key

[[domains]]
name = "string"         # Required
record_types = "string" # "A", "AAAA", or "both"
ttl = integer          # Seconds (60-2147483647)
proxied = boolean      # true/false

# Optional global settings
interval = integer     # Seconds (0 = run once)
verbose = boolean      # true/false
```

### Configuration Validation Rules

#### Cloudflare Section

**Authentication (Choose One):**
```toml
# Option 1: API Token (Recommended)
[cloudflare]
api_token = "required_string_min_40_chars"

# Option 2: Legacy API Key + Email
[cloudflare]
api_key = "required_string_37_chars"
email = "required_valid_email"
```

#### Domain Section

**Required Fields:**
```toml
[[domains]]
name = "required_valid_domain_or_subdomain"
```

**Optional Fields with Defaults:**
```toml
record_types = "both"    # Default: "both"
ttl = 300               # Default: 300 (5 minutes)
proxied = false         # Default: false
```

**Validation Rules:**
- `name`: Must be valid domain/subdomain format
- `record_types`: Must be exactly "A", "AAAA", or "both"
- `ttl`: Must be integer between 60 and 2147483647
- `proxied`: Must be boolean true/false

#### Global Settings

```toml
# Optional settings
interval = 0           # Default: 0 (run once)
verbose = false        # Default: false
```

---

## Exit Codes

The application uses standard exit codes for automation and monitoring:

| Exit Code | Meaning | Description |
|-----------|---------|-------------|
| `0` | Success | Operation completed successfully |
| `1` | General Error | Generic error (check logs for details) |
| `2` | Configuration Error | Invalid configuration file or settings |
| `3` | Authentication Error | Invalid API credentials |
| `4` | Network Error | Network connectivity issues |
| `5` | API Error | Cloudflare API errors |
| `6` | DNS Error | DNS resolution or update errors |
| `127` | Command Not Found | Binary not found or not executable |

### Exit Code Usage Examples

```bash
#!/bin/bash
# Script using exit codes for error handling

cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -once
EXIT_CODE=$?

case $EXIT_CODE in
    0)
        echo "DNS update successful"
        ;;
    2)
        echo "Configuration error - check config file"
        exit 1
        ;;
    3)
        echo "Authentication failed - check API credentials"
        exit 1
        ;;
    4)
        echo "Network error - check connectivity"
        exit 1
        ;;
    *)
        echo "Unknown error (exit code: $EXIT_CODE)"
        exit 1
        ;;
esac
```

---

## Integration Examples

### Systemd Integration

#### Service File

```ini
[Unit]
Description=Cloudflare DDNS Updater
Documentation=https://jlbyh2o.github.io/cf-ddns-updater/
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=30
User=cf-ddns
Group=cf-ddns

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/etc/cf-ddns

[Install]
WantedBy=multi-user.target
```

#### Service Management

```bash
# Enable and start
systemctl enable cf-ddns-updater
systemctl start cf-ddns-updater

# Status and logs
systemctl status cf-ddns-updater
journalctl -u cf-ddns-updater -f

# Reload configuration
systemctl reload cf-ddns-updater

# Restart service
systemctl restart cf-ddns-updater
```

### Docker Integration

#### Dockerfile

```dockerfile
FROM scratch

# Copy binary and config
COPY cf-ddns-updater-linux-amd64 /cf-ddns-updater
COPY cf-ddns.conf /config/cf-ddns.conf

# Run as non-root user
USER 65534:65534

# Expose configuration directory
VOLUME ["/config"]

# Health check
HEALTHCHECK --interval=5m --timeout=30s --start-period=1m \
    CMD /cf-ddns-updater -config /config/cf-ddns.conf -once || exit 1

ENTRYPOINT ["/cf-ddns-updater"]
CMD ["-config", "/config/cf-ddns.conf"]
```

#### Docker Compose

```yaml
version: '3.8'

services:
  cf-ddns-updater:
    image: cf-ddns-updater:latest
    container_name: cf-ddns-updater
    restart: unless-stopped
    volumes:
      - ./config:/config:ro
    environment:
      - TZ=UTC
    healthcheck:
      test: ["/cf-ddns-updater", "-config", "/config/cf-ddns.conf", "-once"]
      interval: 5m
      timeout: 30s
      retries: 3
      start_period: 1m
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    read_only: true
    tmpfs:
      - /tmp
```

### Kubernetes Integration

#### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cf-ddns-updater
  labels:
    app: cf-ddns-updater
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cf-ddns-updater
  template:
    metadata:
      labels:
        app: cf-ddns-updater
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
        fsGroup: 65534
      containers:
      - name: cf-ddns-updater
        image: cf-ddns-updater:latest
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          readOnlyRootFilesystem: true
        volumeMounts:
        - name: config
          mountPath: /config
          readOnly: true
        resources:
          requests:
            memory: "16Mi"
            cpu: "10m"
          limits:
            memory: "32Mi"
            cpu: "50m"
        livenessProbe:
          exec:
            command:
            - /cf-ddns-updater
            - -config
            - /config/cf-ddns.conf
            - -once
          initialDelaySeconds: 60
          periodSeconds: 300
      volumes:
      - name: config
        secret:
          secretName: cf-ddns-config
```

#### ConfigMap and Secret

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: cf-ddns-config
type: Opaque
stringData:
  cf-ddns.conf: |
    [cloudflare]
    api_token = "your_api_token_here"
    
    [[domains]]
    name = "example.com"
    record_types = "both"
    ttl = 300
    proxied = false
    
    interval = 300
    verbose = false
```

### Cron Integration

```bash
# Update DNS every 5 minutes via cron
*/5 * * * * /usr/local/bin/cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -once >> /var/log/cf-ddns-cron.log 2>&1

# Daily log rotation
0 0 * * * logrotate /etc/logrotate.d/cf-ddns-cron
```

### Windows Service Integration

#### NSSM Configuration

```cmd
# Install service
nssm install "Cloudflare DDNS Updater" "C:\cf-ddns\cf-ddns-updater.exe"

# Configure arguments
nssm set "Cloudflare DDNS Updater" Arguments "-config C:\cf-ddns\cf-ddns.conf"

# Set working directory
nssm set "Cloudflare DDNS Updater" AppDirectory "C:\cf-ddns"

# Configure logging
nssm set "Cloudflare DDNS Updater" AppStdout "C:\cf-ddns\logs\stdout.log"
nssm set "Cloudflare DDNS Updater" AppStderr "C:\cf-ddns\logs\stderr.log"

# Set service properties
nssm set "Cloudflare DDNS Updater" DisplayName "Cloudflare DDNS Updater"
nssm set "Cloudflare DDNS Updater" Description "Automatic DNS updater for Cloudflare"
nssm set "Cloudflare DDNS Updater" Start SERVICE_AUTO_START

# Start service
nssm start "Cloudflare DDNS Updater"
```

---

## Monitoring APIs

### Health Check Endpoint (Future Feature)

**Planned HTTP health check endpoint:**

```bash
# Future feature - HTTP health check
curl http://localhost:8080/health
# Response: {"status": "healthy", "last_update": "2023-12-07T10:30:00Z"}
```

### Metrics Export (Future Feature)

**Planned Prometheus metrics:**

```
# HELP cf_ddns_updates_total Total number of DNS updates attempted
# TYPE cf_ddns_updates_total counter
cf_ddns_updates_total{domain="example.com",type="A",status="success"} 42

# HELP cf_ddns_update_duration_seconds Time spent updating DNS records
# TYPE cf_ddns_update_duration_seconds histogram
cf_ddns_update_duration_seconds_bucket{le="1.0"} 10
cf_ddns_update_duration_seconds_bucket{le="2.5"} 15
cf_ddns_update_duration_seconds_bucket{le="+Inf"} 20

# HELP cf_ddns_last_update_timestamp Unix timestamp of last successful update
# TYPE cf_ddns_last_update_timestamp gauge
cf_ddns_last_update_timestamp{domain="example.com"} 1702459800
```

### Log-Based Monitoring

#### Structured Logging

Current log format (text-based):
```
2023-12-07T10:30:00Z [INFO] Starting DNS update process...
2023-12-07T10:30:01Z [INFO] Current IPv4 address: 203.0.113.1
2023-12-07T10:30:02Z [INFO] Successfully updated A record for example.com
```

**Future structured logging (JSON):**
```json
{
  "timestamp": "2023-12-07T10:30:00Z",
  "level": "info",
  "message": "DNS update successful",
  "domain": "example.com",
  "record_type": "A",
  "old_ip": "203.0.113.0",
  "new_ip": "203.0.113.1",
  "duration_ms": 1234
}
```

#### Log Parsing Examples

**ELK Stack (Logstash):**
```ruby
filter {
  if [fields][service] == "cf-ddns-updater" {
    grok {
      match => { 
        "message" => "%{TIMESTAMP_ISO8601:timestamp} \[%{LOGLEVEL:level}\] %{GREEDYDATA:log_message}" 
      }
    }
    
    if "Successfully updated" in [log_message] {
      grok {
        match => { 
          "log_message" => "Successfully updated %{WORD:record_type} record for %{HOSTNAME:domain}" 
        }
      }
    }
  }
}
```

**Splunk Search:**
```
index=infrastructure source="cf-ddns-updater" 
| rex "Successfully updated (?<record_type>\w+) record for (?<domain>\S+)"
| stats count by domain, record_type
| sort -count
```

---

## Extending Functionality

### Plugin Architecture (Future)

**Planned plugin interface:**

```go
// Plugin interface for custom IP detection services
type IPDetector interface {
    GetIPv4(ctx context.Context) (string, error)
    GetIPv6(ctx context.Context) (string, error)
    Name() string
}

// Plugin interface for custom DNS providers
type DNSProvider interface {
    UpdateRecord(ctx context.Context, record DNSRecord) error
    GetRecords(ctx context.Context, domain string) ([]DNSRecord, error)
    Name() string
}
```

### Custom IP Detection Services

**Adding custom IP services (future configuration):**

```toml
[ip_detection]
# Custom IPv4 detection services
[[ip_detection.ipv4_services]]
url = "https://custom-ip-service.com/ipv4"
method = "GET"
response_format = "text"  # or "json"
json_path = ".ip"         # if json format
timeout = 10              # seconds

# Custom IPv6 detection services  
[[ip_detection.ipv6_services]]
url = "https://custom-ip-service.com/ipv6"
method = "GET"
response_format = "json"
json_path = ".ipv6_address"
timeout = 10
```

### Webhook Integration (Future)

**Planned webhook support:**

```toml
[webhooks]
# Pre-update webhook
[[webhooks.pre_update]]
url = "https://webhook.example.com/pre-update"
method = "POST"
headers = {"Authorization" = "Bearer token"}
payload_template = """
{
  "domain": "{{ .Domain }}",
  "old_ip": "{{ .OldIP }}",
  "new_ip": "{{ .NewIP }}",
  "timestamp": "{{ .Timestamp }}"
}
"""

# Post-update webhook
[[webhooks.post_update]]
url = "https://webhook.example.com/post-update"
method = "POST"
on_success = true
on_failure = true
```

### Custom DNS Providers (Future)

**Support for additional DNS providers:**

```toml
[dns_provider]
type = "cloudflare"  # Default

# Future support
# type = "route53"
# type = "azure"
# type = "gcp"
# type = "custom"

[aws_route53]  # Future configuration
access_key_id = "your_key"
secret_access_key = "your_secret"
region = "us-east-1"

[custom_provider]  # Future configuration
api_endpoint = "https://api.example.com"
auth_header = "X-API-Key"
auth_value = "your_key"
```

---

## Error Handling

### Error Response Format

**Standard error output:**
```
Error: [ERROR_TYPE] error_description
```

**Error types:**
- `CONFIG` - Configuration file errors
- `AUTH` - Authentication errors  
- `NETWORK` - Network connectivity errors
- `API` - Cloudflare API errors
- `DNS` - DNS resolution/update errors

### Error Handling in Scripts

```bash
#!/bin/bash
# Robust error handling example

set -euo pipefail  # Exit on error, undefined vars, pipe failures

CONFIG_FILE="/etc/cf-ddns/cf-ddns.conf"
LOG_FILE="/var/log/cf-ddns-script.log"

# Function to log messages
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# Function to handle errors
handle_error() {
    local exit_code=$1
    local error_msg="$2"
    
    log "ERROR: $error_msg (exit code: $exit_code)"
    
    case $exit_code in
        2)
            log "Configuration error - checking config file..."
            if [[ ! -f "$CONFIG_FILE" ]]; then
                log "Config file not found: $CONFIG_FILE"
            else
                log "Config file exists, checking permissions..."
                ls -la "$CONFIG_FILE" | tee -a "$LOG_FILE"
            fi
            ;;
        3)
            log "Authentication error - API credentials may be invalid"
            ;;
        4)
            log "Network error - checking connectivity..."
            if ping -c 1 8.8.8.8 >/dev/null 2>&1; then
                log "Basic internet connectivity OK"
            else
                log "No internet connectivity"
            fi
            ;;
        *)
            log "Unknown error code: $exit_code"
            ;;
    esac
    
    exit $exit_code
}

# Main execution
log "Starting DNS update..."

if cf-ddns-updater -config "$CONFIG_FILE" -once; then
    log "DNS update completed successfully"
else
    exit_code=$?
    handle_error $exit_code "DNS update failed"
fi
```

---

## Development API

### Building from Source

```bash
# Clone repository
git clone https://github.com/jlbyh2o/cf-ddns-updater.git
cd cf-ddns-updater

# Build for current platform
go build -o cf-ddns-updater

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o cf-ddns-updater-linux-amd64

# Build with version info
go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o cf-ddns-updater
```

### Testing

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests (requires config)
go test -tags=integration ./...
```

### Contributing

**Code contribution workflow:**
1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Run linting and tests
5. Submit pull request

**Code style requirements:**
- Follow Go formatting standards (`gofmt`)
- Add comprehensive tests for new features
- Update documentation for API changes
- Use semantic versioning for releases

---

This API reference covers the current stable version 1.0.0. Future versions may include additional features and endpoints as outlined in the "Future" sections above.

For the latest API updates and additional integration examples, see the [GitHub repository](https://github.com/jlbyh2o/cf-ddns-updater) and [community discussions](https://github.com/jlbyh2o/cf-ddns-updater/discussions).