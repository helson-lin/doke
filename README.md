# Doke 🐳

<p align="center">
  <strong>一个强大的 Docker 终端工具：快速转换容器配置为执行命令</strong>
</p>

<p align="center">
  <a href="README_EN.md">English</a> | 中文
</p>

<p align="center">
  <a href="#功能特性">功能特性</a> •
  <a href="#安装">安装</a> •
  <a href="#使用方法">使用方法</a> •
  <a href="#示例">示例</a> •
  <a href="#高级功能">高级功能</a> •
  <a href="#贡献">贡献</a>
</p>

---

## 🚀 功能特性

- **🔧 容器配置转换**: 将运行中的 Docker 容器配置转换为 `docker run` 命令
- **📝 Docker Compose 生成**: 自动生成 Docker Compose YAML 配置文件
- **🔍 容器实时监控**: 实时监控容器状态、资源使用情况和日志输出
- **🧹 资源清理**: 自动清理未使用的 Docker 资源，释放系统空间
- **🌐 镜像源配置**: 自动配置 Docker 镜像源，提升拉取速度
- **🔄 版本管理**: 支持版本查看和更新
- **💻 跨平台支持**: 支持 Linux 和 macOS 系统
- **⚡ 快速高效**: 轻量级命令行工具，执行速度快

[![asciicast](https://asciinema.org/a/6N8kSJONyoUnO2R4Sn6gUeQG2.svg)](https://asciinema.org/a/6N8kSJONyoUnO2R4Sn6gUeQG2)

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
doke <command> --help
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

### 容器实时监控

```bash
# 实时监控容器状态
doke inspect <container_id>

# 监控内容包括：
# - 容器基本信息（名称、ID、镜像、状态）
# - 实时资源使用情况（CPU、内存、网络）
# - 实时日志输出
# - 端口映射和卷挂载信息
```

### 资源清理

```bash
# 清理未使用的 Docker 资源
doke clear

# 清理所有未使用的镜像（包括有标签的镜像）
doke clear --all

# 强制清理，跳过确认提示
doke clear --force

# 查看详细系统信息
doke clear --info
```

### 镜像源配置

```bash
# 自动配置 Docker 镜像源
doke proxy
```

### 语言设置

```bash
# 切换到中文
doke lang zh

# 切换到英文
doke lang en
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

### 3. 实时监控容器

```bash
# 开始监控容器
doke inspect my-container

# 输出示例：
# 📦 容器名称: my-container
# 🏷️  容器ID: abc123def456
# 🖼️  镜像: nginx:latest
# 📍 状态: running
# ⏰ 运行时间: 2h30m15s
# 
# 🌐 端口映射:
#    8080 -> 80/tcp
# 
# 💾 卷挂载:
#    /host/data -> /var/www/html (bind)
# 
# 🔍 开始实时监控容器状态...
# 💡 按 Ctrl+C 退出监控
# 📊 CPU: 15.2% | 💾 内存: 128.5MB/1024.0MB (12.5%) | 🌐 网络: ↓2.1MB ↑0.8MB
```

### 4. 清理未使用资源

```bash
# 清理系统资源
doke clear

# 输出示例：
# 📊 当前 Docker 资源使用情况:
#    📦 容器: 15 个 (运行中: 3, 已停止: 12)
#    🖼️  镜像: 25 个 (悬空: 8) 总大小: 3.2GB
#    🌐 网络: 8 个 (自定义: 3)
#    💾 卷: 12 个
# 
# 🗑️  将要清理的资源:
#    ✅ 清理已停止的容器
#    ✅ 清理悬空镜像
#    ✅ 清理未使用的网络
#    ✅ 清理未使用的卷
# 
# 确认执行清理操作？(y/n): y
# 
# 🧹 清理结果:
# 📦 容器: 12 个
# 🖼️  镜像: 8 个
# 🌐 网络: 2 个
# 💾 卷: 3 个
# 💽 回收空间: 1.8 GB
```

### 5. 配置镜像源

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

## 🔧 高级功能

### 容器实时监控功能

`doke inspect` 命令提供了强大的容器实时监控功能：

- **📊 实时资源监控**: 每2秒更新一次 CPU、内存、网络使用情况
- **📜 实时日志流**: 显示容器的最新日志输出
- **📋 详细信息展示**: 容器基本信息、端口映射、卷挂载等
- **🔄 持续监控**: 持续监控直到用户手动停止（Ctrl+C）

### 智能资源清理

`doke clear` 命令提供了智能的 Docker 资源清理功能：

- **🔍 资源分析**: 自动分析当前系统资源使用情况
- **🎯 精准清理**: 只清理真正未使用的资源
- **📊 清理统计**: 显示详细的清理结果和回收空间
- **⚡ 安全清理**: 提供确认机制，避免误删重要资源

### 多语言支持

Doke 支持中英文双语界面：

- **🌐 自动检测**: 根据系统语言自动选择界面语言
- **🔄 灵活切换**: 使用 `doke lang` 命令随时切换语言
- **📝 完整翻译**: 所有功能和提示信息都有完整的中英文翻译

## 🔧 配置说明

### 支持的镜像源

Doke 会自动配置以下镜像源：
- `https://docker.1ms.run`
- `https://docker.1panel.live`

### 配置文件位置

- **macOS (OrbStack)**: `~/.orbstack/config/docker.json`
- **Linux**: `/etc/docker/daemon.json`

### 语言配置

Doke 支持以下语言：
- **中文**: `zh` (默认)
- **英文**: `en`

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

### v0.0.2 (计划中)
- 🔍 新增容器实时监控功能
- 🧹 新增资源清理功能
- 🌐 新增多语言支持
- 📊 优化资源使用统计显示
- 🐛 修复已知问题

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

