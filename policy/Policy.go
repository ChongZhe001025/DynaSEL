package policy

import (
	"DynaSEL-latest/apply"
	"DynaSEL-latest/parse"
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

func CreateSElinuxPolicyFiles(strConfigDirPath string, strContainerID string) {
	strTeFilePath := ("SysFiles/SELinuxPolicies/.te/container_" + strContainerID + ".te")
	// strModFilePath := ("SysFiles/SELinuxPolicies/.mod/container_" + strContainerID + ".mod")
	strPPFilePath := ("SysFiles/SELinuxPolicies/.pp/container_" + strContainerID + ".pp")

	filePolicyCil, err := os.Create(strTeFilePath)
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

	// LoadPolicyToSELinux(strTeFilePath, strModFilePath, strPPFilePath)

	apply.ApplyPolicyToContainer(strContainerID, strPPFilePath)

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

func LoadPolicyToSELinux(strCilFilePath string) {
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

// func LoadPolicyToSELinux(strTeFilePath string, strModFilePath string, strPPFilePath string) {
// 	// Step 1: 編譯 .te 文件成 .mod 文件
// 	cmdCompile := exec.Command("checkmodule", "-M", "-m", "-o", strModFilePath, strTeFilePath)
// 	cmdCompile.Stdout = os.Stdout
// 	cmdCompile.Stderr = os.Stderr
// 	fmt.Println("Compiling .te file...")
// 	if err := cmdCompile.Run(); err != nil {
// 		fmt.Printf("Failed to compile .te file: %v\n", err)
// 		return
// 	}

// 	// Step 2: 生成 .pp 文件
// 	cmdPackage := exec.Command("semodule_package", "-o", strPPFilePath, "-m", strModFilePath)
// 	cmdPackage.Stdout = os.Stdout
// 	cmdPackage.Stderr = os.Stderr
// 	fmt.Println("Creating .pp package...")
// 	if err := cmdPackage.Run(); err != nil {
// 		fmt.Printf("Failed to create .pp package: %v\n", err)
// 		return
// 	}

// 	// Step 3: 載入 .pp 文件到 SELinux
// 	cmdLoad := exec.Command("semodule", "-i", strPPFilePath)
// 	cmdLoad.Stdout = os.Stdout
// 	cmdLoad.Stderr = os.Stderr
// 	fmt.Println("Loading .pp file into SELinux...")
// 	if err := cmdLoad.Run(); err != nil {
// 		fmt.Printf("Failed to load .pp file into SELinux: %v\n", err)
// 		return
// 	}

// 	fmt.Println("SELinux policy loaded successfully!")
// }
