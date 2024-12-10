package policy

import (
	"DynaSEL-latest/parse"
	"DynaSEL-latest/policy/capability"
	"DynaSEL-latest/policy/device"
	"DynaSEL-latest/policy/mount"
	"DynaSEL-latest/policy/port"
	"fmt"
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

func CreateCilFile(strConfigDirPath string, strContainerID string) {
	filePolicyCil, err := os.Create(strContainerID + ".cil")
	if err != nil {
		return
	}
	defer filePolicyCil.Close()

	strPolicy := fmt.Sprintf("(block %s\n", strContainerID)
	strPolicy += "    (blockinherit container)\n"

	parserResult := parse.GetParserResult()
	parserResult.Parse(strConfigDirPath, strContainerID)

	strPolicy = createPolicy(strPolicy, parserResult.MapStrInspectMounts, parserResult.MapStrConfigMounts, parserResult.MapStrInspectDevices, parserResult.MapStrConfigCaps, parserResult.MapStrInspectPorts)

	strPolicy += ")\n"

	_, err = filePolicyCil.WriteString(strPolicy)
	if err != nil {
		fmt.Println("fail")
	}

	// loadPolicy(filePolicyCil)

}

func createPolicy(strPolicy string, inspect_mounts []map[string]interface{}, config_mounts []map[string]interface{}, devices []map[string]interface{}, capabilities []map[string]interface{}, ports []map[string]interface{}) string {

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

func loadPolicy(filePolicyCil *os.File) {

}
