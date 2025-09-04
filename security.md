---
layout: page
title: "Security Guide"
permalink: /security/
---

# 🔒 Security Guide

Comprehensive security guidelines for deploying and maintaining Cloudflare DDNS Updater safely.

<div class="toc">
<h3>Security Topics</h3>
<ul>
  <li><a href="#authentication-security">Authentication Security</a></li>
  <li><a href="#configuration-security">Configuration Security</a></li>
  <li><a href="#network-security">Network Security</a></li>
  <li><a href="#system-security">System Security</a></li>
  <li><a href="#monitoring-alerting">Monitoring & Alerting</a></li>
  <li><a href="#incident-response">Incident Response</a></li>
</ul>
</div>

## Authentication Security

### API Token Best Practices

#### Use API Tokens (Not API Keys)

<div class="alert alert-success">
<strong>✅ Recommended:</strong> API Tokens with minimal permissions
</div>

```toml
[cloudflare]
api_token = "your_limited_scope_token"
```

<div class="alert alert-danger">
<strong>❌ Avoid:</strong> Global API Keys
</div>

```toml
# Avoid this - too much access
[cloudflare]
api_key = "global_api_key"
email = "your@email.com"
```

#### Token Configuration

**Required Permissions:**
- `Zone:Zone:Read` - To read zone information
- `Zone:DNS:Edit` - To update DNS records

**Zone Resources:**
```
Specific zones: *.example.com, *.mydomain.net
```

**IP Address Filtering (Optional):**
- Add your server's IP addresses for additional security
- Use CIDR notation for ranges: `192.168.1.0/24`

#### Token Rotation

```bash
# Regular token rotation schedule
# 1. Generate new token in Cloudflare dashboard
# 2. Update configuration file
# 3. Test new token
# 4. Revoke old token

# Test new token before deployment
curl -X GET "https://api.cloudflare.com/client/v4/user/tokens/verify" \
     -H "Authorization: Bearer NEW_TOKEN" \
     -H "Content-Type: application/json"
```

**Recommended Rotation Schedule:**
- **Production**: Every 90 days
- **Development**: Every 30 days
- **High-risk environments**: Every 30 days

### Token Storage Security

#### Environment Variables (Future Feature)

```bash
# Secure environment variable storage
export CF_API_TOKEN="your_token"
echo 'export CF_API_TOKEN="your_token"' >> ~/.bashrc

# Use in systemd service
Environment=CF_API_TOKEN=your_token
```

#### Secrets Management

**For Production Environments:**

1. **HashiCorp Vault:**
   ```bash
   # Store token in Vault
   vault kv put secret/cf-ddns api_token="your_token"
   
   # Retrieve in service startup script
   CF_API_TOKEN=$(vault kv get -field=api_token secret/cf-ddns)
   ```

2. **Docker Secrets:**
   ```bash
   # Create secret
   echo "your_token" | docker secret create cf_api_token -
   
   # Use in container
   docker run --secret cf_api_token cf-ddns-updater
   ```

3. **Kubernetes Secrets:**
   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: cf-ddns-secret
   type: Opaque
   data:
     api-token: <base64-encoded-token>
   ```

---

## Configuration Security

### File Permissions

#### Linux Permissions

```bash
# Secure configuration file
sudo chown cf-ddns:cf-ddns /etc/cf-ddns/cf-ddns.conf
sudo chmod 600 /etc/cf-ddns/cf-ddns.conf

# Verify permissions
ls -la /etc/cf-ddns/cf-ddns.conf
# Should show: -rw------- cf-ddns cf-ddns
```

#### Windows Permissions

```cmd
# Remove inheritance and set specific permissions
icacls C:\cf-ddns\cf-ddns.conf /inheritance:d
icacls C:\cf-ddns\cf-ddns.conf /grant:r "System:(F)"
icacls C:\cf-ddns\cf-ddns.conf /grant:r "Administrators:(F)"
icacls C:\cf-ddns\cf-ddns.conf /grant:r "ServiceAccount:(R)"
icacls C:\cf-ddns\cf-ddns.conf /remove "Users"
```

### Configuration Validation

#### Secure Configuration Template

```toml
[cloudflare]
# Use API token with minimal permissions
api_token = "your_limited_scope_token"

[[domains]]
name = "example.com"
record_types = "both"
ttl = 300
proxied = false  # Set to true only if needed

# Security-focused settings
interval = 300      # Don't set too low to avoid rate limiting
verbose = false     # Disable verbose logging in production
```

#### Configuration Checklist

- ✅ **API token has minimal required permissions**
- ✅ **No global API key used**
- ✅ **Configuration file has restrictive permissions (600)**
- ✅ **No sensitive data in version control**
- ✅ **Verbose logging disabled in production**
- ✅ **Update interval not too aggressive**

### Backup and Recovery

```bash
# Secure configuration backup
sudo cp /etc/cf-ddns/cf-ddns.conf /etc/cf-ddns/cf-ddns.conf.backup
sudo chmod 600 /etc/cf-ddns/cf-ddns.conf.backup
sudo chown cf-ddns:cf-ddns /etc/cf-ddns/cf-ddns.conf.backup

# Encrypted backup for off-site storage
gpg --cipher-algo AES256 --compress-algo 1 --s2k-digest-algo SHA256 \
    --cert-digest-algo SHA256 --encrypt --armor \
    -r your@email.com /etc/cf-ddns/cf-ddns.conf
```

---

## Network Security

### Firewall Configuration

#### Outbound Rules (Required)

**Linux (iptables):**
```bash
# Allow HTTPS to Cloudflare API
iptables -A OUTPUT -p tcp --dport 443 -d api.cloudflare.com -j ACCEPT

# Allow DNS resolution
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT
iptables -A OUTPUT -p tcp --dport 53 -j ACCEPT

# Allow HTTPS to IP detection services
iptables -A OUTPUT -p tcp --dport 443 -d v4.fetch-ip.com -j ACCEPT
iptables -A OUTPUT -p tcp --dport 443 -d ipv4.icanhazip.com -j ACCEPT
iptables -A OUTPUT -p tcp --dport 443 -d api.ipify.org -j ACCEPT
iptables -A OUTPUT -p tcp --dport 443 -d ident.me -j ACCEPT
```

**Windows Firewall:**
```cmd
# Allow outbound HTTPS
netsh advfirewall firewall add rule name="CF-DDNS HTTPS" dir=out action=allow protocol=TCP localport=443

# Allow outbound DNS
netsh advfirewall firewall add rule name="CF-DDNS DNS" dir=out action=allow protocol=UDP localport=53
```

#### Network Segmentation

**DMZ Deployment:**
- Place DDNS updater in DMZ network segment
- Limit access to only required external services
- Monitor network traffic for anomalies

**Docker Network Isolation:**
```bash
# Create isolated network
docker network create --driver bridge cf-ddns-net

# Run container in isolated network
docker run --network cf-ddns-net cf-ddns-updater
```

### TLS/SSL Security

#### Certificate Validation

The application automatically validates TLS certificates for all HTTPS connections:

- ✅ **Certificate chain validation**
- ✅ **Hostname verification**
- ✅ **Expiration checking**
- ✅ **Strong cipher suites only**

#### TLS Configuration

```go
// Application uses secure TLS configuration by default
// No user configuration required - handled automatically
```

---

## System Security

### User Account Security

#### Dedicated Service Account

**Linux:**
```bash
# Create system user (done by install script)
sudo useradd --system --no-create-home --shell /bin/false cf-ddns

# Verify account
id cf-ddns
# Should show: no groups except cf-ddns
```

**Windows:**
```cmd
# Create service account
net user cf-ddns-service /add /system
net localgroup "Log on as a service" cf-ddns-service /add
```

#### Privilege Minimization

**Linux systemd service hardening:**
```ini
[Service]
# User/Group isolation
User=cf-ddns
Group=cf-ddns

# File system restrictions
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/etc/cf-ddns
PrivateTmp=true
PrivateDevices=true

# Network restrictions
RestrictAddressFamilies=AF_INET AF_INET6
RestrictNamespaces=true

# System call restrictions
SystemCallFilter=@system-service
SystemCallErrorNumber=EPERM

# Capabilities
NoNewPrivileges=true
CapabilityBoundingSet=
AmbientCapabilities=

# Memory protection
MemoryDenyWriteExecute=true
RestrictRealtime=true
LockPersonality=true

# Misc hardening
RemoveIPC=true
RestrictSUIDSGID=true
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true
```

### Container Security

#### Docker Security

**Dockerfile Security:**
```dockerfile
# Use minimal base image
FROM scratch

# Non-root user
USER 65534:65534

# Read-only root filesystem
WORKDIR /app
COPY cf-ddns-updater /app/
COPY cf-ddns.conf /app/

# Security options
ENTRYPOINT ["/app/cf-ddns-updater", "-config", "/app/cf-ddns.conf"]
```

**Runtime Security:**
```bash
# Run with security options
docker run -d \
    --name cf-ddns-updater \
    --read-only \
    --tmpfs /tmp \
    --user 65534:65534 \
    --cap-drop ALL \
    --security-opt no-new-privileges:true \
    --restart unless-stopped \
    cf-ddns-updater
```

#### Kubernetes Security

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cf-ddns-updater
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
        fsGroup: 65534
        seccompProfile:
          type: RuntimeDefault
      containers:
      - name: cf-ddns-updater
        image: cf-ddns-updater:latest
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          readOnlyRootFilesystem: true
          runAsNonRoot: true
```

---

## Monitoring & Alerting

### Security Monitoring

#### Log Analysis

**Security-relevant log events to monitor:**

1. **Authentication failures:**
   ```
   Error: Authentication failed: Invalid API token
   ```

2. **Rate limiting:**
   ```
   Error: API rate limit exceeded
   ```

3. **Configuration changes:**
   ```
   Configuration file modified
   ```

4. **Unusual IP changes:**
   ```
   IP changed from X.X.X.X to Y.Y.Y.Y (significant change)
   ```

#### Log Monitoring Setup

**Linux (rsyslog):**
```bash
# Configure log forwarding
echo "local1.*  @@your-siem-server:514" >> /etc/rsyslog.conf

# Update service to log to syslog
ExecStart=/usr/local/bin/cf-ddns-updater -config /etc/cf-ddns/cf-ddns.conf -log syslog
```

**ELK Stack Integration:**
```yaml
# Filebeat configuration
filebeat.inputs:
- type: journald
  id: cf-ddns-logs
  include_matches:
    - "UNIT=cf-ddns-updater.service"
```

### Alerting Rules

#### Critical Alerts

1. **Service down for >15 minutes**
2. **Authentication failures >3 in 10 minutes**
3. **No successful DNS updates in >2 hours**
4. **Rate limiting errors**

#### Warning Alerts

1. **IP detection service failures**
2. **DNS update delays >5 minutes**
3. **Configuration file modifications**

**Example (Prometheus AlertManager):**
```yaml
groups:
- name: cf-ddns-alerts
  rules:
  - alert: CFDDNSServiceDown
    expr: up{job="cf-ddns-updater"} == 0
    for: 15m
    annotations:
      summary: "CF DDNS Updater service is down"
      
  - alert: CFDDNSAuthFailure
    expr: increase(cf_ddns_auth_failures_total[10m]) > 3
    annotations:
      summary: "Multiple authentication failures detected"
```

---

## Incident Response

### Security Incident Types

#### 1. API Token Compromise

**Immediate Actions:**
1. **Revoke compromised token** in Cloudflare dashboard
2. **Generate new token** with same permissions
3. **Update configuration** with new token
4. **Restart service** to use new token
5. **Review logs** for unauthorized DNS changes

**Investigation:**
```bash
# Check recent DNS changes in Cloudflare
curl -X GET "https://api.cloudflare.com/client/v4/zones/ZONE_ID/dns_records" \
     -H "Authorization: Bearer NEW_TOKEN" | jq '.result[] | select(.modified_on > "2023-01-01")'

# Review application logs
journalctl -u cf-ddns-updater --since "24 hours ago" | grep -i error
```

#### 2. Unauthorized DNS Changes

**Immediate Actions:**
1. **Verify current DNS records** in Cloudflare dashboard
2. **Check application logs** for recent updates
3. **Restore correct IP addresses** if needed
4. **Investigate source** of unauthorized changes

**Verification Script:**
```bash
#!/bin/bash
# Verify DNS records match expected IPs

DOMAIN="example.com"
EXPECTED_IP="203.0.113.1"
CURRENT_IP=$(dig +short $DOMAIN @8.8.8.8)

if [ "$CURRENT_IP" != "$EXPECTED_IP" ]; then
    echo "ALERT: DNS record for $DOMAIN has unexpected IP: $CURRENT_IP"
    # Send alert notification
fi
```

#### 3. Service Compromise

**Immediate Actions:**
1. **Stop the service** immediately
2. **Isolate the system** from network
3. **Preserve logs** for forensic analysis
4. **Assess scope** of compromise
5. **Rebuild system** from known-good state

### Recovery Procedures

#### Clean Reinstallation

```bash
# Complete removal and reinstall
sudo systemctl stop cf-ddns-updater
sudo systemctl disable cf-ddns-updater

# Remove all files
sudo rm /usr/local/bin/cf-ddns-updater
sudo rm /etc/systemd/system/cf-ddns-updater.service
sudo rm -rf /etc/cf-ddns/

# Fresh installation
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash

# Restore configuration from secure backup
sudo cp /secure/backup/location/cf-ddns.conf /etc/cf-ddns/
sudo chmod 600 /etc/cf-ddns/cf-ddns.conf
sudo chown cf-ddns:cf-ddns /etc/cf-ddns/cf-ddns.conf
```

### Communication Plan

#### Internal Communication

1. **Security Team** - Immediate notification
2. **Network Operations** - Infrastructure impact
3. **Management** - Business impact assessment

#### External Communication

1. **DNS Changes** - Notify dependent services
2. **Regulatory Requirements** - If applicable
3. **Customer Impact** - If services affected

---

## Security Hardening Checklist

### Pre-Deployment

- [ ] **API token created with minimal permissions**
- [ ] **Configuration file has secure permissions**
- [ ] **Service runs as non-root user**
- [ ] **Firewall rules configured (outbound only)**
- [ ] **Log monitoring configured**
- [ ] **Backup procedures established**

### Post-Deployment

- [ ] **Service status monitoring active**
- [ ] **Log analysis rules configured**
- [ ] **Alerting rules tested**
- [ ] **Incident response plan documented**
- [ ] **Regular security reviews scheduled**

### Ongoing Maintenance

- [ ] **Monthly log reviews**
- [ ] **Quarterly token rotation**
- [ ] **Annual security assessments**
- [ ] **Regular backup testing**
- [ ] **Security update monitoring**

---

## Compliance Considerations

### Data Protection

- **No personal data** is processed or stored
- **IP addresses** are technical identifiers, not personal data in most jurisdictions
- **API tokens** should be treated as sensitive credentials

### Audit Requirements

**Log Retention:**
- Maintain logs for compliance periods (typically 1-7 years)
- Secure log storage with integrity protection
- Regular log backup and testing

**Access Controls:**
- Document who has access to configuration files
- Regular access reviews and updates
- Strong authentication for administrative access

### Industry Standards

**ISO 27001 Alignment:**
- Risk assessment and treatment
- Security monitoring and incident response
- Secure configuration management

**NIST Framework:**
- Identify, Protect, Detect, Respond, Recover
- Security controls implementation
- Continuous monitoring

---

## Reporting Security Issues

### Responsible Disclosure

**For security vulnerabilities:**
1. **DO NOT** create public GitHub issues
2. **Use** [GitHub Security Advisories](https://github.com/jlbyh2o/cf-ddns-updater/security/advisories/new)
3. **Provide** detailed vulnerability information
4. **Allow** reasonable time for fixes before disclosure

### What to Include

- **Vulnerability description**
- **Steps to reproduce**
- **Impact assessment**
- **Suggested remediation**
- **Your contact information**

**Response Timeline:**
- Initial response: 48 hours
- Status update: 7 days
- Resolution target: 30 days for critical issues

---

Security is a shared responsibility. Stay vigilant, keep systems updated, and follow security best practices!