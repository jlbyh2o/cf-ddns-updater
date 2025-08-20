@echo off
REM Build script for Windows
echo Building Cloudflare DDNS Updater...

REM Clean previous builds
if exist "bin" rmdir /s /q bin
mkdir bin

REM Build for Windows 64-bit
echo Building for Windows 64-bit...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o bin\cf-ddns-updater-windows-amd64.exe .
if %errorlevel% neq 0 (
    echo Failed to build for Windows 64-bit
    exit /b 1
)

REM Build for Windows ARM64
echo Building for Windows ARM64...
set GOOS=windows
set GOARCH=arm64
go build -ldflags "-s -w" -o bin\cf-ddns-updater-windows-arm64.exe .
if %errorlevel% neq 0 (
    echo Failed to build for Windows ARM64
    exit /b 1
)

REM Build for Linux x86-64
echo Building for Linux x86-64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o bin/cf-ddns-updater-linux-amd64 .
if %errorlevel% neq 0 (
    echo Failed to build for Linux x86-64
    exit /b 1
)

REM Build for Linux ARM
echo Building for Linux ARM...
set GOOS=linux
set GOARCH=arm
go build -ldflags "-s -w" -o bin/cf-ddns-updater-linux-arm .
if %errorlevel% neq 0 (
    echo Failed to build for Linux ARM
    exit /b 1
)

REM Build for Linux ARM64
echo Building for Linux ARM64...
set GOOS=linux
set GOARCH=arm64
go build -ldflags "-s -w" -o bin/cf-ddns-updater-linux-arm64 .
if %errorlevel% neq 0 (
    echo Failed to build for Linux ARM64
    exit /b 1
)

REM Copy example config
copy cf-ddns.conf.example bin\

echo.
echo Build completed successfully!
echo Binaries are available in the 'bin' directory:
echo - cf-ddns-updater-windows-amd64.exe (Windows 64-bit)
echo - cf-ddns-updater-windows-arm64.exe (Windows ARM64)
echo - cf-ddns-updater-linux-amd64 (Linux x86-64)
echo - cf-ddns-updater-linux-arm (Linux ARM)
echo - cf-ddns-updater-linux-arm64 (Linux ARM64)
echo - cf-ddns.conf.example (Example configuration)
echo.
pause