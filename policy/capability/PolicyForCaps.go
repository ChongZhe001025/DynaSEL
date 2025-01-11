package capability

import (
	"fmt"
	"strings"
)

func CreatePolicyFromConfig(capabilities []map[string]interface{}, strPolicy string) (string, error) {
	for _, capsMap := range capabilities {
		strPolicy += "    (deny container_t self (capability("

		highRiskFiltered := filterHighRiskCapabilities(capsMap["bounding"].([]interface{}))
		for _, cap := range highRiskFiltered {
			cap, ok := cap.(string)
			if !ok {
				fmt.Printf("Expected string in highRiskFiltered, got %T", cap)
			}
			strPolicy += " " + strings.ToLower(cap[4:])
		}

		strPolicy += (" )))\n")
	}
	return strPolicy, nil
}

var highRiskCaps = map[string]bool{
	"CAP_SYS_ADMIN":       true, // 全系統管理權限，最危險
	"CAP_SYS_MODULE":      true, // 加載/卸載內核模組
	"CAP_NET_ADMIN":       true, // 網路配置更改
	"CAP_SYS_RAWIO":       true, // 原始I/O訪問
	"CAP_SYS_TIME":        true, // 更改系統時間
	"CAP_SYS_BOOT":        true, // 重啟系統
	"CAP_SYS_RESOURCE":    true, // 設置資源限制
	"CAP_SYS_NICE":        true, // 調整任務優先級
	"CAP_IPC_LOCK":        true, // 鎖定共享記憶體段
	"CAP_AUDIT_CONTROL":   true, // 寫入審計記錄
	"CAP_DAC_OVERRIDE":    true, // 規避文件訪問權限檢查
	"CAP_CHOWN":           true, // 修改文件所有者
	"CAP_SETUID":          true, // 更改用戶ID
	"CAP_SETGID":          true, // 更改組ID
	"CAP_LINUX_IMMUTABLE": true, // 修改不可變文件屬性
	"CAP_FSETID":          true, // 設置文件系統的 UID 和 GID 位。
}

func filterHighRiskCapabilities(caps []interface{}) []interface{} {
	var highRiskFiltered []interface{}
	for _, cap := range caps {
		if highRiskCaps[cap.(string)] {
			highRiskFiltered = append(highRiskFiltered, cap)
		} else {
			// fmt.Printf("Allowed capability: %s\n", cap)
		}
	}
	return highRiskFiltered
}
