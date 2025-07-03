# Doke 🐳

<p align="center">
  <strong>A powerful Docker terminal tool: quickly convert container configurations to executable commands</strong>
</p>

<p align="center">
  English | <a href="README.md">中文</a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#examples">Examples</a> •
  <a href="#advanced-features">Advanced Features</a> •
  <a href="#contributing">Contributing</a>
</p>

---

## 🚀 Features

- **🔧 Container Configuration Conversion**: Convert running Docker container configurations to `docker run` commands
- **📝 Docker Compose Generation**: Automatically generate Docker Compose YAML configuration files
- **🔍 Real-time Container Monitoring**: Real-time monitoring of container status, resource usage, and log output
- **🧹 Resource Cleanup**: Automatically clean up unused Docker resources to free up system space
- **🌐 Registry Mirror Configuration**: Automatically configure Docker registry mirrors to improve pull speeds
- **🔄 Version Management**: Support for version viewing and updates
- **💻 Cross-platform Support**: Support for Linux and macOS systems
- **⚡ Fast and Efficient**: Lightweight command-line tool with fast execution

[![asciicast](https://asciinema.org/a/6N8kSJONyoUnO2R4Sn6gUeQG2.svg)](https://asciinema.org/a/6N8kSJONyoUnO2R4Sn6gUeQG2)

## 📦 Installation

### Install using Homebrew

```bash
brew install helson-lin/tap/doke
```

### Manual Installation

1. Download the latest release
```bash
# Download the binary for your platform
wget https://github.com/helson-lin/doke/releases/latest/download/doke-{version}-{platform}.tar.gz
tar -xzf doke-{version}-{platform}.tar.gz
sudo mv doke /usr/local/bin/
```

2. Verify installation
```bash
doke --version
```

### Build from Source

```bash
git clone https://github.com/helson-lin/doke.git
cd doke
go build -o doke
sudo mv doke /usr/local/bin/
```

## 🎯 Usage

### Basic Commands

```bash
# View version information
doke -v
doke --version
doke version

# View help information
doke --help
doke <command> --help
```

### Container Configuration Conversion

```bash
# Convert container to docker run command
doke command <container_id>

# Use alias
doke c <container_id>

# Generate Docker Compose file
doke command <container_id> -j
doke c <container_id> --json
```

### Real-time Container Monitoring

```bash
# Monitor container status in real-time
doke inspect <container_id>

# Monitoring includes:
# - Container basic information (name, ID, image, status)
# - Real-time resource usage (CPU, memory, network)
# - Real-time log output
# - Port mappings and volume mounts
```

### Resource Cleanup

```bash
# Clean up unused Docker resources
doke clear

# Clean up all unused images (including tagged images)
doke clear --all

# Force cleanup, skip confirmation prompts
doke clear --force

# Show detailed system information
doke clear --info
```

### Registry Mirror Configuration

```bash
# Automatically configure Docker registry mirrors
doke proxy
```

### Language Settings

```bash
# Switch to Chinese
doke lang zh

# Switch to English
doke lang en
```

## 📋 Examples

### 1. Convert Container Configuration

```bash
# View running containers
docker ps

# Convert container configuration to docker run command
doke command abc123def456

# Example output:
# docker run --name my-container -d -p 8080:80 -v /host/path:/container/path nginx:latest
```

### 2. Generate Docker Compose File

```bash
# Generate Docker Compose configuration
doke command abc123def456 -j

# Example output:
# Write Docker Compose configuration to file my-container.yml? (y/n): y
# Docker Compose configuration successfully written to file: my-container.yml
```

### 3. Real-time Container Monitoring

```bash
# Start monitoring container
doke inspect my-container

# Example output:
# 📦 Container name: my-container
# 🏷️  Container ID: abc123def456
# 🖼️  Image: nginx:latest
# 📍 Status: running
# ⏰ Uptime: 2h30m15s
# 
# 🌐 Port mappings:
#    8080 -> 80/tcp
# 
# 💾 Volume mounts:
#    /host/data -> /var/www/html (bind)
# 
# 🔍 Starting real-time container monitoring...
# 💡 Press Ctrl+C to exit monitoring
# 📊 CPU: 15.2% | 💾 Memory: 128.5MB/1024.0MB (12.5%) | 🌐 Network: ↓2.1MB ↑0.8MB
```

### 4. Clean Up Unused Resources

```bash
# Clean up system resources
doke clear

# Example output:
# 📊 Current Docker resource usage:
#    📦 Containers: 15 (Running: 3, Stopped: 12)
#    🖼️  Images: 25 (Dangling: 8) Total size: 3.2GB
#    🌐 Networks: 8 (Custom: 3)
#    💾 Volumes: 12
# 
# 🗑️  Resources to be cleaned:
#    ✅ Clean stopped containers
#    ✅ Clean dangling images
#    ✅ Clean unused networks
#    ✅ Clean unused volumes
# 
# Confirm cleanup operation? (y/n): y
# 
# 🧹 Cleanup results:
# 📦 Containers: 12
# 🖼️  Images: 8
# 🌐 Networks: 2
# 💾 Volumes: 3
# 💽 Space reclaimed: 1.8 GB
```

### 5. Configure Registry Mirrors

```bash
# Configure Docker registry mirrors
doke proxy

# Example output:
# 🚀 Starting Docker registry mirror configuration...
# 📁 Configuration file path: /Users/user/.orbstack/config/docker.json
# ✅ Successfully read configuration file
# 📋 Found existing registry mirrors: [https://docker.1panel.live]
# 🔄 Adding registry mirrors: [https://docker.1ms.run https://docker.1panel.live]
# ✅ Added registry mirror: https://docker.1ms.run
# ⏭️  Registry mirror already exists, skipping: https://docker.1panel.live
# 🎉 Docker registry mirror configuration complete!
```

## 🔧 Advanced Features

### Real-time Container Monitoring

The `doke inspect` command provides powerful real-time container monitoring capabilities:

- **📊 Real-time Resource Monitoring**: Updates CPU, memory, and network usage every 2 seconds
- **📜 Real-time Log Streaming**: Displays the latest log output from containers
- **📋 Detailed Information Display**: Container basic info, port mappings, volume mounts, etc.
- **🔄 Continuous Monitoring**: Monitors continuously until manually stopped (Ctrl+C)

### Intelligent Resource Cleanup

The `doke clear` command provides intelligent Docker resource cleanup:

- **🔍 Resource Analysis**: Automatically analyzes current system resource usage
- **🎯 Precise Cleanup**: Only cleans up truly unused resources
- **📊 Cleanup Statistics**: Shows detailed cleanup results and reclaimed space
- **⚡ Safe Cleanup**: Provides confirmation mechanism to avoid deleting important resources

### Multi-language Support

Doke supports bilingual interface in Chinese and English:

- **🌐 Auto Detection**: Automatically selects interface language based on system language
- **🔄 Flexible Switching**: Use `doke lang` command to switch languages anytime
- **📝 Complete Translation**: All features and messages have complete Chinese and English translations

## 🔧 Configuration

### Supported Registry Mirrors

Doke automatically configures the following registry mirrors:
- `https://docker.1ms.run`
- `https://docker.1panel.live`

### Configuration File Locations

- **macOS (OrbStack)**: `~/.orbstack/config/docker.json`
- **Linux**: `/etc/docker/daemon.json`

### Language Configuration

Doke supports the following languages:
- **Chinese**: `zh` (default)
- **English**: `en`

## 🛠️ Development

### Requirements

- Go 1.21+
- Docker

### Build

```bash
# Build for current platform
go build -o doke

# Build for all platforms
bash build.sh v1.0.0
```

### Testing

```bash
# Run tests
go test ./...

# Test specific functionality
go test ./cmd
```

## 📝 Changelog

### v0.0.2 (Planned)
- 🔍 Added real-time container monitoring
- 🧹 Added resource cleanup functionality
- 🌐 Added multi-language support
- 📊 Improved resource usage statistics display
- 🐛 Fixed known issues

### v0.0.1
- ✨ Initial release
- 🔧 Support for container configuration conversion
- 📝 Support for Docker Compose generation
- 🌐 Support for registry mirror configuration
- 🔄 Support for version viewing

## 🤝 Contributing

Issues and Pull Requests are welcome!

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - Powerful command-line framework
- [Docker](https://www.docker.com/) - Containerization platform
- [OrbStack](https://orbstack.dev/) - Docker alternative for macOS

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/helson-lin">helson-lin</a>
</p> 