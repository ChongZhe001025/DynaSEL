package mount

import (
	"fmt"
)

func CreatePolicyFromConfig(mounts []map[string]interface{}, strPolicy string) (string, error) {
	policySet := make(map[string]bool)

	if len(mounts) > 0 {
		for _, item := range mounts {
			destination, ok := item["destination"].(string)
			if !ok {
				return strPolicy, fmt.Errorf("invalid destination type in mounts, expected string")
			}

			options, ok := item["options"].([]interface{})
			if !ok {
				return strPolicy, fmt.Errorf("invalid options type in mounts, expected array of strings")
			}

			var optionStrings []string
			for _, opt := range options {
				if optStr, ok := opt.(string); ok {
					optionStrings = append(optionStrings, optStr)
				}
			}

			var policy string
			switch destination {
			case "/boot", "/root", "/etc/shadow", "/etc/passwd", "/etc/group":
				// 禁止存取敏感目錄
				policy = "    (deny container_t file_t (file (read write execute open)))\n"
			case "/lib", "/usr/bin":
				// 限制執行，僅允許只讀存取
				if contains(optionStrings, "nosuid") && contains(optionStrings, "nodev") {
					policy = "    (allow container_t lib_t (file (read)))\n"
				} else {
					policy = "    (deny container_t lib_t (file (read write execute)))\n"
				}
			case "/dev":
				// 僅允許必要設備
				if contains(optionStrings, "nosuid") {

					policy = "    (allow container_t null_device_t (chr_file (read write)))\n"
				} else {
					policy = "    (deny container_t null_device_t (chr_file (read write execute)))\n"
				}
			case "/var/lib/docker":
				// 禁止操作 Docker 資料目錄
				policy = "    (deny container_t docker_data_t (dir (read write execute)))\n"
			case "/sys/kernel/security":
				// 嚴格禁止存取安全模塊
				policy = "    (deny container_t security_t (dir (read write execute)))\n"
			case "/proc/kcore":
				// 禁止存取內存相關文件
				policy = "    (deny container_t proc_t (file (read write execute)))\n"
			case "/proc", "/sys":
				if contains(optionStrings, "nosuid") && contains(optionStrings, "noexec") && contains(optionStrings, "nodev") {
					policy = "    (allow container_t proc_t (dir (read)))\n"
				} else {
					policy = "    (deny container_t proc_t (dir (read write execute)))\n"
				}
			default:
				if source, ok := item["source"].(string); ok && source == "bind" {
					if destination == "/etc/hostname" || destination == "/etc/hosts" {
						policy = "    (allow container_t etc_t (file (read)))\n"
					} else {
						policy = "    (deny container_t file_t (file (read write execute)))\n"
					}
				}
			}
			if policy != "" && !policySet[policy] {
				policySet[policy] = true
				strPolicy += policy
			}
		}
	}
	return strPolicy, nil
}

// contains 檢查字串切片中是否包含指定字串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
