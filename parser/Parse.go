package parser

import (
	// "DynaSEL-latest/monitor"
	"DynaSEL-latest/parser/config"
	"encoding/json"
	"os"
	"path/filepath"
)

type ParserResult struct {
	MapStrConfigMounts  []map[string]interface{}
	MapStrConfigCaps    []map[string]interface{}
	MapStrConfigDevices []map[string]interface{}
}

func GetParserResult() *ParserResult {
	return &ParserResult{}
}

func (r *ParserResult) Parse(strConfigDirPath string, strContainerID string) {
	strConfigContent := getConfigFileContent(strConfigDirPath, strContainerID)

	mapStrConfigJson := parseToJson(string(strConfigContent))

	r.MapStrConfigMounts, _ = config.GetMountsFromConfig(mapStrConfigJson)
	r.MapStrConfigCaps, _ = config.GetCapsFromConfig(mapStrConfigJson)
	r.MapStrConfigDevices, _ = config.GetDevicesFromConfig(mapStrConfigJson)
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
