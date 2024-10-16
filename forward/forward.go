package forward

import (
	"R200-tool/adb"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
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
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("输入不能为空，请重新输入")
		return getInput(prompt)
	}
	return input
}

func EnableForward(callback func()) {
	fmt.Println("开启Forward功能")

	// 推送 curl 程序并设置权限
	cmd := exec.Command(adb.AdbPath, "push", "./forward/curl", "/usr/bin/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error pushing curl: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已推送curl程序到CPE中")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	cmd = exec.Command(adb.AdbPath, "shell", "chmod +x /usr/bin/curl")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting permissions for curl: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已设置curl程序权限")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 推送 forward 程序并设置权限
	cmd = exec.Command(adb.AdbPath, "push", "./forward/forward", "/usr/bin/")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error pushing forward: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已推送forward程序到CPE中")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	cmd = exec.Command(adb.AdbPath, "shell", "chmod +x /usr/bin/forward")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting permissions for forward: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已设置forward程序权限")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 提示用户选择推送方式
	var pushType string
	for {
		fmt.Println("请选择推送方式:")
		fmt.Println("1. pushplus")
		fmt.Println("2. 钉钉机器人")
		fmt.Println("3. wxpush")
		fmt.Println("4. bark")
		fmt.Println("5. gotify")
		pushType = getInput("输入推送方式编号并按回车: ")
		if pushType == "1" || pushType == "2" || pushType == "3" || pushType == "4" || pushType == "5" {
			break
		}
		fmt.Println("无效的选择，请重新输入")
	}

	// 根据选择进行配置
	var configContent string
	switch pushType {
	case "1":
		token := getInput("请输入pushplus的token: ")
		configContent = fmt.Sprintf("#推送方式配置：1：pushplus；2：钉钉机器人；3：wxpush; 4:bark; 5: gotify\npush_type:%s\n#pushplus配置\ntoken:%s\n", pushType, token)
	case "2":
		dingtalk := getInput("请输入钉钉机器人的token: ")
		secret := getInput("请输入钉钉机器人的secret: ")
		configContent = fmt.Sprintf("#推送方式配置：1：pushplus；2：钉钉机器人；3：wxpush; 4:bark; 5: gotify\npush_type:%s\n#钉钉机器人配置\ndingtalk:%s\nsecret:%s\n", pushType, dingtalk, secret)
	case "3":
		wxpush := getInput("请输入wxpush的wxpush_appToken: ")
		wxpushUids := getInput("请输入wxpush的wxpush_uids: ")
		configContent = fmt.Sprintf("#推送方式配置：1：pushplus；2：钉钉机器人；3：wxpush; 4:bark; 5: gotify\npush_type:%s\n#wxpush配置\nwxpush:%s\nwxpush_uids:%s\n", pushType, wxpush, wxpushUids)
	case "4":
		barkKey := getInput("请输入bark的bark_key: ")
		configContent = fmt.Sprintf("#推送方式配置：1：pushplus；2：钉钉机器人；3：wxpush; 4:bark; 5: gotify\npush_type:%s\n#bark配置\nbark_key:%s\n", pushType, barkKey)
	case "5":
		gotifyToken := getInput("请输入gotify的gotify_token: ")
		gotifyURL := getInput("请输入gotify的gotify_url: ")
		configContent = fmt.Sprintf("#推送方式配置：1：pushplus；2：钉钉机器人；3：wxpush; 4:bark; 5: gotify\npush_type:%s\n#gotify配置\ngotify_token:%s\ngotify_url:%s\n", pushType, gotifyToken, gotifyURL)
	}

	// 写入配置文件
	configFilePath := "./forward/ipv6_forward.config"
	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Printf("Error creating config file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	_, err = file.WriteString(configContent)
	if err != nil {
		return
	}

	// 推送配置文件到设备
	cmd = exec.Command(adb.AdbPath, "push", configFilePath, "/home/root/r200/")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error pushing config file: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已推送配置文件到设备")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 运行 forward 程序并显示输出
	fmt.Println("forward程序已运行")
	cmd = exec.Command(adb.AdbPath, "shell", "/usr/bin/forward")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating StdoutPipe for forward: %v\n", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating StderrPipe for forward: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting forward: %v\n", err)
		return
	}

	done := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- scanner.Err()
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- scanner.Err()
	}()

	go func() {
		if err := cmd.Wait(); err != nil {
			done <- err
		} else {
			done <- nil
		}
	}()

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("停止显示forward程序的交互信息")
		// 使用 adb 执行 Ctrl+C 停止 forward 程序
		stopCmd := exec.Command(adb.AdbPath, "shell", "killall forward")
		stopOutput, stopErr := stopCmd.CombinedOutput()
		if stopErr != nil {
			fmt.Printf("Error stopping forward: %v\n", stopErr)
		} else if len(stopOutput) == 0 {
			fmt.Println("forward程序已停止")
		} else {
			fmt.Printf("Output: %s\n", stopOutput) // 打印输出
		}
	case err := <-done:
		if err != nil {
			fmt.Printf("Error running forward: %v\n", err)
		} else {
			fmt.Println("forward程序已运行完毕")
		}
	}

	// 提示用户选择是否开启自启
	enableAutostart := getInput("是否开启forward自启? 输入1开启: ")
	if enableAutostart == "1" {
		cmd = exec.Command(adb.AdbPath, "shell", "systemctl start forward")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error enabling forward autostart: %v\n", err)
		} else if len(output) == 0 {
			fmt.Println("forward已设置为开机自启")
			// 等待用户输入回车
			fmt.Println("按回车键返回菜单...")
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadString('\n')

			// 清空屏幕并返回主菜单
			clearScreen()
			callback()
		} else {
			fmt.Printf("Output: %s\n", output) // 打印输出
		}
	} else {
		fmt.Println("未开启forward自启")
	}
}
