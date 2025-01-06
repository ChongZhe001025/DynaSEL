package policy

import (
	"DynaSEL-latest/applicator"
	"DynaSEL-latest/parser"
	"DynaSEL-latest/policy/capability"
	"DynaSEL-latest/policy/device"
	"DynaSEL-latest/policy/mount"
	"DynaSEL-latest/policy/port"
	"fmt"
	"os"
	"os/exec"
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

func CreateSELinuxPolicyCil(strConfigDirPath string, strContainerID string) {
	strCilFilePath := ("SysFiles/SELinuxPolicies/containerPolicies/container_" + strContainerID + ".cil")

	filePolicyCil, err := os.Create(strCilFilePath)
	if err != nil {
		return
	}
	defer filePolicyCil.Close()

	strPolicy := fmt.Sprintf("(block container_%s\n", strContainerID)
	strPolicy += "    (blockinherit container)\n"

	parserResult := parser.GetParserResult()
	parserResult.Parse(strConfigDirPath, strContainerID)

	strPolicy = createPolicy(strPolicy, parserResult.MapStrInspectMounts, parserResult.MapStrConfigMounts, parserResult.MapStrInspectDevices, parserResult.MapStrConfigCaps, parserResult.MapStrInspectPorts)

	strPolicy += ")\n"

	_, err = filePolicyCil.WriteString(strPolicy)
	if err != nil {
		fmt.Println("fail")
	}

	loadPolicyToSELinux(strCilFilePath)

	applicator.ApplyPolicyToContainer(strContainerID)

}

func createPolicy(strPolicy string, inspect_mounts []map[string]interface{}, config_mounts []map[string]interface{}, devices []map[string]interface{}, capabilities []map[string]interface{}, ports []map[string]interface{}) string {

	// // Mounts inspect
	strPolicy, _ = mount.CreatePolicyFromInspectMounts(inspect_mounts, strPolicy)

	// // Mounts config
	strPolicy, _ = mount.CreatePolicyFromConfigMounts(config_mounts, strPolicy)

	// // Devices
	strPolicy, _ = device.CreatePolicyFromInspect(devices, strPolicy)

	// // Caps
	strPolicy, _ = capability.CreatePolicyFromConfig(capabilities, strPolicy)

	//Ports
	strPolicy, _ = port.CreatePolicyFromInspect(ports, strPolicy)
	return strPolicy
}

func loadPolicyToSELinux(strCilFilePath string) {
	cmdLoad := exec.Command("semodule", "-i", strCilFilePath)
	cmdLoad.Stdout = os.Stdout
	cmdLoad.Stderr = os.Stderr
	fmt.Println("Loading .cil file into SELinux...")
	if err := cmdLoad.Run(); err != nil {
		fmt.Printf("Failed to load .cil file into SELinux: %v\n", err)
		return
	}

	fmt.Println("SELinux policy loaded successfully!")
}
