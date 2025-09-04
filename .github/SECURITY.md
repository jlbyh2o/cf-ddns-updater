# Security Policy

## Supported Versions

We actively support the following versions of Cloudflare DDNS Updater with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

### How to Report

We take security vulnerabilities seriously. If you discover a security vulnerability in Cloudflare DDNS Updater, please report it responsibly:

**🔒 For Security Issues:**
- **Email**: [Create a GitHub Security Advisory](https://github.com/jlbyh2o/cf-ddns-updater/security/advisories/new)
- **Alternative**: Open a private issue by emailing the maintainer directly

**📋 For Non-Security Issues:**
- Open a regular [GitHub Issue](https://github.com/jlbyh2o/cf-ddns-updater/issues/new)

### What to Include

When reporting a security vulnerability, please include:

1. **Description**: Clear description of the vulnerability
2. **Impact**: Potential impact and attack scenarios
3. **Reproduction**: Step-by-step instructions to reproduce
4. **Environment**: 
   - Operating system and version
   - Application version
   - Go version (if building from source)
5. **Proof of Concept**: Code or screenshots demonstrating the issue
6. **Suggested Fix**: If you have ideas for remediation

### Response Timeline

We are committed to responding to security reports promptly:

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution**: Within 30 days for critical issues

### Disclosure Policy

We follow responsible disclosure practices:

1. **Private Reporting**: Report vulnerabilities privately first
2. **Coordinated Disclosure**: We'll work with you on timing
3. **Public Disclosure**: After fix is released and users have time to update
4. **Credit**: We'll acknowledge your contribution (if desired)

## Security Best Practices

### For Users

**🔐 Configuration Security:**
- Store configuration files with restricted permissions (`600` or `640`)
- Use environment variables for sensitive data when possible
- Regularly rotate Cloudflare API tokens
- Use API tokens with minimal required permissions

**🛡️ System Security:**
- Run the service with a dedicated non-root user when possible
- Keep the application updated to the latest version
- Monitor logs for unusual activity
- Use firewall rules to restrict network access

**📁 File Permissions:**
```bash
# Recommended file permissions
sudo chmod 600 /etc/cf-ddns-updater/cf-ddns.conf
sudo chown cf-ddns:cf-ddns /etc/cf-ddns-updater/cf-ddns.conf
```

### For Developers

**🔍 Code Security:**
- All dependencies are regularly updated
- Static analysis tools are used in CI/CD
- Input validation on all user-provided data
- Secure handling of API credentials

**🚀 Build Security:**
- Reproducible builds with checksums
- Signed releases (planned for future versions)
- Automated security scanning in GitHub Actions

## Security Features

### Current Security Measures

- **🔒 Secure API Communication**: All Cloudflare API calls use HTTPS
- **🛡️ Input Validation**: Configuration and IP addresses are validated
- **📝 Minimal Logging**: Sensitive data is not logged
- **🔐 Token Handling**: API tokens are handled securely in memory
- **⚡ Fail-Safe**: Application fails securely on errors

### Planned Security Enhancements

- **📋 Configuration Validation**: Enhanced config file validation
- **🔍 Security Scanning**: Automated dependency vulnerability scanning
- **📦 Signed Releases**: GPG-signed release binaries
- **🛡️ Sandboxing**: Containerized deployment options

## Threat Model

### Assets Protected
- Cloudflare API tokens
- DNS record configurations
- System configuration files
- Network communication

### Potential Threats
- **Configuration Exposure**: Unauthorized access to config files
- **Token Compromise**: API token theft or misuse
- **Man-in-the-Middle**: Interception of API communications
- **Privilege Escalation**: Running with excessive permissions
- **Denial of Service**: Resource exhaustion attacks

### Mitigations
- Secure file permissions and ownership
- HTTPS-only API communication
- Input validation and sanitization
- Principle of least privilege
- Rate limiting and error handling

## Security Contacts

- **Primary**: GitHub Security Advisories (preferred)
- **Maintainer**: Available through GitHub profile
- **Community**: GitHub Discussions for general security questions

## Acknowledgments

We appreciate the security research community and will acknowledge researchers who report vulnerabilities responsibly:

- Hall of Fame for security researchers (coming soon)
- Public acknowledgment in release notes
- CVE attribution when applicable

---

**Last Updated**: August 2025  
**Next Review**: January 2026

> 🔒 **Remember**: When in doubt, report it! We'd rather investigate a false positive than miss a real security issue.