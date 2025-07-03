package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/helson-lin/doke/i18n"
	"github.com/spf13/cobra"
)

// 清理统计信息
type CleanupStats struct {
	ContainersRemoved int
	ImagesRemoved     int
	NetworksRemoved   int
	VolumesRemoved    int
	SpaceReclaimed    uint64 // 回收的空间（字节）
}

// 执行清理操作
func performCleanup(all bool, force bool) {
	fmt.Println(i18n.T("clear.starting"))
	fmt.Println(strings.Repeat("=", 60))

	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf(i18n.T("error.docker_client", err))
		fmt.Println()
		return
	}
	defer cli.Close()

	ctx := context.Background()
	stats := &CleanupStats{}

	// 显示当前资源使用情况
	showCurrentUsage(ctx, cli)

	if !force {
		if !confirmCleanup(all) {
			fmt.Println(i18n.T("clear.cancelled"))
			return
		}
	}

	fmt.Println("\n" + i18n.T("clear.starting_cleanup"))

	// 1. 清理停止的容器
	cleanupContainers(ctx, cli, stats, all)

	// 2. 清理未使用的镜像
	cleanupImages(ctx, cli, stats, all)

	// 3. 清理未使用的网络
	cleanupNetworks(ctx, cli, stats)

	// 4. 清理未使用的卷
	cleanupVolumes(ctx, cli, stats)

	// 显示清理结果
	showCleanupResults(stats)
}

// 显示当前资源使用情况
func showCurrentUsage(ctx context.Context, cli *client.Client) {
	fmt.Println(i18n.T("clear.current_usage"))

	// 容器统计
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err == nil {
		running := 0
		stopped := 0
		for _, c := range containers {
			if c.State == "running" {
				running++
			} else {
				stopped++
			}
		}
		fmt.Printf("   " + i18n.T("clear.containers", len(containers), running, stopped) + "\n")
	}

	// 镜像统计
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err == nil {
		var totalSize int64
		danglingCount := 0
		for _, img := range images {
			totalSize += img.Size
			if len(img.RepoTags) == 0 || (len(img.RepoTags) == 1 && img.RepoTags[0] == "<none>:<none>") {
				danglingCount++
			}
		}
		fmt.Printf("   " + i18n.T("clear.images", len(images), danglingCount, float64(totalSize)/1024/1024/1024) + "\n")
	}

	// 网络统计
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err == nil {
		customNetworks := 0
		for _, net := range networks {
			if !isSystemNetwork(net.Name) {
				customNetworks++
			}
		}
		fmt.Printf("   " + i18n.T("clear.networks", len(networks), customNetworks) + "\n")
	}

	// 卷统计
	volumesList, err := cli.VolumeList(ctx, volume.ListOptions{})
	if err == nil {
		fmt.Printf("   " + i18n.T("clear.volumes", len(volumesList.Volumes)) + "\n")
	}

	fmt.Println()
}

// 确认清理操作
func confirmCleanup(all bool) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(i18n.T("clear.confirm_title"))
	fmt.Println("   " + i18n.T("clear.confirm_containers"))
	if all {
		fmt.Println("   " + i18n.T("clear.confirm_images_all"))
	} else {
		fmt.Println("   " + i18n.T("clear.confirm_images_dangling"))
	}
	fmt.Println("   " + i18n.T("clear.confirm_networks"))
	fmt.Println("   " + i18n.T("clear.confirm_volumes"))

	fmt.Print("\n" + i18n.T("clear.confirm_prompt"))
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// 清理停止的容器
func cleanupContainers(ctx context.Context, cli *client.Client, stats *CleanupStats, all bool) {
	fmt.Print(i18n.T("clear.cleaning_containers"))

	// 获取停止的容器
	containerFilters := filters.NewArgs()
	containerFilters.Add("status", "exited")
	containerFilters.Add("status", "created")

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: containerFilters,
	})

	if err != nil {
		fmt.Printf(i18n.T("error.container_list", err))
		fmt.Println()
		return
	}

	for _, c := range containers {
		err := cli.ContainerRemove(ctx, c.ID, container.RemoveOptions{})
		if err != nil {
			fmt.Printf(i18n.T("error.container_remove", c.ID[:12], err))
			fmt.Println()
		} else {
			stats.ContainersRemoved++
		}
	}

	fmt.Printf(i18n.T("clear.containers_removed", stats.ContainersRemoved))
	fmt.Println()
}

// 清理未使用的镜像
func cleanupImages(ctx context.Context, cli *client.Client, stats *CleanupStats, all bool) {
	if all {
		fmt.Print(i18n.T("clear.cleaning_images_all"))
	} else {
		fmt.Print(i18n.T("clear.cleaning_images_dangling"))
	}

	pruneFilters := filters.NewArgs()
	if !all {
		pruneFilters.Add("dangling", "true")
	}

	report, err := cli.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		fmt.Printf(i18n.T("error.image_prune", err))
		fmt.Println()
		return
	}

	stats.ImagesRemoved = len(report.ImagesDeleted)
	stats.SpaceReclaimed += report.SpaceReclaimed

	fmt.Printf(i18n.T("clear.images_removed", stats.ImagesRemoved, float64(report.SpaceReclaimed)/1024/1024))
	fmt.Println()
}

// 清理未使用的网络
func cleanupNetworks(ctx context.Context, cli *client.Client, stats *CleanupStats) {
	fmt.Print(i18n.T("clear.cleaning_networks"))

	report, err := cli.NetworksPrune(ctx, filters.NewArgs())
	if err != nil {
		fmt.Printf(i18n.T("error.network_prune", err))
		fmt.Println()
		return
	}

	stats.NetworksRemoved = len(report.NetworksDeleted)
	fmt.Printf(i18n.T("clear.networks_removed", stats.NetworksRemoved))
	fmt.Println()
}

// 清理未使用的卷
func cleanupVolumes(ctx context.Context, cli *client.Client, stats *CleanupStats) {
	fmt.Print(i18n.T("clear.cleaning_volumes"))

	report, err := cli.VolumesPrune(ctx, filters.NewArgs())
	if err != nil {
		fmt.Printf(i18n.T("error.volume_prune", err))
		fmt.Println()
		return
	}

	stats.VolumesRemoved = len(report.VolumesDeleted)
	stats.SpaceReclaimed += report.SpaceReclaimed

	fmt.Printf(i18n.T("clear.volumes_removed", stats.VolumesRemoved, float64(report.SpaceReclaimed)/1024/1024))
	fmt.Println()
}

// 判断是否为系统网络
func isSystemNetwork(name string) bool {
	systemNetworks := []string{"bridge", "host", "none"}
	for _, sysNet := range systemNetworks {
		if name == sysNet {
			return true
		}
	}
	return false
}

// 显示清理结果
func showCleanupResults(stats *CleanupStats) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(i18n.T("clear.results_title"))
	fmt.Printf(i18n.T("clear.results_stats") + "\n")
	fmt.Printf("   " + i18n.T("clear.results_containers", stats.ContainersRemoved) + "\n")
	fmt.Printf("   " + i18n.T("clear.results_images", stats.ImagesRemoved) + "\n")
	fmt.Printf("   " + i18n.T("clear.results_networks", stats.NetworksRemoved) + "\n")
	fmt.Printf("   " + i18n.T("clear.results_volumes", stats.VolumesRemoved) + "\n")
	fmt.Printf("   " + i18n.T("clear.results_space", float64(stats.SpaceReclaimed)/1024/1024) + "\n")

	if stats.ContainersRemoved+stats.ImagesRemoved+stats.NetworksRemoved+stats.VolumesRemoved == 0 {
		fmt.Println(i18n.T("clear.already_clean"))
	}
}

// 显示详细的系统信息
func showSystemInfo() {
	fmt.Println(i18n.T("sysinfo.getting"))

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf(i18n.T("error.docker_client", err))
		fmt.Println()
		return
	}
	defer cli.Close()

	ctx := context.Background()

	// 获取系统信息
	info, err := cli.Info(ctx)
	if err != nil {
		fmt.Printf(i18n.T("error.system_info", err))
		fmt.Println()
		return
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf(i18n.T("sysinfo.title") + "\n")
	fmt.Printf("   " + i18n.T("sysinfo.version", info.ServerVersion) + "\n")
	fmt.Printf("   " + i18n.T("sysinfo.containers_total", info.Containers, info.ContainersRunning, info.ContainersPaused, info.ContainersStopped) + "\n")
	fmt.Printf("   " + i18n.T("sysinfo.images_total", info.Images) + "\n")
	fmt.Printf("   " + i18n.T("sysinfo.storage_driver", info.Driver) + "\n")
	fmt.Printf("   " + i18n.T("sysinfo.root_dir", info.DockerRootDir) + "\n")

	// 获取磁盘使用情况
	diskUsage, err := cli.DiskUsage(ctx, types.DiskUsageOptions{})
	if err == nil {
		var totalSize int64
		totalSize += diskUsage.LayersSize
		for _, img := range diskUsage.Images {
			totalSize += img.Size
		}
		for _, vol := range diskUsage.Volumes {
			totalSize += vol.UsageData.Size
		}

		fmt.Printf("   " + i18n.T("sysinfo.disk_usage", float64(totalSize)/1024/1024/1024) + "\n")
		fmt.Printf("     " + i18n.T("sysinfo.disk_layers", float64(diskUsage.LayersSize)/1024/1024/1024) + "\n")
		fmt.Printf("     " + i18n.T("sysinfo.disk_cache", float64(diskUsage.BuilderSize)/1024/1024/1024) + "\n")
	}

	fmt.Printf("   " + i18n.T("sysinfo.system_time", time.Now().Format("2006-01-02 15:04:05")) + "\n")
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: i18n.T("clear.short"),
	Long:  i18n.T("clear.long"),
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		force, _ := cmd.Flags().GetBool("force")
		info, _ := cmd.Flags().GetBool("info")

		if info {
			showSystemInfo()
		} else {
			performCleanup(all, force)
		}
	},
}

func init() {
	clearCmd.Flags().BoolP("all", "a", false, i18n.T("clear.flag.all"))
	clearCmd.Flags().BoolP("force", "f", false, i18n.T("clear.flag.force"))
	clearCmd.Flags().BoolP("info", "i", false, i18n.T("clear.flag.info"))
	rootCmd.AddCommand(clearCmd)
}
