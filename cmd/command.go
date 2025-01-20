package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"log"
	"strings"
	"gopkg.in/yaml.v3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
)

var containerId string
var isCompose bool = false

type Service struct {
	Image         string            `yaml:"image"`
	ContainerName string            `yaml:"container_name,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Environment   map[string]string `yaml:"environment,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	Command       string            `yaml:"command,omitempty"`
	Networks      []string          `yaml:"networks,omitempty"`
	HealthCheck   *HealthCheck      `yaml:"healthcheck,omitempty"`
}

type HealthCheck struct {
	Test        []string `yaml:"test"`
	Interval    string   `yaml:"interval,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"start_period,omitempty"`
}

type DockerCompose struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]struct{} `yaml:"networks,omitempty"`
}


func init() {
	dockerCommand.PersistentFlags().BoolVarP(&isCompose, "json", "j", false, "is export docker compose yaml file")
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
		if isCompose {
			fmt.Println("ok")
			yamlData, err := getDockerComposeYaml(config)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			if yamlData != "" {
				// fmt.Println(yamlData)
				err := writeDockerComposeYaml(config.Name, yamlData)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
			}
		} else {
			// 打印容器配置信息
			runCommand := generateRunCommand(config)
			fmt.Println(runCommand)
		}
	},
}

// 生成 Docker Compose YAML 文件并写入
func writeDockerComposeYaml(containerName string, yamlData string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	fileName := fmt.Sprintf("%s.yaml", containerName)
	filePath := filepath.Join(currentDir, fileName)

	fmt.Printf("是否将 Docker Compose 配置写入文件 %s？(y/n): ", filePath)

	var confirm string // Declare confirm outside the if block
	_, err = fmt.Scanln(&confirm) // Use _ to ignore the return value of a
	if err != nil {
		return fmt.Errorf("读取用户输入失败: %v", err)
	}

	if strings.ToLower(confirm) != "y" {
		fmt.Println("用户取消操作。")
		return nil
	}

	err = os.WriteFile(filePath, []byte(yamlData), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("Docker Compose 配置已成功写入文件: %s\n", fileName)
	return nil
}

// 生成 docker compose yaml 文件
func getDockerComposeYaml(config *types.ContainerJSON) (string, error) {
	// 解析容器配置
	container := config.Config
	hostConfig := config.HostConfig

	// 构建 Service
	service := Service{
		Image:       container.Image,
		ContainerName: strings.TrimPrefix(config.Name, "/"), // 去掉容器名前缀的 "/"
		Command:     strings.Join(container.Cmd, " "),
	}

	// 解析端口映射
	for port, bindings := range hostConfig.PortBindings {
		for _, binding := range bindings {
			service.Ports = append(service.Ports, fmt.Sprintf("%s:%s", binding.HostPort, port))
		}
	}

	// 解析环境变量
	service.Environment = make(map[string]string)
	for _, env := range container.Env {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			service.Environment[parts[0]] = parts[1]
		}
	}

	// 解析挂载卷
	for _, mount := range hostConfig.Mounts {
		service.Volumes = append(service.Volumes, fmt.Sprintf("%s:%s", mount.Source, mount.Target))
	}

	// 解析网络
	if hostConfig.NetworkMode != "" {
		service.Networks = []string{string(hostConfig.NetworkMode)}
	}

	// 解析健康检查
	if container.Healthcheck != nil {
		service.HealthCheck = &HealthCheck{
			Test:        container.Healthcheck.Test,
			Interval:    fmt.Sprintf("%dns", container.Healthcheck.Interval),
			Timeout:     fmt.Sprintf("%dns", container.Healthcheck.Timeout),
			Retries:     container.Healthcheck.Retries,
			StartPeriod: fmt.Sprintf("%dns", container.Healthcheck.StartPeriod),
		}
	}

	// 构建 Docker Compose 配置
	compose := DockerCompose{
		Version: "3.8",
		Services: map[string]Service{
			"app": service,
		},
	}

	// 添加网络配置
	if len(service.Networks) > 0 {
		compose.Networks = make(map[string]struct{})
		for _, network := range service.Networks {
			compose.Networks[network] = struct{}{}
		}
	}

	// 转换为 YAML
	yamlData, err := yaml.Marshal(&compose)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %v", err)
	}

	return string(yamlData), nil
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
