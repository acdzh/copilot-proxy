package main

import (
	"flag"
	"strings"

	"copilot-proxy/scan"
)

var GITHUB_COPILOT_API_HOST = "https://copilot-proxy.githubusercontent.com"

func main() {

	// 定义命令行参数
	isRecoverFlag := flag.Bool("r", false, "是否恢复原始 api 地址")
	newApiHostFlag := flag.String("u", GITHUB_COPILOT_API_HOST, "你的 api 地址")
	flag.Parse()

	isRecover := *isRecoverFlag
	newApiHost := strings.TrimRight(*newApiHostFlag, "/")

	// scan.ScanVSCode()
	scan.ScanJetBrains(newApiHost, isRecover)
	scan.ScanAndroidStudio(newApiHost, isRecover)
}
