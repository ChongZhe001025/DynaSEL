package inspect

import (
	"DynaSEL-latest/parser/auxiliary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func ParseToJson(data string) []map[string]interface{} {
	var parsedData map[string]interface{}
	json.Unmarshal([]byte(data), &parsedData)
	return []map[string]interface{}{parsedData}
}

// mounts
func GetMountsFromInspect(data []map[string]interface{}) ([]map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}

	if mounts, ok := data[0]["Mounts"].([]interface{}); ok {
		return auxiliary.ConvertToMountList(mounts), nil
	}

	return nil, fmt.Errorf("failed to parse mounts from data")
}

// devices
func GetDevicesFromInspect(data []map[string]interface{}) ([]map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}

	if hostConfig, ok := data[0]["HostConfig"].(map[string]interface{}); ok {
		if devices, ok := hostConfig["Devices"].([]interface{}); ok {
			return auxiliary.ConvertToDeviceList(devices), nil
		}
	}

	return nil, fmt.Errorf("failed to parse devices from data")
}

// ports
func GetPortsFromInspect(data []map[string]interface{}) ([]map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}

	var ports []map[string]interface{}
	if networkSettings, ok := data[0]["NetworkSettings"].(map[string]interface{}); ok {
		if portSettings, ok := networkSettings["Ports"].(map[string]interface{}); ok {
			for key := range portSettings {
				containerPort := strings.Split(key, "/")
				if len(containerPort) != 2 {
					continue
				}
				portNumber := containerPort[0]
				protocol := containerPort[1]

				newPort := map[string]interface{}{
					"portNumber": portNumber,
					"protocol":   protocol,
				}

				ports = append(ports, newPort)
			}
		}
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("failed to parse ports from data")
	}
	return ports, nil
}
