package scan

import "fmt"

func ScanVSCode() {
	// vscodeExtsDir := filepath.Join(homeDir, ".vscode", "extensions")
	// exist, _ := isExist(vscodeExtsDir) // vscode 扩展文件夹是否存在
	// if !exist {
	// 	return files
	// }
	// extDirs, _ := ioutil.ReadDir(vscodeExtsDir) // 遍历 vscode 扩展文件夹
	// for _, extDir := range extDirs {
	// 	if strings.HasPrefix(extDir.Name(), "github.copilot") { // 找到所有 github.copilot 开头的文件夹
	// 		distDir := filepath.Join(vscodeExtsDir, extDir.Name(), "dist")
	// 		exist, _ := isExist(distDir) // dist 文件夹是否存在
	// 		if exist {
	// 			distFilePaths, _ := ioutil.ReadDir(distDir) // 遍历 dist 文件夹
	// 			for _, distFilePath := range distFilePaths {
	// 				if strings.HasSuffix(distFilePath.Name(), ".js") { // 找到所有 js 文件
	// 					files = append(files, filepath.Join(distDir, distFilePath.Name()))
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	fmt.Println("vscode 请直接修改配置如下")
	fmt.Println("\"github.copilot.advanced\":{ \"debug.testOverrideProxyUrl\":\"你的 host\", \"debug.overrideProxyUrl\": \"你的 host\" }")
	fmt.Println("")
}
