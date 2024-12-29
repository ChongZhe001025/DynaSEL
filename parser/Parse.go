package parser

import (
	// "DynaSEL-latest/monitor"
	"DynaSEL-latest/parser/config"
	"DynaSEL-latest/parser/inspect"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ParserResult struct {
	MapStrConfigMounts   []map[string]interface{}
	MapStrConfigCaps     []map[string]interface{}
	MapStrInspectMounts  []map[string]interface{}
	MapStrInspectDevices []map[string]interface{}
	MapStrInspectPorts   []map[string]interface{}
}

func GetParserResult() *ParserResult {
	return &ParserResult{}
}

func (r *ParserResult) Parse(strConfigDirPath string, strContainerID string) {
	strConfigContent := getConfigFileContent(strConfigDirPath, strContainerID)
	strInspectContent := getInspectContent(strContainerID)

	mapStrInspectJson := parseToJson(string(strInspectContent))
	mapStrConfigJson := parseToJson(string(strConfigContent))

	r.MapStrConfigMounts, _ = config.GetMountsFromConfig(mapStrConfigJson)
	r.MapStrConfigCaps, _ = config.GetCapsFromConfig(mapStrConfigJson)
	r.MapStrInspectMounts, _ = inspect.GetMountsFromInspect(mapStrInspectJson)
	r.MapStrInspectDevices, _ = inspect.GetDevicesFromInspect(mapStrInspectJson)
	r.MapStrInspectPorts, _ = inspect.GetPortsFromInspect(mapStrInspectJson)
}

// internal functions
func parseToJson(data string) []map[string]interface{} {
	var parsedData map[string]interface{}
	json.Unmarshal([]byte(data), &parsedData)
	return []map[string]interface{}{parsedData}
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
