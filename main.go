package main

import (
	"DynaSEL-latest/parse_oci_json/config"
	"DynaSEL-latest/parse_oci_json/inspect"
	"DynaSEL-latest/policy"
	"bytes"
	"encoding/json"
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

		strConfigContent := getConfigFileContent(strConfigDirPath, strContainerID)

		strInspectContent := getInspectContent(strContainerID)

		mapStrInspectJson := parseToJson(string(strInspectContent))
		mapStrConfigJson := parseToJson(string(strConfigContent))

		mapStrConfigMounts, _ := config.GetMountsFromConfig(mapStrConfigJson)
		// fmt.Printf("mapStrConfigMounts: ")
		// fmt.Println(mapStrConfigMounts)

		mapStrConfigCaps, _ := config.GetCapsFromConfig(mapStrConfigJson)
		// fmt.Printf("mapStrConfigCaps: ")
		// fmt.Println(mapStrConfigCaps)

		mapStrInspectMounts, _ := inspect.GetMountsFromInspect(mapStrInspectJson)
		// fmt.Printf("mapStrInspectMounts: ")
		// fmt.Println(mapStrInspectMounts)

		mapStrInspectDevices, _ := inspect.GetDevicesFromInspect(mapStrInspectJson)
		// fmt.Printf("mapStrInspectDevices: ")
		// fmt.Println(mapStrInspectDevices)

		mapStrInspectPorts, _ := inspect.GetPortsFromInspect(mapStrInspectJson)
		// fmt.Printf("mapStrInspectPorts: ")
		// fmt.Println(mapStrInspectPorts)

		filePolicyCil, err := os.Create(strContainerID + ".cil")
		if err != nil {
			return
		}
		defer filePolicyCil.Close()

		strPolicy := fmt.Sprintf("(block %s\n", strContainerID)
		strPolicy += "    (blockinherit container)\n"

		strPolicy = policy.CreatePolicy(strPolicy, mapStrInspectMounts, mapStrConfigMounts, mapStrInspectDevices, mapStrConfigCaps, mapStrInspectPorts)

		strPolicy += ")\n"

		_, err = filePolicyCil.WriteString(strPolicy)
		if err != nil {
			fmt.Println("fail")
		}

		policy.LoadPolicy(filePolicyCil)

	}
}

// internal functions
func parseToJson(data string) []map[string]interface{} {
	var parsedData map[string]interface{}
	json.Unmarshal([]byte(data), &parsedData)
	return []map[string]interface{}{parsedData}
}

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

func getConfigFileContent(strConfigDirPath string, strContainerID string) string {
	strConfigFilePath := filepath.Join("/"+strConfigDirPath+"/"+strContainerID, "config.json")
	byteConfigFileContent, _ := os.ReadFile(strConfigFilePath)
	strConfigFileContent := string(byteConfigFileContent)

	return strConfigFileContent
}

func getInspectContent(strContainerID string) string {
	InspectContainerCmd := exec.Command("docker", "inspect", strContainerID)
	strInspectOutput, err := InspectContainerCmd.Output()
	if err != nil {
		fmt.Println("can't get containerID")
	}

	parseOutput := strings.TrimSpace(string(strInspectOutput))
	parseOutput = parseOutput[1 : len(parseOutput)-1]

	return parseOutput
}
