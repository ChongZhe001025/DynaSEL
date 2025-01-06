package mount

import (
	"DynaSEL-latest/policy/auxiliary"
	"fmt"
	"strings"
)

const (
	CONFIG_CONTAINER   = "/etc"
	HOME_CONTAINER     = "/home"
	LOG_CONTAINER      = "/var/log"
	TMP_CONTAINER      = "/tmp"
	TEMPLATE_PLAYBOOK  = "/usr/share/udica/ansible/deploy-module.yml"
	VARIABLE_FILE_NAME = "variables-deploy-module.yml"
)

var perms = map[string]string{
	"dir_rw":    "add_name create getattr ioctl lock open read remove_name rmdir search setattr write",
	"dir_ro":    "getattr ioctl lock open read search",
	"file_rw":   "append create getattr ioctl lock map open read rename setattr unlink write",
	"file_ro":   "getattr ioctl lock open read",
	"fifo_rw":   "getattr read write append ioctl lock open",
	"fifo_ro":   "getattr open read lock ioctl",
	"socket_rw": "append getattr open read write",
	"socket_ro": "getattr open read",
}
var templatesToLoad []string

func CreatePolicyFromInspectMounts(mounts []map[string]interface{}, strPolicy string) (string, error) {
	var contexts []string
	var contextsRW []string

	for _, item := range mounts {
		source, ok := item["source"].(string)
		if !ok || !strings.Contains(source, "/") {
			continue
		}
		rw, ok := item["RW"].(bool)
		if !ok {
			rw = false
		}

		switch source {
		case LOG_CONTAINER:
			if rw {
				strPolicy += "    (blockinherit log_rw_container)\n"
			} else {
				strPolicy += "    (blockinherit log_container)\n"
			}
			addTemplate("log_container")
			continue
		case HOME_CONTAINER:
			if rw {
				strPolicy += "    (blockinherit home_rw_container)\n"
			} else {
				strPolicy += "    (blockinherit home_container)\n"
			}
			addTemplate("home_container")
			continue
		case TMP_CONTAINER:
			if rw {
				strPolicy += "    (blockinherit tmp_rw_container)\n"
			} else {
				strPolicy += "    (blockinherit tmp_container)\n"
			}
			addTemplate("tmp_container")
			continue
		case CONFIG_CONTAINER:
			if rw {
				strPolicy += "    (blockinherit config_rw_container)\n"
			} else {
				strPolicy += "    (blockinherit config_container)\n"
			}
			addTemplate("config_container")
			continue
		default:
			if rw {
				contextsRW = append(contextsRW, auxiliary.ListContexts(source)...)
			} else {
				contexts = append(contexts, auxiliary.ListContexts(source)...)
			}
		}
	}

	strPolicy = appendContextsToPolicy(strPolicy, auxiliary.SortedUnique(contextsRW), perms, true)

	strPolicy = appendContextsToPolicy(strPolicy, auxiliary.SortedUnique(contexts), perms, false)

	return strPolicy, nil
}

func appendContextsToPolicy(policy string, contexts []string, perms map[string]string, isRW bool) string {
	for _, context := range contexts {
		dirPerm := perms["dir_ro"]
		filePerm := perms["file_ro"]
		fifoPerm := perms["fifo_ro"]
		socketPerm := perms["socket_ro"]
		if isRW {
			dirPerm = perms["dir_rw"]
			filePerm = perms["file_rw"]
			fifoPerm = perms["fifo_rw"]
			socketPerm = perms["socket_rw"]
		}

		policy += fmt.Sprintf("    (allow process %s ( dir ( %s )))\n", context, dirPerm)
		policy += fmt.Sprintf("    (allow process %s ( file ( %s )))\n", context, filePerm)
		policy += fmt.Sprintf("    (allow process %s ( fifo_file ( %s )))\n", context, fifoPerm)
		policy += fmt.Sprintf("    (allow process %s ( sock_file ( %s )))\n", context, socketPerm)
	}
	loadTemplates()
	return policy
}

func addTemplate(template string) {
	templatesToLoad = append(templatesToLoad, template)
}

func loadTemplates() {
	if len(templatesToLoad) == 0 {
		// fmt.Println("No templates to load.")
		return
	}
	// fmt.Println("Templates to load:")
	// for _, template := range templatesToLoad {
	// 	fmt.Printf(" - %s.cil\n", template)
	// }
}

//-------------------------------------------------------------------------------

func CreatePolicyFromConfigMounts(mounts []map[string]interface{}, strPolicy string) (string, error) {
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
			case "/proc", "/sys":
				if contains(optionStrings, "nosuid") && contains(optionStrings, "noexec") && contains(optionStrings, "nodev") {
					policy = "    (allow container_t proc_t (dir (read)))\n"
				} else {
					policy = "    (deny container_t proc_t (dir (read write execute)))\n"
				}
			case "/dev":
				if contains(optionStrings, "nosuid") {
					policy = "    (allow container_t null_device_t (chr_file (read write)))\n"
				} else {
					policy = "    (deny container_t null_device_t (chr_file (read write execute)))\n"
				}
			case "/dev/pts":
				policy = "    (allow container_t null_device_t (chr_file (read write)))\n"
			case "/dev/shm":
				if !contains(optionStrings, "noexec") {
					policy = "    (allow container_t tmpfs_t (dir (read write execute)))\n"
				} else {
					policy = "    (allow container_t tmpfs_t (dir (read write)))\n"
				}
			case "/dev/mqueue":
				if !contains(optionStrings, "noexec") {
					policy = "    (allow container_t mqueue_t (dir (read write execute)))\n"
				} else {
					policy = "    (allow container_t mqueue_t (dir (read write)))\n"
				}
			default:
				if source, ok := item["source"].(string); ok && source == "bind" {
					if destination == "/etc/hostname" || destination == "/etc/hosts" {
						policy = "    (allow container_t etc_t (file (read)))\n"
					} else {
						policy = "    (deny container_t etc_t (file (read write execute)))\n"
					}
				}
				// else {
				// 	policy = "    (deny container_t file (read write execute open))\n    (deny container_t dir (read write remove))\n    (deny container_t process (transition execute))\n"
				// }
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
