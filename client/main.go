package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var homeDir, _ = os.UserHomeDir()
var GITHUB_COPILOT_API_HOST = "https://copilot-proxy.githubusercontent.com"

var LAST_MODIFIED_API_HOST_PATH_SUFFIX = ".replaced_api_host.txt"

func isExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// 用于获取上一次改动成了什么域名
func getLastModifiedApiHostWithJsPath(jsPath string) string {
	lastApiHostCachePath := jsPath + LAST_MODIFIED_API_HOST_PATH_SUFFIX
	exist, err := isExist(lastApiHostCachePath)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	if !exist {
		return ""
	}
	content, err := ioutil.ReadFile(lastApiHostCachePath)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	contentStr := string(content)
	return contentStr
}

// 用于设置本次改动成了什么域名, 便于恢复
func setLastModifiedApiHostWithJsPath(jsPath, newApiHost string) {
	lastApiHostCachePath := jsPath + LAST_MODIFIED_API_HOST_PATH_SUFFIX
	err := ioutil.WriteFile(lastApiHostCachePath, []byte(newApiHost), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// 删除修改记录
func rmLastModifiedApiHostWithJsPath(jsPath string) {
	lastApiHostCachePath := jsPath + LAST_MODIFIED_API_HOST_PATH_SUFFIX
	err := os.Remove(lastApiHostCachePath)
	if err != nil {
		log.Fatal(err)
	}
}

func modifyJsWithPath(jsPath, newApiHost string) {
	jsContent, err := ioutil.ReadFile(jsPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	jsContentStr := string(jsContent)

	// modify
	lastModifiedApiHost := getLastModifiedApiHostWithJsPath(jsPath)
	var modifiedContentStr string = jsContentStr

	// 预防漏网之鱼
	modifiedContentStr = strings.ReplaceAll(modifiedContentStr, GITHUB_COPILOT_API_HOST, newApiHost)
	if lastModifiedApiHost != "" {
		modifiedContentStr = strings.ReplaceAll(modifiedContentStr, lastModifiedApiHost, newApiHost)
	}
	setLastModifiedApiHostWithJsPath(jsPath, newApiHost)

	// write
	if modifiedContentStr == jsContentStr {
		fmt.Println("No need to modify js file: " + jsPath)
	} else {
		err = ioutil.WriteFile(jsPath, []byte(modifiedContentStr), 0644)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("Modified js file: " + jsPath + "\n\t" + func() string {
			if lastModifiedApiHost == "" {
				return GITHUB_COPILOT_API_HOST
			} else {
				return lastModifiedApiHost
			}
		}() + " --> " + newApiHost)
	}
}

func recoverJsWithPath(jsPath string) {
	jsContent, err := ioutil.ReadFile(jsPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	jsContentStr := string(jsContent)

	// recover
	lastModifiedApiHost := getLastModifiedApiHostWithJsPath(jsPath)
	if lastModifiedApiHost == "" {
		fmt.Println("Can not find recover log, continue.")
		return
	}
	modifiedContentStr := strings.ReplaceAll(jsContentStr, lastModifiedApiHost, GITHUB_COPILOT_API_HOST)
	rmLastModifiedApiHostWithJsPath(jsPath)

	// write
	if modifiedContentStr == jsContentStr {
		fmt.Println("No need to recover js file: " + jsPath)
	} else {
		err = ioutil.WriteFile(jsPath, []byte(modifiedContentStr), 0644)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("Recovered js file: " + jsPath + "\n\t" + lastModifiedApiHost + " --> " + GITHUB_COPILOT_API_HOST)
	}
}

func getVisualStudioCodeExtJsPathList() []string {
	var files []string

	vscodeExtsDir := filepath.Join(homeDir, ".vscode", "extensions")
	exist, _ := isExist(vscodeExtsDir) // vscode 扩展文件夹是否存在
	if !exist {
		return files
	}
	extDirs, _ := ioutil.ReadDir(vscodeExtsDir) // 遍历 vscode 扩展文件夹
	for _, extDir := range extDirs {
		if strings.HasPrefix(extDir.Name(), "github.copilot") { // 找到所有 github.copilot 开头的文件夹
			distDir := filepath.Join(vscodeExtsDir, extDir.Name(), "dist")
			exist, _ := isExist(distDir) // dist 文件夹是否存在
			if exist {
				distFilePaths, _ := ioutil.ReadDir(distDir) // 遍历 dist 文件夹
				for _, distFilePath := range distFilePaths {
					if strings.HasSuffix(distFilePath.Name(), ".js") { // 找到所有 js 文件
						files = append(files, filepath.Join(distDir, distFilePath.Name()))
					}
				}
			}
		}
	}

	return files
}

func getJetBrainsPluginJsPathList() []string {
	var files []string

	jetBrainsDir := filepath.Join(homeDir, "Library", "Application Support", "JetBrains")
	exist, _ := isExist(jetBrainsDir)
	if !exist {
		return files
	}

	ideDirs, _ := ioutil.ReadDir(jetBrainsDir)
	for _, ideDir := range ideDirs {
		distDir := filepath.Join(jetBrainsDir, ideDir.Name(), "plugins", "github-copilot-intellij", "copilot-agent", "dist")
		exist, _ := isExist(distDir)
		if exist {
			distFilePaths, _ := ioutil.ReadDir(distDir) // 遍历 dist 文件夹
			for _, distFilePath := range distFilePaths {
				if strings.HasSuffix(distFilePath.Name(), ".js") { // 找到所有 js 文件, 修改
					files = append(files, filepath.Join(distDir, distFilePath.Name()))
				}
			}
		}
	}

	return files
}

func main() {

	// 定义命令行参数
	isRecoverFlag := flag.Bool("r", false, "是否恢复原始 api 地址")
	newApiHostFlag := flag.String("u", GITHUB_COPILOT_API_HOST, "你的 api 地址")
	flag.Parse()

	isRecover := *isRecoverFlag
	newApiHost := strings.TrimRight(*newApiHostFlag, "/")

	files := append(getVisualStudioCodeExtJsPathList(), getJetBrainsPluginJsPathList()...)
	for _, file := range files {
		if isRecover {
			recoverJsWithPath(file)
		} else {
			modifyJsWithPath(file, newApiHost)
		}
	}
}
