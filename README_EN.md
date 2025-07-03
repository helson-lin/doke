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
  <a href="#contributing">Contributing</a>
</p>

---

## 🚀 Features

- **🔧 Container Configuration Conversion**: Convert running Docker container configurations to `docker run` commands
- **📝 Docker Compose Generation**: Automatically generate Docker Compose YAML configuration files
- **🌐 Registry Mirror Configuration**: Automatically configure Docker registry mirrors to improve pull speeds
- **🔄 Version Management**: Support for version viewing and updates
- **💻 Cross-platform Support**: Support for Linux and macOS systems
- **⚡ Fast and Efficient**: Lightweight command-line tool with fast execution

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
doke command --help
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

### Registry Mirror Configuration

```bash
# Automatically configure Docker registry mirrors
doke proxy
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

### 3. Configure Registry Mirrors

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

## 🔧 Configuration

### Supported Registry Mirrors

Doke automatically configures the following registry mirrors:
- `https://docker.1ms.run`
- `https://docker.1panel.live`

### Configuration File Locations

- **macOS (OrbStack)**: `~/.orbstack/config/docker.json`
- **Linux**: `/etc/docker/daemon.json`

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