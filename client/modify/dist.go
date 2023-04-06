package modify

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"copilot-proxy/utils"
)

var GITHUB_COPILOT_API_HOST = "https://copilot-proxy.githubusercontent.com"
var LAST_MODIFIED_API_HOST_PATH_SUFFIX = ".replaced_api_host.txt"

// 用于获取上一次改动成了什么域名
func getLastModifiedApiHostWithJsPath(jsPath string) string {
	lastApiHostCachePath := jsPath + LAST_MODIFIED_API_HOST_PATH_SUFFIX
	exist, err := utils.IsExist(lastApiHostCachePath)
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

func modifyJs(jsPath, newApiHost string) {
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

func recoverJs(jsPath string) {
	jsContent, err := ioutil.ReadFile(jsPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	jsContentStr := string(jsContent)

	// recover
	lastModifiedApiHost := getLastModifiedApiHostWithJsPath(jsPath)
	if lastModifiedApiHost == "" {
		fmt.Println("Can not find js recover log, continue: " + jsPath)
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

func ModifyDistDir(distDirPath, newApiHost string) {
	exist, _ := utils.IsExist(distDirPath)
	if exist {
		distFilePaths, _ := ioutil.ReadDir(distDirPath) // 遍历 dist 文件夹
		for _, distFilePath := range distFilePaths {
			if strings.HasSuffix(distFilePath.Name(), ".js") { // 找到所有 js 文件, 修改
				modifyJs(filepath.Join(distDirPath, distFilePath.Name()), newApiHost)
			}
		}
	}
}

func RecoverDistDir(distDirPath string) {
	exist, _ := utils.IsExist(distDirPath)
	if exist {
		distFilePaths, _ := ioutil.ReadDir(distDirPath) // 遍历 dist 文件夹
		for _, distFilePath := range distFilePaths {
			if strings.HasSuffix(distFilePath.Name(), ".js") { // 找到所有 js 文件, 修改
				recoverJs(filepath.Join(distDirPath, distFilePath.Name()))
			}
		}
	}
}
