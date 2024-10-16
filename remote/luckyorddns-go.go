package remote

import (
	"R200-tool/adb"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// 清空控制台屏幕
func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	default:
		fmt.Println("Unsupported platform")
	}
}

// 获取用户输入
func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}

func EnableIPv6Penetration(callback func()) {
	fmt.Println("请选择ipv6穿透工具:")
	fmt.Println("1. lucky")
	fmt.Println("2. ddns-go")

	var option string
	for {
		option = getInput("输入选项编号并按回车: ")
		if option == "1" || option == "2" {
			break
		}
		fmt.Println("无效的选择，请重新输入")
	}

	switch option {
	case "1":
		fmt.Println("配置lucky...")
		cmd := exec.Command(adb.AdbPath, "push", "./remote/lucky", "/cache/")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error pushing lucky: %v\n", err)
		} else if len(output) == 0 {
			fmt.Println("已推送lucky文件到设备")
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}

		cmd = exec.Command(adb.AdbPath, "push", "./remote/lucky.service", "/lib/systemd/system/")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error pushing lucky service: %v\n", err)
		} else if len(output) == 0 {
			fmt.Println("已推送lucky服务文件到设备")
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}

		cmd = exec.Command(adb.AdbPath, "shell", "systemctl start lucky")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error starting lucky: %v\n", err)
		} else if len(output) == 0 {
			fmt.Println("lucky已启动")
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}

		// 提示用户选择是否开启自启
		enableAutostart := getInput("是否开启lucky自启? 输入1开启: ")
		if enableAutostart == "1" {
			cmd = exec.Command(adb.AdbPath, "shell", "systemctl enable lucky")
			output, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error enabling lucky autostart: %v\n", err)
			} else if len(output) == 0 {
				fmt.Println("lucky已设置为开机自启")
			} else {
				fmt.Printf("Output: %s\n", output) // 打印输出
			}
		} else {
			fmt.Println("未开启lucky自启")
		}

	case "2":
		fmt.Println("配置ddns-go...")
		// 推送ddns-go文件到设备
		cmd := exec.Command(adb.AdbPath, "push", "./remote/ddnsgo", "/cache/")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error pushing ddns-go: %v\n", err)
		} else if len(output) == 0 {
			fmt.Println("已推送ddns-go文件到设备")
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}

		cmd = exec.Command(adb.AdbPath, "shell", "chmod +x /cache/ddnsgo/ddns-go")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error setting permissions for ddns-go: %v\n", err)
		} else if len(output) == 0 {
			fmt.Println("已设置ddns-go程序权限")
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}

		// 检查SELINUX模式
		cmd = exec.Command(adb.AdbPath, "shell", "getenforce")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error checking SELINUX status: %v\n", err)
		} else {
			selinuxStatus := strings.TrimSpace(string(output))
			if selinuxStatus == "Enforcing" {
				fmt.Println("SELINUX模式当前为Enforcing")
				disableSelinux := getInput("是否关闭SELINUX模式并重启CPE? 输入1确认: ")
				if disableSelinux == "1" {
					cmd = exec.Command(adb.AdbPath, "shell", "setenforce 0")
					output, err = cmd.CombinedOutput()
					if err != nil {
						fmt.Printf("Error disabling SELINUX: %v\n", err)
					} else {
						cmd = exec.Command(adb.AdbPath, "shell", "sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config")
						output, err = cmd.CombinedOutput()
						if err != nil {
							fmt.Printf("Error updating SELINUX config: %v\n", err)
						} else {
							fmt.Println("SELINUX模式已关闭，请重启CPE以应用更改")
							cmd = exec.Command(adb.AdbPath, "shell", "reboot")
							err := cmd.Run()
							if err != nil {
								return
							}
							fmt.Println("设备重启后，请重新运行脚本以完成ddns-go安装服务")
							return
						}
					}
				} else {
					fmt.Println("未关闭SELINUX模式，ddns-go自启可能无法正常工作")
				}
			} else {
				fmt.Println("SELINUX模式当前已关闭")
				// 运行ddns-go安装服务
				cmd = exec.Command(adb.AdbPath, "shell", "/cache/ddnsgo/ddns-go -s install")
				output, err = cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("Error installing ddns-go service: %v\n", err)
				} else if len(output) == 0 {
					fmt.Println("ddns-go已经运行")
				} else {
					fmt.Printf("Output: %s\n", output) // 打印输出
				}
				// 等待用户输入回车
				fmt.Println("按回车键返回菜单...")
				reader := bufio.NewReader(os.Stdin)
				_, _ = reader.ReadString('\n')

				// 清空屏幕并返回主菜单
				clearScreen()
				callback()
			}
		}

	default:
		fmt.Println("无效的选择")
	}
}
