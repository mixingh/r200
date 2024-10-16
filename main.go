package main

import (
	"R200-tool/adb"
	"R200-tool/forward"
	"R200-tool/remote"
	"R200-tool/ssh"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("启动程序...")
	if !adb.CheckDeviceConnected() {
		fmt.Println("没有设备连接，请检查5037端口是否被其他程序占用。")
		// 添加等待用户输入，防止窗口立即关闭
		fmt.Println("按回车键退出...")
		_, err := fmt.Scanln()
		if err != nil {
			fmt.Println("等待用户输入时出错:", err)
			return
		}
		return
	} else {
		fmt.Println("设备已连接")
		adb.GrantPermissions()
		checkAndDisableSELinux()
		showMenu()
	}
}

func checkAndDisableSELinux() {
	cmd := exec.Command(adb.AdbPath, "shell", "getenforce")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("检查SELINUX状态时出错: %v\n", err)
		return
	}
	selinuxStatus := strings.TrimSpace(string(output))
	if selinuxStatus == "Enforcing" {
		fmt.Println("SELINUX模式当前为Enforcing")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("是否关闭SELINUX模式并重启设备? 输入1确认，输入其他内容退出程序: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "1" {
			cmd = exec.Command(adb.AdbPath, "shell", "setenforce 0")
			output, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("关闭SELINUX时出错: %v\n", err)
				return
			}
			cmd = exec.Command(adb.AdbPath, "shell", "sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config")
			output, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("更新SELINUX配置时出错: %v\n", err)
				return
			}
			fmt.Println("SELINUX模式已关闭，请重启设备以应用更改")
			fmt.Println("按回车键重启设备...设备重启后，请重新运行脚本以完成配置")
			_, _ = reader.ReadString('\n') // 等待用户按回车键
			cmd = exec.Command(adb.AdbPath, "shell", "reboot")
			err := cmd.Run()
			if err != nil {
				return
			}
		} else {
			fmt.Println("未关闭SELINUX模式，程序将退出")
			os.Exit(0)
		}
	} else {
		fmt.Println("SELINUX模式当前已关闭")
	}
}

func showMenu() {
	fmt.Println("请选择一个功能:")
	fmt.Println("0. 退出脚本")
	fmt.Println("1. 开启ssh(dropbear)")
	fmt.Println("2. 开启短信转发(forward)")
	fmt.Println("3. 开启IPv6穿透(lucky or ddns-go)")
	fmt.Println("4. 开启IPv4穿透(nps客户端)")
	// 其他功能选项...

	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("读取用户输入时出错:", err)
		return
	}

	switch choice {
	case 0:
		fmt.Println("拜拜")
	case 1:
		ssh.EnableSSH(showMenu)
		//showMenu()

	case 2:
		forward.EnableForward(showMenu)

	case 3:
		remote.EnableIPv6Penetration(showMenu)

	case 4:
		remote.EnableNPS(showMenu)
	// 其他功能实现...
	default:
		fmt.Println("无效的选择")
	}
}
