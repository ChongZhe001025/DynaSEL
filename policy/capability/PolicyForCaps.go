package capability

import (
	"fmt"
)

func CreatePolicyFromConfig(capabilities []map[string]interface{}, strPolicy string) (string, error) {
	strPolicy += "(block container_caps\n"
	strPolicy += "    (blockinherit)\n"

	for i, capsMap := range capabilities {
		strPolicy += fmt.Sprintf("    ; Process %d capabilities\n", i+1)
		strPolicy += "    (allow process self (capability ("

		filteredCaps := filterCapabilities(capsMap["bounding"].([]interface{}))
		for _, cap := range filteredCaps {
			cap, ok := cap.(string)
			if !ok {
				fmt.Printf("Expected string in filteredCaps, got %T", cap)
			}
			strPolicy += fmt.Sprintf(" " + cap[4:])
		}

		strPolicy += (" )))\n")
	}
	fmt.Println(strPolicy)
	return strPolicy, nil
}

var highRiskCaps = map[string]bool{
	"CAP_SYS_ADMIN":  true,
	"CAP_SYS_MODULE": true,
	"CAP_NET_ADMIN":  true,
}

var mediumRiskCaps = map[string]bool{
	"CAP_NET_RAW":      true,
	"CAP_DAC_OVERRIDE": true,
	"CAP_SYS_PTRACE":   true,
}

func filterCapabilities(caps []interface{}) []interface{} {
	var filtered []interface{}
	for _, cap := range caps {
		if !highRiskCaps[cap.(string)] && !mediumRiskCaps[cap.(string)] {
			filtered = append(filtered, cap)
		} else {
			// fmt.Printf("Filtered out capability: %s\n", cap)
		}
	}
	return filtered
}
