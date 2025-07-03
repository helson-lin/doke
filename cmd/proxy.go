package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/helson-lin/doke/i18n"
	"github.com/spf13/cobra"
)

var registryMirrors = []string{"https://docker.1ms.run", "https://docker.1panel.live"}

func linuxDockerCheck() (bool, string) {
	_, err := exec.Command("docker", "version").Output()
	if err != nil {
		rootCmd.PrintErrf(i18n.T("error.docker_version_check", err))
		return false, i18n.T("error.docker_version_check_failed")
	}
	// 2. 检查docker是否正在运行
	dockerStatus, err := exec.Command("systemctl", "is-active", "docker").Output()
	if err != nil {
		rootCmd.PrintErrf(i18n.T("error.docker_status_check", err))
		return false, i18n.T("error.docker_status_check_failed")
	}

	if string(dockerStatus) == "active" {
		rootCmd.Println(i18n.T("proxy.docker_running"))
		return true, i18n.T("proxy.status_running")
	} else {
		rootCmd.Println(i18n.T("proxy.docker_not_running"))
		return true, i18n.T("proxy.status_stopped")
	}
}

func linuxProxy() {
	fmt.Println(i18n.T("proxy.linux_start_config"))

	// 1. 读取/etc/docker/daemon.json
	fmt.Println(i18n.T("proxy.reading_config"))
	daemonJson, err := os.ReadFile("/etc/docker/daemon.json")
	if err != nil {
		fmt.Println(i18n.T("error.read_daemon_json"))
		return
	}
	fmt.Println(i18n.T("proxy.read_config_success"))

	// 2. 解析daemon.json
	fmt.Println(i18n.T("proxy.parsing_config"))
	var daemonConfig map[string]interface{}
	err = json.Unmarshal(daemonJson, &daemonConfig)
	if err != nil {
		fmt.Println(i18n.T("error.parse_daemon_json"))
		return
	}
	fmt.Println(i18n.T("proxy.parse_config_success"))

	// 3. 检查是否存在registry-mirrors配置
	fmt.Println(i18n.T("proxy.checking_mirrors"))
	mirrors, ok := daemonConfig["registry-mirrors"]
	if !ok {
		fmt.Println(i18n.T("proxy.no_existing_mirrors"))
		daemonConfig["registry-mirrors"] = registryMirrors
	} else {
		// 4. 检查mirrors是否为空
		if mirrors == nil {
			fmt.Println(i18n.T("proxy.empty_mirrors"))
			daemonConfig["registry-mirrors"] = registryMirrors
		} else {
			// 5. 检查mirrors是否为数组
			mirrorsArray, ok := mirrors.([]interface{})
			if !ok {
				fmt.Println(i18n.T("proxy.invalid_mirrors_format"))
				daemonConfig["registry-mirrors"] = registryMirrors
			} else {
				// 6. 检查mirrorsArray是否为空
				if len(mirrorsArray) == 0 {
					fmt.Println(i18n.T("proxy.empty_mirrors_list"))
					daemonConfig["registry-mirrors"] = registryMirrors
				} else {
					fmt.Printf(i18n.T("proxy.existing_mirrors", mirrorsArray))

					// 7. 追加新的镜像源
					fmt.Printf(i18n.T("proxy.adding_mirrors", registryMirrors))
					for _, mirror := range registryMirrors {
						// 检查是否已存在
						exists := false
						for _, existing := range mirrorsArray {
							if existing == mirror {
								exists = true
								break
							}
						}
						if !exists {
							mirrorsArray = append(mirrorsArray, mirror)
							fmt.Printf(i18n.T("proxy.mirror_added", mirror))
						} else {
							fmt.Printf(i18n.T("proxy.mirror_exists", mirror))
						}
					}
					daemonConfig["registry-mirrors"] = mirrorsArray
					fmt.Printf(i18n.T("proxy.final_mirrors", mirrorsArray))
				}
			}
		}
	}

	// 8. 将daemon.json写回文件
	fmt.Println(i18n.T("proxy.saving_config"))
	daemonJson, err = json.MarshalIndent(daemonConfig, "", "  ")
	if err != nil {
		fmt.Println(i18n.T("error.serialize_config"))
		return
	}
	fmt.Println(i18n.T("proxy.serialize_success"))

	// 9. 将daemon.json写回文件
	err = os.WriteFile("/etc/docker/daemon.json", daemonJson, 0644)
	if err != nil {
		fmt.Println(i18n.T("error.write_daemon_json"))
		return
	}
	fmt.Println(i18n.T("proxy.write_config_success"))

	// 10. 执行 systemctl daemon-reload
	fmt.Println(i18n.T("proxy.reloading_daemon"))
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		fmt.Println(i18n.T("error.daemon_reload"))
		return
	}
	fmt.Println(i18n.T("proxy.daemon_reload_success"))

	// 11. 执行 systemctl restart docker
	fmt.Println(i18n.T("proxy.restarting_docker"))
	err = exec.Command("systemctl", "restart", "docker").Run()
	if err != nil {
		fmt.Println(i18n.T("error.restart_docker"))
		return
	}
	fmt.Println(i18n.T("proxy.restart_docker_success"))
	fmt.Println(i18n.T("proxy.linux_config_complete"))
}

func macDockerCheck() (bool, string) {
	// 1. 检查是否安装了orbStack
	_, err := exec.Command("orbctl", "version").Output()
	if err != nil {
		rootCmd.PrintErrln(i18n.T("error.orbstack_version_check"))
		return false, i18n.T("error.orbstack_version_check_failed")
	}
	// 2. 检查orbStack状态
	orbStack, err := exec.Command("orbctl", "status").Output()
	if err != nil {
		return true, i18n.T("proxy.status_stopped")
	}
	if string(orbStack) == "Running" {
		rootCmd.Println(i18n.T("proxy.orbstack_running"))
		return true, string(orbStack)
	} else if string(orbStack) == "Stopped" {
		rootCmd.Println(i18n.T("proxy.orbstack_stopped"))
		return true, string(orbStack)
	}
	return true, string(orbStack)
}

func updateObrStack() {
	fmt.Println(i18n.T("proxy.start_config"))

	// 1. 获取用户主目录并构建正确的路径
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(i18n.T("error.get_home_dir", err))
		return
	}
	dockerConfigPath := filepath.Join(homeDir, ".orbstack", "config", "docker.json")
	fmt.Printf(i18n.T("proxy.config_path", dockerConfigPath))

	// 2. 读取~/.orbstack/config/docker.json
	dockerJson, err := os.ReadFile(dockerConfigPath)
	if err != nil {
		fmt.Printf(i18n.T("error.read_config_file", err))
		return
	}
	fmt.Println(i18n.T("proxy.read_config_success"))

	// 3. 解析现有的配置
	var config map[string]interface{}
	err = json.Unmarshal(dockerJson, &config)
	if err != nil {
		fmt.Printf(i18n.T("error.parse_config_file", err))
		return
	}
	fmt.Println(i18n.T("proxy.parse_config_success"))

	// 4. 获取现有的 registry-mirrors 配置
	existingMirrors, ok := config["registry-mirrors"]
	var mirrorsArray []interface{}

	if ok && existingMirrors != nil {
		// 如果存在现有的镜像配置，转换为数组
		if mirrors, ok := existingMirrors.([]interface{}); ok {
			mirrorsArray = mirrors
			fmt.Printf(i18n.T("proxy.existing_mirrors", mirrorsArray))
		} else {
			fmt.Println(i18n.T("proxy.invalid_mirrors_format"))
			mirrorsArray = []interface{}{}
		}
	} else {
		fmt.Println(i18n.T("proxy.no_existing_mirrors"))
		mirrorsArray = []interface{}{}
	}

	// 5. 追加新的镜像源
	fmt.Printf(i18n.T("proxy.adding_mirrors", registryMirrors))
	for _, mirror := range registryMirrors {
		// 检查是否已存在
		exists := false
		for _, existing := range mirrorsArray {
			if existing == mirror {
				exists = true
				break
			}
		}
		if !exists {
			mirrorsArray = append(mirrorsArray, mirror)
			fmt.Printf(i18n.T("proxy.mirror_added", mirror))
		} else {
			fmt.Printf(i18n.T("proxy.mirror_exists", mirror))
		}
	}

	// 6. 更新配置
	config["registry-mirrors"] = mirrorsArray
	fmt.Printf(i18n.T("proxy.final_mirrors", mirrorsArray))

	// 7. 将修改后的配置转回JSON
	updatedJson, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println(i18n.T("error.serialize_config", err))
		return
	}
	fmt.Println(i18n.T("proxy.serialize_success"))

	// 8. 写回~/.orbstack/config/docker.json
	err = os.WriteFile(dockerConfigPath, updatedJson, 0644)
	if err != nil {
		fmt.Printf(i18n.T("error.write_config_file", err))
		return
	}
	fmt.Println(i18n.T("proxy.write_config_success"))

	// 9. 重启orbStack
	fmt.Println(i18n.T("proxy.restarting_docker"))
	err = exec.Command("orbctl", "restart", "docker").Run()
	if err != nil {
		fmt.Println(i18n.T("error.restart_docker", err))
		return
	}
	fmt.Println(i18n.T("proxy.restart_docker_success"))
	fmt.Println(i18n.T("proxy.config_complete"))
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: i18n.T("proxy.short"),
	Long:  i18n.T("proxy.long"),
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 判断操作系统
		os := runtime.GOOS
		if os == "linux" {
			isInstall, _ := linuxDockerCheck()
			if isInstall {
				linuxProxy()
			} else {
				rootCmd.PrintErrln(i18n.T("error.docker_not_installed"))
			}
		} else if os == "darwin" {
			isInstall, _ := macDockerCheck()
			if isInstall {
				updateObrStack()
			} else {
				rootCmd.PrintErrln(i18n.T("error.orbstack_not_installed"))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
