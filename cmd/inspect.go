package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/helson-lin/doke/i18n"
	"github.com/spf13/cobra"
)

// 实时检测容器的运行情况，当用户主动终止时退出
func inspectContainer(containerId string) {
	// 1. 获取容器配置
	config, err := getDockerContainerConfig(containerId)
	if err != nil {
		rootCmd.PrintErrf(i18n.T("error.container_config", err) + "\n")
		return
	}

	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		rootCmd.PrintErrf(i18n.T("error.docker_client", err) + "\n")
		return
	}
	defer cli.Close()

	// 打印基本信息
	printContainerBasicInfo(config)

	// 设置信号处理，允许用户使用 Ctrl+C 退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建用于停止监控的 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动实时监控
	go monitorContainerStats(ctx, cli, containerId)
	go monitorContainerLogs(ctx, cli, containerId)

	fmt.Println("\n" + i18n.T("inspect.monitoring_start"))
	fmt.Println(i18n.T("inspect.monitoring_tip"))
	fmt.Println(strings.Repeat("=", 60))

	// 等待用户中断信号
	<-sigChan
	fmt.Println("\n\n" + i18n.T("inspect.monitoring_stopped"))
}

// 打印容器基本信息
func printContainerBasicInfo(config *types.ContainerJSON) {
	fmt.Printf(i18n.T("inspect.container_name", strings.TrimPrefix(config.Name, "/")) + "\n")
	fmt.Printf(i18n.T("inspect.container_id", config.ID[:12]) + "\n")
	fmt.Printf(i18n.T("inspect.image", config.Config.Image) + "\n")
	fmt.Printf(i18n.T("inspect.status", config.State.Status) + "\n")

	if config.State.Running {
		startedAt, err := time.Parse(time.RFC3339Nano, config.State.StartedAt)
		if err == nil {
			fmt.Printf(i18n.T("inspect.uptime", time.Since(startedAt).Round(time.Second)) + "\n")
		}
	}

	// 端口映射
	if len(config.HostConfig.PortBindings) > 0 {
		fmt.Printf(i18n.T("inspect.port_mappings") + "\n")
		for containerPort, hostBindings := range config.HostConfig.PortBindings {
			for _, hostBinding := range hostBindings {
				fmt.Printf(i18n.T("inspect.port_mapping", hostBinding.HostPort, containerPort) + "\n")
			}
		}
	}

	// 卷挂载
	if len(config.Mounts) > 0 {
		fmt.Printf(i18n.T("inspect.volume_mounts") + "\n")
		for _, mount := range config.Mounts {
			fmt.Printf(i18n.T("inspect.volume_mount", mount.Source, mount.Destination, mount.Type) + "\n")
		}
	}
}

// 实时监控容器统计信息
func monitorContainerStats(ctx context.Context, cli *client.Client, containerId string) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			stats, err := cli.ContainerStats(ctx, containerId, false)
			if err != nil {
				fmt.Printf(i18n.T("error.container_stats", err) + "\n")
				continue
			}

			var statsData types.StatsJSON
			decoder := json.NewDecoder(stats.Body)
			if err := decoder.Decode(&statsData); err != nil {
				stats.Body.Close()
				continue
			}
			stats.Body.Close()

			// 计算 CPU 使用率
			cpuPercent := calculateCPUPercent(&statsData)

			// 计算内存使用率
			memUsage := float64(statsData.MemoryStats.Usage)
			memLimit := float64(statsData.MemoryStats.Limit)
			memPercent := (memUsage / memLimit) * 100

			// 格式化内存大小
			memUsageMB := memUsage / 1024 / 1024
			memLimitMB := memLimit / 1024 / 1024

			// 网络 I/O
			var rxBytes, txBytes uint64
			for _, network := range statsData.Networks {
				rxBytes += network.RxBytes
				txBytes += network.TxBytes
			}

			// 清屏并显示统计信息
			fmt.Printf("\r" + i18n.T("inspect.stats_format",
				cpuPercent,
				memUsageMB, memLimitMB, memPercent,
				float64(rxBytes)/1024/1024, float64(txBytes)/1024/1024))
		}
	}
}

// 计算 CPU 使用率
func calculateCPUPercent(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

	if systemDelta > 0 && cpuDelta > 0 {
		return (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}

// 实时监控容器日志
func monitorContainerLogs(ctx context.Context, cli *client.Client, containerId string) {
	// 获取最近的日志
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "10", // 显示最近10行
		Timestamps: true,
	}

	logs, err := cli.ContainerLogs(ctx, containerId, options)
	if err != nil {
		fmt.Printf("\n" + i18n.T("error.container_logs", err) + "\n")
		return
	}
	defer logs.Close()

	fmt.Printf("\n" + i18n.T("inspect.recent_logs") + "\n")
	fmt.Println(strings.Repeat("-", 60))

	// 读取日志流
	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := logs.Read(buffer)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("\n" + i18n.T("error.read_logs", err) + "\n")
				}
				return
			}

			if n > 0 {
				// Docker 日志格式包含头部信息，需要跳过前8个字节
				logContent := string(buffer[8:n])
				fmt.Print(logContent)
			}
		}
	}
}

var inspectCommand = &cobra.Command{
	Use:   "inspect [container_name_or_id]",
	Short: i18n.T("inspect.short"),
	Long:  i18n.T("inspect.long"),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerName := args[0]
		inspectContainer(containerName)
	},
}

func init() {
	// 在命令执行前更新国际化文本
	inspectCommand.Short = i18n.T("inspect.short")
	inspectCommand.Long = i18n.T("inspect.long")

	rootCmd.AddCommand(inspectCommand)
}
