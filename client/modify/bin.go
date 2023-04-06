package modify

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"copilot-proxy/utils"
)

var payload = `
#!/bin/bash

# 获取当前脚本所在目录
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# 进入脚本所在目录
cd "${DIR}"
cd "../dist"

# 启动 agent.js
node agent.js
`
var BACKUP_SUFFIX = ".backup"

func copy(destPath, srcPath string) (written int64, err error) {
	src, _ := os.Open(srcPath)
	defer src.Close()
	dst, _ := os.Create(destPath)
	defer dst.Close()

	return io.Copy(dst, src)
}

func modify(binPath string) {
	exist, _ := utils.IsExist(binPath)
	if exist {
		backupBinPath := binPath + BACKUP_SUFFIX
		exist, _ = utils.IsExist(backupBinPath)
		if !exist {
			copy(backupBinPath, binPath)
		}
		os.Remove(binPath)
	}
	ioutil.WriteFile(binPath, []byte(payload), 0777)
	fmt.Println("Modified: " + binPath)
}

func recover(binPath string) {
	backupBinPath := binPath + BACKUP_SUFFIX
	exist, _ := utils.IsExist(backupBinPath)
	if !exist {
		fmt.Println("Can not find bin recover backup, continue: " + binPath)
		return
	}

	os.Remove(binPath)
	os.Rename(backupBinPath, binPath)
	fmt.Println("Recover: " + binPath)
}

func ModifyBinDir(binDirPath string) {
	exist, _ := utils.IsExist(binDirPath)
	if !exist {
		return
	}
	for _, bin := range []string{"copilot-agent-linux", "copilot-agent-macos-arm64", "copilot-agent-macos"} {
		binPath := filepath.Join(binDirPath, bin)
		modify(binPath)
	}
}

func RecoverBinDir(binDirPath string) {
	exist, _ := utils.IsExist(binDirPath)
	if !exist {
		return
	}
	for _, bin := range []string{"copilot-agent-linux", "copilot-agent-macos-arm64", "copilot-agent-macos"} {
		binPath := filepath.Join(binDirPath, bin)
		recover(binPath)
	}
}
