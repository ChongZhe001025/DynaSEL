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
				policy = "    (deny container_t file_t (file (read write execute open)))\n"
			case "/lib", "/usr/bin":
				if contains(optionStrings, "nosuid") && contains(optionStrings, "nodev") {
					policy = "    (allow container_t lib_t (file (read)))\n"
				} else {
					policy = "    (deny container_t lib_t (file (read write execute)))\n"
				}
			case "/dev":
				if contains(optionStrings, "nosuid") {

					policy = "    (allow container_t null_device_t (chr_file (read write)))\n"
				} else {
					policy = "    (deny container_t null_device_t (chr_file (read write execute)))\n"
				}
			case "/var/lib/docker":
				policy = "    (deny container_t docker_data_t (dir (read write execute)))\n"
			case "/sys/kernel/security":
				policy = "    (deny container_t security_t (dir (read write execute)))\n"
			case "/proc/kcore":
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
