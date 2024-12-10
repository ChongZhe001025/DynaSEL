package device

import (
	"DynaSEL-latest/policy/auxiliary"
)

var perms = map[string]string{
	"device_rw": "getattr read write append ioctl lock open",
}

func CreatePolicyFromInspect(devices []map[string]interface{}, strPolicy string) (string, error) {
	var contexts []string
	for _, item := range devices {
		if pathOnHost, ok := item["path"].(string); ok && pathOnHost != "" {
			contexts = append(contexts, auxiliary.ListContexts(pathOnHost)...)
		}
	}

	for _, context := range auxiliary.SortedUnique(contexts) {
		strPolicy += ("    (allow process %s ( blk_file ( %s )))\n" + context + perms["device_rw"])
		strPolicy += ("    (allow process %s ( chr_file ( %s )))\n" + context + perms["device_rw"])
	}
	strPolicy += "\n)"
	return strPolicy, nil
}
