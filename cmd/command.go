package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
)

var containerId string

func init() {
	rootCmd.AddCommand(dockerCommand)
}

var dockerCommand = &cobra.Command{
	Use:     "command [container id]",
	Aliases: []string{"c"}, // 添加别名 v
	Short:   "Convert Docker container to docker run command",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 从 args 中获取 containerId
		containerId := args[0]
		// 获取容器配置
		config, err := getDockerContainerConfig(containerId)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		// 打印容器配置信息
		runCommand := generateRunCommand(config)
		fmt.Println(runCommand)
	},
}

// 获取容器的配置信息
func getDockerContainerConfig(containerID string) (*types.ContainerJSON, error) {
	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	// 获取容器的详细信息
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %v", err)
	}

	return &containerInfo, nil
}

// 转换为命令行
func generateRunCommand(config *types.ContainerJSON) string {
	// LogObject(config)
	var cmd strings.Builder

	cmd.WriteString("docker run")
	// 添加容器的名称
	if config.Name != "" {
		cmd.WriteString(fmt.Sprintf(" --name %s -d", strings.ReplaceAll(config.Name, "/", "")))
	}
	// 绑定端口映射
	for key, value := range config.HostConfig.PortBindings {
		if value[0].HostPort != "" {
			parts := strings.Split(string(key), "/")
			cmd.WriteString(fmt.Sprintf(" -p %s:%s", value[0].HostPort, parts[0]))
		}
	}
	// 绑定映射目录
	for _, mount := range config.Mounts {
		cmd.WriteString(fmt.Sprintf(" -v %s:%s", mount.Source, mount.Destination))
	}

	// 绑定 device 映射
	for _, device := range config.HostConfig.Devices {
		cmd.WriteString(fmt.Sprintf("--device %s:%s", device.PathOnHost, device.PathInContainer))
	}

	// Add CPU limit
	if config.HostConfig.NanoCPUs > 0 {
		cpuLimit := fmt.Sprintf("%f", float64(config.HostConfig.NanoCPUs)/1_000_000_000) // Convert NanoCPUs to CPU units
		cmd.WriteString(fmt.Sprintf(" --cpus=%s", cpuLimit))
	}

	// Add Memory limit
	if config.HostConfig.Memory > 0 {
		memoryLimit := units.BytesSize(float64(config.HostConfig.Memory))
		cmd.WriteString(fmt.Sprintf(" --memory=%s", memoryLimit))
	}

	// Add other options from config (e.g., restart policy, network)
	if config.HostConfig.RestartPolicy.Name != "" {
		cmd.WriteString(fmt.Sprintf(" --restart %s", config.HostConfig.RestartPolicy.Name))
	}

	if config.HostConfig.NetworkMode != "" {
		cmd.WriteString(fmt.Sprintf(" --network %s", config.HostConfig.NetworkMode))
	}

	// 设置容器的镜像
	if config.Config.Image != "" {
		cmd.WriteString(fmt.Sprintf(" %s", config.Config.Image))
	}

	return cmd.String()
}

func LogObject[T any](info T) {
	jsonData, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonData))
}
