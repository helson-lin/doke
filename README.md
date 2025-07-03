# Doke 🐳

<p align="center">
  <strong>一个强大的 Docker 终端工具：快速转换容器配置为执行命令</strong>
</p>

<p align="center">
  <a href="#功能特性">功能特性</a> •
  <a href="#安装">安装</a> •
  <a href="#使用方法">使用方法</a> •
  <a href="#示例">示例</a> •
  <a href="#贡献">贡献</a>
</p>

---

## 🚀 功能特性

- **🔧 容器配置转换**: 将运行中的 Docker 容器配置转换为 `docker run` 命令
- **📝 Docker Compose 生成**: 自动生成 Docker Compose YAML 配置文件
- **🌐 镜像源配置**: 自动配置 Docker 镜像源，提升拉取速度
- **🔄 版本管理**: 支持版本查看和更新
- **💻 跨平台支持**: 支持 Linux 和 macOS 系统
- **⚡ 快速高效**: 轻量级命令行工具，执行速度快

## 📦 安装

### 使用 Homebrew 安装

```bash
brew install helson-lin/tap/doke
```

### 手动安装

1. 下载最新版本
```bash
# 下载对应平台的二进制文件
wget https://github.com/helson-lin/doke/releases/latest/download/doke-{version}-{platform}.tar.gz
tar -xzf doke-{version}-{platform}.tar.gz
sudo mv doke /usr/local/bin/
```

2. 验证安装
```bash
doke --version
```

### 从源码编译

```bash
git clone https://github.com/helson-lin/doke.git
cd doke
go build -o doke
sudo mv doke /usr/local/bin/
```

## 🎯 使用方法

### 基本命令

```bash
# 查看版本信息
doke -v
doke --version
doke version

# 查看帮助信息
doke --help
doke command --help
```

### 容器配置转换

```bash
# 转换容器为 docker run 命令
doke command <container_id>

# 使用别名
doke c <container_id>

# 生成 Docker Compose 文件
doke command <container_id> -j
doke c <container_id> --json
```

### 镜像源配置

```bash
# 自动配置 Docker 镜像源
doke proxy
```

## 📋 示例

### 1. 转换容器配置

```bash
# 查看运行中的容器
docker ps

# 转换容器配置为 docker run 命令
doke command abc123def456

# 输出示例：
# docker run --name my-container -d -p 8080:80 -v /host/path:/container/path nginx:latest
```

### 2. 生成 Docker Compose 文件

```bash
# 生成 Docker Compose 配置
doke command abc123def456 -j

# 输出示例：
# 是否将 Docker Compose 配置写入文件 my-container.yml？(y/n): y
# Docker Compose 配置已成功写入文件: my-container.yml
```

### 3. 配置镜像源

```bash
# 配置 Docker 镜像源
doke proxy

# 输出示例：
# 🚀 开始配置 Docker 镜像源...
# 📁 配置文件路径: /Users/user/.orbstack/config/docker.json
# ✅ 成功读取配置文件
# 📋 发现现有镜像源: [https://docker.1panel.live]
# 🔄 正在添加镜像源: [https://docker.1ms.run https://docker.1panel.live]
# ✅ 添加镜像源: https://docker.1ms.run
# ⏭️  镜像源已存在，跳过: https://docker.1panel.live
# 🎉 Docker 镜像源配置完成！
```

## 🔧 配置说明

### 支持的镜像源

Doke 会自动配置以下镜像源：
- `https://docker.1ms.run`
- `https://docker.1panel.live`

### 配置文件位置

- **macOS (OrbStack)**: `~/.orbstack/config/docker.json`
- **Linux**: `/etc/docker/daemon.json`

## 🛠️ 开发

### 环境要求

- Go 1.21+
- Docker

### 构建

```bash
# 构建当前平台
go build -o doke

# 构建所有平台
bash build.sh v1.0.0
```

### 测试

```bash
# 运行测试
go test ./...

# 测试特定功能
go test ./cmd
```

## 📝 更新日志

### v0.0.1
- ✨ 初始版本发布
- 🔧 支持容器配置转换
- 📝 支持 Docker Compose 生成
- 🌐 支持镜像源配置
- 🔄 支持版本查看

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Cobra](https://github.com/spf13/cobra) - 强大的命令行框架
- [Docker](https://www.docker.com/) - 容器化平台
- [OrbStack](https://orbstack.dev/) - macOS 上的 Docker 替代方案

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/helson-lin">helson-lin</a>
</p>

