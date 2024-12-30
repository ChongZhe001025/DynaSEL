package main

import (
	"DynaSEL-latest/applicator"
	"DynaSEL-latest/monitor"
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

	for _, strContainerID := range strArrContainerID {
		policy.CreateSElinuxPolicyCil(strConfigDirPath, strContainerID)

		applicator.ApplyPolicyToContainer(strContainerID)
	}

	monitor.MonitorMobyDir(strConfigDirPath)
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
