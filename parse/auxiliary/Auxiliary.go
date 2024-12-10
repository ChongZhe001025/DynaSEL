package auxiliary

func ConvertToDeviceList(devices []interface{}) []map[string]interface{} {
	deviceList := []map[string]interface{}{}
	for _, device := range devices {
		if deviceMap, ok := device.(map[string]interface{}); ok {
			deviceList = append(deviceList, deviceMap)
		}
	}
	return deviceList
}

func ConvertToMountList(mounts []interface{}) []map[string]interface{} {
	mountList := []map[string]interface{}{}
	for _, mount := range mounts {
		if mountMap, ok := mount.(map[string]interface{}); ok {
			mountList = append(mountList, mountMap)
		}
	}
	return mountList
}
