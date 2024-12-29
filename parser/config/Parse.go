package config

import (
	"DynaSEL-latest/parser/auxiliary"
	"encoding/json"
	"errors"
	"fmt"
)

func ParseToJson(data string) []map[string]interface{} {
	var parsedData map[string]interface{}
	json.Unmarshal([]byte(data), &parsedData)
	return []map[string]interface{}{parsedData}
}

// mounts
func GetMountsFromConfig(data []map[string]interface{}) ([]map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}

	if mounts, ok := data[0]["mounts"].([]interface{}); ok {
		return auxiliary.ConvertToMountList(mounts), nil
	}

	return nil, fmt.Errorf("failed to parse mounts from data")
}

// caps
func GetCapsFromConfig(data []map[string]interface{}) ([]map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}
	if process, ok := data[0]["process"].(map[string]interface{}); ok {
		if caps, ok := process["capabilities"].(map[string]interface{}); ok {
			return []map[string]interface{}{caps}, nil
		}
	}

	return nil, fmt.Errorf("failed to parse bounding from data")
}
