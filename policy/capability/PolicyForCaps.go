package capability

import (
	"fmt"
)

func CreatePolicyFromConfig(capabilities []map[string]interface{}, strPolicy string) (string, error) {

	for _, capsMap := range capabilities {
		strPolicy += "    (deny process self (capability ("

		highRiskFiltered := filterHighRiskCapabilities(capsMap["bounding"].([]interface{}))
		for _, cap := range highRiskFiltered {
			cap, ok := cap.(string)
			if !ok {
				fmt.Printf("Expected string in highRiskFiltered, got %T", cap)
			}
			strPolicy += " " + cap[4:]
		}

		strPolicy += (" )))\n")
	}
	return strPolicy, nil
}

var highRiskCaps = map[string]bool{
	"CAP_SYS_ADMIN":          true, // 全系統管理權限，最危險
	"CAP_SYS_MODULE":         true, // 加載/卸載內核模組
	"CAP_NET_ADMIN":          true, // 網路配置更改
	"CAP_SYS_RAWIO":          true, // 原始I/O訪問
	"CAP_SYS_TIME":           true, // 更改系統時間
	"CAP_SYS_BOOT":           true, // 重啟系統
	"CAP_SYS_RESOURCE":       true, // 設置資源限制
	"CAP_SYS_NICE":           true, // 調整任務優先級
	"CAP_IPC_LOCK":           true, // 鎖定共享記憶體段
	"CAP_MAC_ADMIN":          true, // 修改MAC配置
	"CAP_AUDIT_WRITE":        true, // 寫入審計記錄
	"CAP_AUDIT_READ":         true, // 讀取審計記錄
	"CAP_CHECKPOINT_RESTORE": true, // 遷移/還原進程
	"CAP_BPF":                true, // 訪問 eBPF 功能
	"CAP_SYSLOG":             true, // 操作系統日誌
	"CAP_DAC_OVERRIDE":       true, // 規避文件訪問權限檢查
	"CAP_CHOWN":              true, // 修改文件所有者
	"CAP_SETUID":             true, // 更改用戶ID
	"CAP_SETGID":             true, // 更改組ID
	"CAP_LINUX_IMMUTABLE":    true, // 修改不可變文件屬性
	"CAP_BLOCK_SUSPEND":      true, // 阻止系統掛起
}

// var mediumRiskCaps = map[string]bool{
// 	"CAP_NET_RAW":         true, // 原始套接字訪問
// 	"CAP_DAC_OVERRIDE":    true, // 規避文件訪問權限檢查
// 	"CAP_SYS_PTRACE":      true, // 訪問其他進程記憶體
// 	"CAP_AUDIT_CONTROL":   true, // 管理審計設置
// 	"CAP_CHOWN":           true, // 修改文件所有者
// 	"CAP_FOWNER":          true, // 避免文件權限檢查
// 	"CAP_SETUID":          true, // 更改用戶ID
// 	"CAP_SETGID":          true, // 更改組ID
// 	"CAP_NET_BIND_SERVICE": true, // 綁定低編號端口
// 	"CAP_LINUX_IMMUTABLE": true, // 修改不可變文件屬性
// 	"CAP_MKNOD":           true, // 創建特殊文件
// 	"CAP_WAKE_ALARM":      true, // 設置實時鬧鐘
// 	"CAP_BLOCK_SUSPEND":   true, // 阻止系統掛起
// 	"CAP_PERFMON":         true, // 性能監控
// }

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

// func filterCapabilities(caps []interface{}) []interface{} {
// 	var filtered []interface{}
// 	for _, cap := range caps {
// 		if !highRiskCaps[cap.(string)] && !mediumRiskCaps[cap.(string)] {
// 			filtered = append(filtered, cap)
// 		} else {
// 			// fmt.Printf("Filtered out capability: %s\n", cap)
// 		}
// 	}
// 	return filtered
// }
