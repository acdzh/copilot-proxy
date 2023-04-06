package scan

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"copilot-proxy/modify"
	"copilot-proxy/utils"
)

func ScanAndroidStudio(newApiHost string, isRecover bool) {
	var homeDir, _ = os.UserHomeDir()
	googleDir := filepath.Join(homeDir, "Library", "Application Support", "Google")
	exist, _ := utils.IsExist(googleDir)
	if !exist {
		return
	}
	programPaths, _ := ioutil.ReadDir(googleDir)
	for _, programPath := range programPaths {
		if strings.HasPrefix(programPath.Name(), "AndroidStudio") {
			copilotAgentPath := filepath.Join(googleDir, programPath.Name(), "plugins", "github-copilot-intellij", "copilot-agent")
			exist, _ := utils.IsExist(copilotAgentPath)
			if exist {
				if isRecover {
					modify.RecoverBinDir(filepath.Join(copilotAgentPath, "bin"))
					modify.RecoverDistDir(filepath.Join(copilotAgentPath, "dist"))
				} else {
					modify.ModifyBinDir(filepath.Join(copilotAgentPath, "bin"))
					modify.ModifyDistDir(filepath.Join(copilotAgentPath, "dist"), newApiHost)
				}
			}
		}
	}
}
