package remote

import (
	"R200-tool/adb"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// 清空控制台屏幕

// 获取用户输入
func getInput1(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("输入不能为空，请重新输入")
			continue
		}
		return input
	}
}

func EnableNPS(callback func()) {
	// 提示用户输入nps服务端的IP、端口和密钥
	npsIP := getInput1("请输入nps服务端的IP: ")
	npsPort := getInput1("请输入nps服务端的端口: ")
	npsKey := getInput1("请输入nps服务端的密钥: ")

	// 读取npc.sh文件并替换相应的值
	npcFile := "./remote/nps/npc.sh"
	file, err := os.Open(npcFile)
	if err != nil {
		fmt.Printf("Error opening npc.sh: %v\n", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, "ExecStart=/home/root/r200/npc -server=ip:端口 -vkey=密钥 -type=tcp", fmt.Sprintf("ExecStart=/home/root/r200/npc -server=%s:%s -vkey=%s -type=tcp", npsIP, npsPort, npsKey), 1)
		content.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading npc.sh: %v\n", err)
		return
	}

	// 写回npc.sh文件
	err = os.WriteFile(npcFile, []byte(content.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing npc.sh: %v\n", err)
		return
	}

	// 推送nps文件夹到设备
	cmd := exec.Command(adb.AdbPath, "push", "./remote/nps", "/cache")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error pushing nps folder: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已推送nps文件夹到设备")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 设置npc权限
	cmd = exec.Command(adb.AdbPath, "shell", "chmod +x /cache/nps/npc")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error setting permissions for npc: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已设置npc权限")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 推送npc.sh文件到设备
	cmd = exec.Command(adb.AdbPath, "push", "./remote/nps/npc.sh", "/cache/nps/")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error pushing npc.sh: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已推送npc.sh文件到设备")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	// 运行npc.sh脚本并显示输出
	cmd = exec.Command(adb.AdbPath, "shell", "dos2unix /cache/nps/npc.sh")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error converting npc.sh to Unix format: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("已转换npc.sh为Unix格式")
	} else {
		fmt.Printf("Output: %s\n", output) // 打印输出
	}

	cmd = exec.Command(adb.AdbPath, "shell", "sh /cache/nps/npc.sh")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running npc.sh: %v\n", err)
	} else if len(output) == 0 {
		fmt.Println("npc.sh脚本已运行")
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
