---
layout: home
title: "Cloudflare DDNS Updater"
---

# 🚀 Cloudflare DDNS Updater

**A reliable and lightweight Dynamic DNS updater for Cloudflare** - Keep your domains pointing to your dynamic IP address automatically!

[![Build Status](https://github.com/jlbyh2o/cf-ddns-updater/actions/workflows/build.yml/badge.svg)](https://github.com/jlbyh2o/cf-ddns-updater/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/jlbyh2o/cf-ddns-updater)](https://github.com/jlbyh2o/cf-ddns-updater/blob/main/go.mod)
[![Release](https://img.shields.io/github/v/release/jlbyh2o/cf-ddns-updater)](https://github.com/jlbyh2o/cf-ddns-updater/releases)
[![Downloads](https://img.shields.io/github/downloads/jlbyh2o/cf-ddns-updater/total)](https://github.com/jlbyh2o/cf-ddns-updater/releases)

## ✨ Key Features

<div class="features-grid">
  <div class="feature-card">
    <h3>🔄 Automatic Updates</h3>
    <p>Keeps your DNS records in sync with your dynamic IP address without manual intervention.</p>
  </div>
  
  <div class="feature-card">
    <h3>🛡️ Reliable</h3>
    <p>Uses multiple IP detection services with fallback mechanisms for maximum uptime.</p>
  </div>
  
  <div class="feature-card">
    <h3>⚙️ Flexible</h3>
    <p>Supports IPv4 (A records), IPv6 (AAAA records), or both with customizable settings.</p>
  </div>
  
  <div class="feature-card">
    <h3>🔐 Secure</h3>
    <p>Uses API token authentication with minimal required permissions and secure practices.</p>
  </div>
  
  <div class="feature-card">
    <h3>🐧 Linux Ready</h3>
    <p>Complete systemd integration with security hardening and automated installation.</p>
  </div>
  
  <div class="feature-card">
    <h3>🚀 Lightweight</h3>
    <p>Single binary with no dependencies, minimal resource usage, and cross-platform support.</p>
  </div>
</div>

## 🚀 Quick Start

### Linux One-Line Install
```bash
curl -sSL https://raw.githubusercontent.com/jlbyh2o/cf-ddns-updater/main/install-linux.sh | sudo bash
```

### Manual Installation
1. [Download the latest release](https://github.com/jlbyh2o/cf-ddns-updater/releases/latest)
2. Extract and place the binary in your PATH
3. Follow our [detailed installation guide](installation.html)

## 📖 Documentation

<div class="docs-grid">
  <a href="getting-started.html" class="doc-card">
    <h3>🏁 Getting Started</h3>
    <p>Set up your first DDNS configuration in minutes</p>
  </a>
  
  <a href="installation.html" class="doc-card">
    <h3>📦 Installation</h3>
    <p>Install on Linux, Windows, or build from source</p>
  </a>
  
  <a href="configuration.html" class="doc-card">
    <h3>⚙️ Configuration</h3>
    <p>Complete configuration reference and examples</p>
  </a>
  
  <a href="troubleshooting.html" class="doc-card">
    <h3>🔧 Troubleshooting</h3>
    <p>Solve common issues and debug problems</p>
  </a>
  
  <a href="security.html" class="doc-card">
    <h3>🔒 Security</h3>
    <p>Security best practices and configuration</p>
  </a>
  
  <a href="api.html" class="doc-card">
    <h3>🔌 API Reference</h3>
    <p>Integration details and API documentation</p>
  </a>
</div>

## 🌐 Supported Platforms

- **Linux**: x86_64, ARM, ARM64 (Ubuntu, Debian, CentOS, RHEL, Fedora, Arch)
- **Windows**: x86_64, ARM64 (Windows 10, 11, Server)
- **Docker**: Multi-architecture container images available

## 💡 Use Cases

- **Home Server**: Keep your home server accessible with a dynamic IP
- **Self-Hosted Services**: Maintain domain access for self-hosted applications  
- **Remote Access**: Ensure SSH, VPN, or remote desktop connectivity
- **Web Hosting**: Point your domain to dynamic hosting environments
- **Development**: Keep development servers accessible during testing

## 🤝 Community & Support

- **📖 Documentation**: Comprehensive guides and tutorials
- **🐛 Issue Tracker**: [Report bugs and request features](https://github.com/jlbyh2o/cf-ddns-updater/issues)
- **💬 Discussions**: [Community Q&A and help](https://github.com/jlbyh2o/cf-ddns-updater/discussions)
- **🔒 Security**: [Report security issues responsibly](security.html)

## ⭐ Star Us on GitHub

If you find this project helpful, please consider giving us a star on GitHub! It helps others discover the project and motivates continued development.

[⭐ Star on GitHub](https://github.com/jlbyh2o/cf-ddns-updater){:.btn .btn-primary}

---

*Cloudflare DDNS Updater is open source software licensed under GPL-3.0*