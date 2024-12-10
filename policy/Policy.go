package policy

import (
	"DynaSEL-latest/policy/capability"
	"DynaSEL-latest/policy/device"
	"DynaSEL-latest/policy/mount"
	"DynaSEL-latest/policy/port"
	"os"
)

const (
	CONFIG_CONTAINER   = "/etc"
	HOME_CONTAINER     = "/home"
	LOG_CONTAINER      = "/var/log"
	TMP_CONTAINER      = "/tmp"
	TEMPLATE_PLAYBOOK  = "/usr/share/udica/ansible/deploy-module.yml"
	VARIABLE_FILE_NAME = "variables-deploy-module.yml"
)

var TEMPLATES_STORE string

// var templatesToLoad []string

func CreatePolicy(strPolicy string, inspect_mounts []map[string]interface{}, config_mounts []map[string]interface{}, devices []map[string]interface{}, capabilities []map[string]interface{}, ports []map[string]interface{}) string {

	// Mounts inspect
	strPolicy, _ = mount.CreatePolicyFromInspectMounts(inspect_mounts, strPolicy)

	// Mounts config
	strPolicy, _ = mount.CreatePolicyFromConfigMounts(config_mounts, strPolicy)

	// Devices
	strPolicy, _ = device.CreatePolicyFromInspect(devices, strPolicy)

	// Caps
	strPolicy, _ = capability.CreatePolicyFromConfig(capabilities, strPolicy)

	//Ports
	strPolicy, _ = port.CreatePolicyFromInspect(ports, strPolicy)

	return strPolicy
}

func LoadPolicy(filePolicyCil *os.File) {

}
