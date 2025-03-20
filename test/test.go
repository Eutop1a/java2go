package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// TEMPLATE_PATH 模板文件所在路径
const TEMPLATE_PATH = "services"

// isWindows 判断是否为 Windows 系统
func isWindows() bool {
	return os.PathSeparator == '\\'
}

// getHomePath 获取用户主目录
func getHomePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("无法获取用户主目录:", err)
		return ""
	}
	return home
}

// getTemplatePath 获取模板路径
func getTemplatePath() string {
	tmpDir := os.TempDir()
	if !isWindows() {
		tmpDir = filepath.Join(tmpDir, ".")
	}
	return tmpDir
}

func main() {
	// 准备输出文件
	outputFilePath := filepath.Join(getHomePath(), "hello.xml")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("创建文件时出错:", err)
		return
	}
	defer outputFile.Close()

	// 创建数据模型
	dataMap := map[string]interface{}{
		"title": "123",
	}

	// 加载模板文件
	templatePath := filepath.Join(TEMPLATE_PATH, "重庆邮电大学试题.ftl")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println("解析模板文件时出错:", err)
		return
	}

	// 执行模板并将结果写入文件
	err = tmpl.Execute(outputFile, dataMap)
	if err != nil {
		fmt.Println("执行模板时出错:", err)
		return
	}

	fmt.Println("文件生成成功:", outputFilePath)
}
