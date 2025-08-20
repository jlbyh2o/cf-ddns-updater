#!/bin/bash

# Build script for Unix-like systems
echo "Building Cloudflare DDNS Updater..."

# Clean previous builds
rm -rf bin
mkdir -p bin

# Build for Windows 64-bit
echo "Building for Windows 64-bit..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/cf-ddns-updater-windows-amd64.exe .
if [ $? -ne 0 ]; then
    echo "Failed to build for Windows 64-bit"
    exit 1
fi

# Build for Linux x86-64
echo "Building for Linux x86-64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/cf-ddns-updater-linux-amd64 .
if [ $? -ne 0 ]; then
    echo "Failed to build for Linux x86-64"
    exit 1
fi

# Copy example config
cp cf-ddns.conf.example bin/

echo ""
echo "Build completed successfully!"
echo "Binaries are available in the 'bin' directory:"
echo "- cf-ddns-updater-windows-amd64.exe (Windows 64-bit)"
echo "- cf-ddns-updater-linux-amd64 (Linux x86-64)"
echo "- cf-ddns.conf.example (Example configuration)"
echo ""

# Make the Linux binary executable
chmod +x bin/cf-ddns-updater-linux-amd64

echo "Linux binary has been made executable."