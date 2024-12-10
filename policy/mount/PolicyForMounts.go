package mount

import (
	"DynaSEL-latest/policy/auxiliary"
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

func CreatePolicyFromInspectMounts(mounts []map[string]interface{}, strPolicy string) (string, error) {
	var contexts []string
	var contextsRW []string

	for _, item := range mounts {
		source := item["source"].(string)
		rw := item["RW"].(bool)

		if !strings.Contains(source, "/") {
			continue
		}

		switch source {
		case LOG_CONTAINER:
			if !rw {
				strPolicy += "    (blockinherit log_container)\n"
			} else {
				strPolicy += "    (blockinherit log_rw_container)\n"
			}
			continue

		case HOME_CONTAINER:
			if !rw {
				strPolicy += "    (blockinherit home_container)\n"
			} else {
				strPolicy += "    (blockinherit home_rw_container)\n"
			}
			continue

		case TMP_CONTAINER:
			if !rw {
				strPolicy += "    (blockinherit tmp_container)\n"
			} else {
				strPolicy += "    (blockinherit tmp_rw_container)\n"
			}
			continue

		case CONFIG_CONTAINER:
			if !rw {
				strPolicy += "    (blockinherit config_container)\n"
			} else {
				strPolicy += "    (blockinherit config_rw_container)\n"
			}
			continue
		}

		if rw {
			contextsRW = append(contextsRW, auxiliary.ListContexts(source)...)
		} else {
			contexts = append(contexts, auxiliary.ListContexts(source)...)
		}
	}

	for _, context := range auxiliary.SortedUnique(contextsRW) {
		strPolicy += ("    (allow process %s ( dir ( %s )))\n" + context + perms["dir_rw"])
		strPolicy += ("    (allow process %s ( file ( %s )))\n" + context + perms["file_rw"])
		strPolicy += ("    (allow process %s ( fifo_file ( %s )))\n" + context + perms["fifo_rw"])
		strPolicy += ("    (allow process %s ( sock_file ( %s )))\n" + context + perms["socket_rw"])
	}

	for _, context := range auxiliary.SortedUnique(contexts) {
		strPolicy += ("    (allow process %s ( dir ( %s )))\n" + context + perms["dir_ro"])
		strPolicy += ("    (allow process %s ( file ( %s )))\n" + context + perms["file_ro"])
		strPolicy += ("    (allow process %s ( fifo_file ( %s )))\n" + context + perms["fifo_ro"])
		strPolicy += ("    (allow process %s ( sock_file ( %s )))\n" + context + perms["socket_ro"])
	}
	return strPolicy, nil
}

func CreatePolicyFromConfigMounts(mounts []map[string]interface{}, strPolicy string) (string, error) {
	if len(mounts) > 0 {
		for _, item := range mounts {
			if destination, ok := item["destination"].(string); ok {
				if destination == "/proc" || destination == "/sys" {
					strPolicy += "    (allow container_t proc_t (dir (read)))\n"
				}
				if destination == "/dev" {
					strPolicy += "    (allow container_t null_device_t (chr_file (read write)))\n"
				}

				if destination == "/dev/pts" {
					strPolicy += "    (allow container_t null_device_t (chr_file (read write)))\n"
				}

				if destination == "/dev/shm" {
					strPolicy += "    (allow container_t tmpfs_t (dir (read write execute)))\n"

				}

				if destination == "/dev/mqueue" {
					strPolicy += "    (allow container_t mqueue_t (dir (read write execute)))\n"
				}

				if destination == "bind" {
					if destination == "/etc/hostname" || destination == "/etc/hosts" {
						strPolicy += "    (allow container_t etc_t (file (read)))\n"
					}
				}
			}
		}
	}
	return strPolicy, nil
}
