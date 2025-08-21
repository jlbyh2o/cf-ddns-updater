# Cloudflare DDNS Updater

A reliable and configurable Dynamic DNS updater for Cloudflare written in Go. This tool automatically updates your Cloudflare DNS records with your current public IP address, supporting both IPv4 (A records) and IPv6 (AAAA records).

## Features

- **Reliable**: Uses multiple IP detection services for redundancy, with fetch-ip.com as the default
- **Configurable**: Support for A records, AAAA records, or both
- **Cross-platform**: Builds for Linux (x86-64, ARM, ARM64) and Windows (x86-64, ARM64)
- **Flexible Authentication**: Supports both API tokens and API key/email combinations
- **Continuous Mode**: Can run continuously with configurable intervals
- **Comprehensive Logging**: Detailed logging with optional file output
- **Error Handling**: Robust error handling and recovery

## Quick Start

1. **Download or Build**
   - Download pre-built binaries from releases, or
   - Build from source using the provided build scripts

2. **Create Configuration**
   ```bash
   cp cf-ddns.conf.example cf-ddns.conf
   ```

3. **Edit Configuration**
   - Add your Cloudflare API credentials
   - Configure your domains and record types

4. **Run**
   ```bash
   # Run once (Linux x86-64)
   ./cf-ddns-updater-linux-amd64 -config cf-ddns.conf
   
   # Run once (Linux ARM)
   ./cf-ddns-updater-linux-arm -config cf-ddns.conf
   
   # Run once (Linux ARM64)
   ./cf-ddns-updater-linux-arm64 -config cf-ddns.conf
   
   # Run continuously (Windows)
   cf-ddns-updater-windows-amd64.exe -config cf-ddns.conf
   ```

## Configuration

### Basic Configuration

Create a `cf-ddns.conf` file based on the example:

```toml
[cloudflare]
api_token = "your_cloudflare_api_token_here"

[[domains]]
name = "example.com"
record_types = "both"
ttl = 300
proxied = false
```

### Configuration Options

#### Cloudflare Section

- `api_token` (string): Cloudflare API token (recommended)
- `api_key` (string): Cloudflare API key (legacy, requires email)
- `email` (string): Cloudflare account email (required with api_key)
- `zone_id` (string, optional): Zone ID (auto-detected if not provided)

#### Domain Configuration

- `name` (string): Domain or subdomain name
- `record_types` (string): "A", "AAAA", or "both" (default: "both")
- `ttl` (int): DNS record TTL in seconds (default: 300)
- `proxied` (bool): Whether to proxy through Cloudflare (default: false)

#### Global Options

- `interval` (int): Update interval in seconds (0 = run once, default: 0)
- `verbose` (bool): Enable verbose logging (default: false)

### Authentication Methods

#### API Token (Recommended)

1. Go to [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens)
2. Create a custom token with:
   - **Permissions**: `Zone:DNS:Edit`, `Zone:Zone:Read`
   - **Zone Resources**: Include specific zones or all zones
3. Use the token in your config:

```toml
[cloudflare]
api_token = "your_token_here"
```

#### API Key + Email (Legacy)

```toml
[cloudflare]
api_key = "your_global_api_key"
email = "your@email.com"
```

## Command Line Options

```bash
./cf-ddns-updater [options]

Options:
  -config string
        Path to configuration file (default "cf-ddns.conf")
  -verbose
        Enable verbose logging
  -log string
        Log file path (optional, logs to stdout if not specified)
  -once
        Run once and exit (ignore interval setting)
```

## Usage Examples

### Run Once
```bash
./cf-ddns-updater -config cf-ddns.conf -once
```

### Continuous Mode
```bash
# Update every 5 minutes (300 seconds)
./cf-ddns-updater -config cf-ddns.conf
```

### With Logging
```bash
./cf-ddns-updater -config cf-ddns.conf -verbose -log /var/log/ddns.log
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

## Building from Source

### Prerequisites

- Go 1.19 or later
- Git

### Build Instructions

#### Windows
```cmd
git clone https://github.com/jlbyh2o/cf-ddns-updater.git
cd cf-ddns-updater
build.bat
```

#### Linux/macOS
```bash
git clone https://github.com/jlbyh2o/cf-ddns-updater.git
cd cf-ddns-updater
chmod +x build.sh
./build.sh
```

Built binaries will be available in the `bin/` directory:
- `cf-ddns-updater-windows-amd64.exe` (Windows x86-64)
- `cf-ddns-updater-windows-arm64.exe` (Windows ARM64)
- `cf-ddns-updater-linux-amd64` (Linux x86-64)
- `cf-ddns-updater-linux-arm` (Linux ARM)
- `cf-ddns-updater-linux-arm64` (Linux ARM64)

## Linux Installation

For Linux systems, we provide a complete installation package that integrates with systemd and follows Linux filesystem standards.

### System Installation (Recommended)

The system installation places the binary in `/usr/local/bin`, configuration in `/etc/cf-ddns`, and sets up a systemd service.

#### Prerequisites

- Linux system with systemd
- Root access (sudo)
- Go 1.19+ (for building from source)

#### Installation Steps

1. **Clone and Build**
   ```bash
   git clone https://github.com/jlbyh2o/cf-ddns-updater.git
   cd cf-ddns-updater
   make build
   ```

2. **Install System-wide**
   ```bash
   sudo make install
   ```
   
   Or use the installation script:
   ```bash
   chmod +x install.sh
   sudo ./install.sh
   ```

3. **Configure**
   ```bash
   # Copy example configuration
   sudo cp /etc/cf-ddns/cf-ddns.conf.example /etc/cf-ddns/cf-ddns.conf
   
   # Edit configuration
   sudo nano /etc/cf-ddns/cf-ddns.conf
   ```

4. **Enable and Start Service**
   ```bash
   sudo systemctl enable cf-ddns-updater
   sudo systemctl start cf-ddns-updater
   ```

5. **Check Status**
   ```bash
   sudo systemctl status cf-ddns-updater
   sudo journalctl -u cf-ddns-updater -f
   ```

#### Installation Paths

- **Binary**: `/usr/local/bin/cf-ddns-updater`
- **Configuration**: `/etc/cf-ddns/cf-ddns.conf`
- **Service**: `/etc/systemd/system/cf-ddns-updater.service`
- **User**: `cf-ddns` (created automatically)
- **Logs**: Available via `journalctl -u cf-ddns-updater`

#### Configuration File Locations

The application searches for configuration files in this order:

1. `/etc/cf-ddns/cf-ddns.conf` (system-wide)
2. `./cf-ddns.conf` (current directory)
3. `cf-ddns.conf` relative to executable

You can override this with the `-config` flag:
```bash
cf-ddns-updater -config /path/to/your/cf-ddns.conf
#### Makefile Targets

```bash
make build      # Build the binary for Linux x86-64
make build-all  # Build for all supported architectures
make install    # Install system-wide with systemd service
make install-dev# Install without systemd service
make uninstall  # Remove installation
make clean      # Clean build artifacts
make help       # Show available targets
```

#### Uninstallation

```bash
sudo make uninstall
```

Or using the installation script:
```bash
sudo ./install.sh uninstall
```

### Manual Installation

If you prefer manual installation:

1. **Build and copy binary**
   ```bash
   make build
   sudo cp cf-ddns-updater /usr/local/bin/
   sudo chmod +x /usr/local/bin/cf-ddns-updater
   ```

2. **Create configuration directory**
   ```bash
   sudo mkdir -p /etc/cf-ddns
   sudo cp cf-ddns.conf.example /etc/cf-ddns/
   ```

3. **Create systemd service** (optional)
   ```bash
   sudo cp cf-ddns-updater.service /etc/systemd/system/
   sudo systemctl daemon-reload
   ```

### Security Features

The systemd service includes comprehensive security hardening:

- Runs as dedicated `cf-ddns` user
- No new privileges
- Private temporary directory
- Protected system directories
- Restricted system calls
- Network access limited to IPv4/IPv6
- Memory execution protection

## Deployment

### Linux Service (systemd)

Create `/etc/systemd/system/cf-ddns-updater.service`:

```ini
[Unit]
Description=Cloudflare DDNS Updater
After=network.target

[Service]
Type=simple
User=ddns
WorkingDirectory=/opt/cf-ddns-updater
ExecStart=/opt/cf-ddns-updater/cf-ddns-updater-linux-amd64 -config /opt/cf-ddns-updater/cf-ddns.conf
Restart=always
RestartSec=30

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable cf-ddns-updater
sudo systemctl start cf-ddns-updater
```

### Windows Service

Use a tool like [NSSM](https://nssm.cc/) to run as a Windows service:

```cmd
nssm install "Cloudflare DDNS Updater" "C:\path\to\cf-ddns-updater-windows-amd64.exe"
nssm set "Cloudflare DDNS Updater" Arguments "-config C:\path\to\cf-ddns.conf"
nssm start "Cloudflare DDNS Updater"
```

### Docker

Create a `Dockerfile`:

```dockerfile
FROM scratch
COPY cf-ddns-updater-linux-amd64 /cf-ddns-updater
COPY cf-ddns.conf /cf-ddns.conf
ENTRYPOINT ["/cf-ddns-updater", "-config", "/cf-ddns.conf"]
```

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Verify your API token has the correct permissions
   - Check that the token isn't expired
   - Ensure zone access is granted

2. **IP Detection Failures**
   - Check internet connectivity
   - Verify firewall settings
   - The tool tries multiple services for redundancy (fetch-ip.com, icanhazip.com, ipify.org, ident.me)

3. **DNS Update Failures**
   - Verify domain ownership in Cloudflare
   - Check zone ID (if manually specified)
   - Ensure record name matches exactly

### Debug Mode

Run with verbose logging to see detailed information:

```bash
./cf-ddns-updater -config cf-ddns.conf -verbose -once
```

## Security Considerations

- Store configuration files securely with appropriate permissions
- Use API tokens instead of global API keys when possible
- Limit API token permissions to only what's needed
- Consider using environment variables for sensitive data
- Regularly rotate API credentials

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0).

Cloudflare Dynamic DNS Updater
Copyright (C) 2025

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

## Support

If you encounter issues or need help:

1. Check the troubleshooting section
2. Review the logs with verbose mode enabled
3. Open an issue with detailed information about your setup and the problem