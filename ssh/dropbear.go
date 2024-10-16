package ssh

import (
	"R200-tool/adb"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

func EnableSSH(callback func()) {
	fmt.Println("开启SSH功能")

	// 赋予设备写权限
	cmd := exec.Command(adb.AdbPath, "shell", "mount -o remount,rw /")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error granting permissions: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已赋予设备写权限")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 检查并创建必要的目录
	directories := []string{
		"/usr/sbin",
		"/usr/bin",
		"/usr/libexec",
		"/home/root/r200",
	}
	for _, dir := range directories {
		cmd = exec.Command(adb.AdbPath, "shell", "mkdir -p", dir)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
		} else if len(output) == 0 {
			fmt.Printf("已检查并创建目录 %s\n", dir)
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}
	}
	fmt.Println("已检查并创建必要的目录")

	// 推送压缩包到设备中
	cmd = exec.Command(adb.AdbPath, "push", "./ssh/dropbear.zip", "/home/root/r200")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error pushing zip file: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("压缩包已推送到设备")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 解压压缩包到 r200 目录
	cmd = exec.Command(adb.AdbPath, "shell", "unzip /home/root/r200/dropbear.zip -d /home/root/r200")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error unzipping file: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("压缩包已解压到 r200 目录")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 部署文件到对应文件夹
	cmds := []struct {
		src  string
		dest string
	}{
		{"/home/root/r200/dropbear/sbin/*", "/usr/sbin"},
		{"/home/root/r200/dropbear/bin/*", "/usr/bin"},
		{"/home/root/r200/dropbear/libexec/*", "/usr/libexec"},
	}
	for _, c := range cmds {
		cmd = exec.Command(adb.AdbPath, "shell", "cp", c.src, c.dest)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error copying files from %s to %s: %v\n", c.src, c.dest, err)
		} else if len(output) == 0 {
			fmt.Printf("已将文件从 %s 部署到 %s\n", c.src, c.dest)
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}
	}
	fmt.Println("已部署文件到目标设备")

	// 修改权限
	permissions := []string{
		"/usr/bin/scp",
		"/usr/bin/dbclient",
		"/usr/bin/dropbearconvert",
		"/usr/bin/dropbearkey",
		"/usr/sbin/dropbear",
		"/usr/libexec/sftp-server",
	}
	for _, file := range permissions {
		cmd = exec.Command(adb.AdbPath, "shell", "chmod 0777", file)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error setting permissions for %s: %v\n", file, err)
		} else if len(output) == 0 {
			fmt.Printf("已设置 %s 的权限\n", file)
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}
	}
	fmt.Println("已经设置好权限")

	// 生成密钥
	cmd = exec.Command(adb.AdbPath, "shell", "mkdir -p /etc/dropbear && dropbearkey -t rsa -f /etc/dropbear/dropbear_rsa_host_key")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("密钥生成成功")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 修改 root 用户密码
	cmd = exec.Command(adb.AdbPath, "shell", "echo -e 'online\nonline' | passwd root")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error changing root password: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("root 用户密码已修改为 online")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 临时运行 dropbear
	//cmd = exec.Command(adb.AdbPath, "shell", "dropbear")
	//output, err = cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Printf("Error starting dropbear: %v\n", err)
	//} else if len(output) == 0 {
	//	fmt.Println("已临时运行 dropbear，请使用其他应用测试是否能够通过 SSH 链接设备")
	//} else {
	//	fmt.Printf("Output: %s\n", output) // 打印输出
	//}

	// 推送服务文件到目标设备
	cmd = exec.Command(adb.AdbPath, "push", "./ssh/dropbear.service", "/lib/systemd/system/")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error pushing service file: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("自启动服务文件已推送")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 启动服务
	cmd = exec.Command(adb.AdbPath, "shell", "systemctl start dropbear")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error starting dropbear service: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("dropbear 服务已启动，请使用其他应用测试是否能够通过 SSH 链接设备")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 检查 dropbear 是否在运行
	cmd = exec.Command(adb.AdbPath, "shell", "ps | grep dropbear")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error checking dropbear process: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("dropbear 未能启动，请手动检查")
		return
	} else {
		fmt.Printf("dropbear 进程信息: %s\n", output) // 打印输出
	}

	// 等待用户确认 SSH 连接成功
	var sshSuccess int
	fmt.Println("通过 SSH 链接成功请按 1")
	_, err = fmt.Scan(&sshSuccess)
	if err != nil {
		return
	}
	if sshSuccess == 1 {
		fmt.Println("正在进行设置 dropbear 开机自启")
		cmd = exec.Command(adb.AdbPath, "shell", "systemctl enable dropbear")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error enabling dropbear service: %v\n", err)
		} else if len(output) == 0 {
			fmt.Println("dropbear 服务已设置为开机自启，请重启设备进行 SSH 链接测试开机自启是否成功")
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}

		// 删除 /home/root/r200/ 下的 dropbear 文件夹和压缩包
		cmd = exec.Command(adb.AdbPath, "shell", "rm -rf /home/root/r200/dropbear /home/root/r200/dropbear.zip")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error deleting dropbear files: %v\n", err)
		} else if len(output) == 0 {
			fmt.Println("已删除 /home/root/r200/ 下的 dropbear 文件夹和压缩包")
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
	} else {
		fmt.Println("设置 SSH 链接失败，请自行检查原因")
	}
}
