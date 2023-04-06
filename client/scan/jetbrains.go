package scan

import (
	"os"
	"path/filepath"
	"strings"

	"copilot-proxy/modify"
	"copilot-proxy/utils"
)

func ScanJetBrains(newApiHost string, isRecover bool) {
	var homeDir, _ = os.UserHomeDir()
	jetBrainsDir := filepath.Join(homeDir, "Library", "Application Support", "JetBrains")
	exist, _ := utils.IsExist(jetBrainsDir)
	if !exist {
		return
	}

	filepath.Walk(jetBrainsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasSuffix(path, "node_modules") ||
				strings.HasSuffix(path, "JavascriptLanguage") ||
				strings.HasSuffix(path, "javascript") ||
				strings.HasSuffix(path, "javascript-impl") ||
				strings.HasSuffix(path, "lib") ||
				strings.HasSuffix(path, "jbr") ||
				strings.HasSuffix(path, "maven") ||
				strings.HasSuffix(path, "gradle") ||
				strings.HasSuffix(path, "textmate") {
				// 防止进入黑洞
				return filepath.SkipDir
			}
			if strings.HasSuffix(path, filepath.Join("github-copilot-intellij", "copilot-agent")) {
				if isRecover {
					modify.RecoverBinDir(filepath.Join(path, "bin"))
					modify.RecoverDistDir(filepath.Join(path, "dist"))
				} else {
					modify.ModifyBinDir(filepath.Join(path, "bin"))
					modify.ModifyDistDir(filepath.Join(path, "dist"), newApiHost)
				}
			}
		}

		return nil
	})
}
