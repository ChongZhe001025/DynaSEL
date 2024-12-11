package main

import (
	// "DynaSEL-latest/monitor"
	// "DynaSEL-latest/automation/test"
	"DynaSEL-latest/policy"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	strConfigDirPath := getConfigDirPath()
	strArrContainerID := getArrContainerID(strConfigDirPath)

	strArrConfigParentDirPath := []string{}

	for _, strContainerID := range strArrContainerID {
		// policy.CreateSElinuxPolicyFiles(strConfigDirPath, strContainerID)

		strCilFilePath := ("SysFiles/SELinuxPolicies/.cil/container_" + strContainerID + ".cil")

		policy.LoadPolicyToSELinux(strCilFilePath)

		// test.TestApplyPolicyToContainer(strContainerID, strPPFilePath)
		// automation.ApplyPolicyToContainer(strContainerID, strPPFilePath)

		strArrConfigParentDirPath = append(strArrConfigParentDirPath, (strConfigDirPath + "/" + strContainerID))
	}

	// monitor.MonitorConfigJson(strArrConfigParentDirPath)
}

// internal functions
func getArrContainerID(strConfigDirPath string) []string {
	var strArrContainerID []string

	entries, err := os.ReadDir(strConfigDirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "rootfs" {
			strArrContainerID = append(strArrContainerID, entry.Name())
		}
	}
	return strArrContainerID
}

func getConfigDirPath() string {
	findConfigCmd := exec.Command("bash", "-c", "find / -type f -name \"config.json\" 2>/dev/null | tail -n 1")

	var out bytes.Buffer
	findConfigCmd.Stdout = &out

	err := findConfigCmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}

	strSecondLastDir := strings.TrimSpace(out.String())
	strSecondLastDir = filepath.Dir(strSecondLastDir)

	parts := strings.Split(strSecondLastDir, string(os.PathSeparator))

	strSecondLastDir = filepath.Join(parts[:len(parts)-1]...)
	return "/" + strSecondLastDir
}
