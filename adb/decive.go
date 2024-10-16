package adb

import (
	"fmt"
	"os/exec"
	"strings"
)

const AdbPath = "adb/adb.exe"

func CheckDeviceConnected() bool {
	cmd := exec.Command(AdbPath, "devices")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("检查设备时出错:", err)
		return false
	}
	fmt.Println("命令输出:\n", string(output)) // 显示输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "device") && !strings.Contains(line, "List of devices attached") {
			return true
		}
	}
	return false
}

//func CheckAndKillPort() {
//	cmd := exec.Command("netstat", "-ano")
//	output, err := cmd.CombinedOutput() // 使用 CombinedOutput 捕获 stdout 和 stderr
//	if err != nil {
//		fmt.Println("检查端口时出错:", err)
//		return
//	}
//	fmt.Println("命令输出:\n", string(output)) // 显示输出
//	lines := strings.Split(string(output), "\n")
//	currentPID := fmt.Sprintf("%d", os.Getpid()) // 获取当前进程的 PID
//	for _, line := range lines {
//		if strings.Contains(line, ":5037") {
//			fields := strings.Fields(line)
//			if len(fields) > 4 {
//				pid := fields[len(fields)-1]
//				if pid != "0" && pid != currentPID { // 确保 PID 有效且不是当前进程
//					killCmd := exec.Command("taskkill", "/PID", pid, "/F")
//					err := killCmd.Run()
//					if err != nil {
//						fmt.Printf("无法结束进程 %s: %v\n", pid, err)
//					} else {
//						fmt.Printf("端口 5037 被进程 ID %s 占用，已结束进程\n", pid)
//					}
//				}
//			}
//		}
//	}
//}

func GrantPermissions() {
	cmd := exec.Command(AdbPath, "shell", "mount -o remount,rw /")
	output, err := cmd.CombinedOutput() // 使用 CombinedOutput 捕获 stdout 和 stderr
	if err != nil {
		fmt.Println("赋予权限时出错:", err)
	} else {
		fmt.Println("命令输出:", string(output)) // 显示输出
		fmt.Println("已赋予设备写权限")
	}
}
