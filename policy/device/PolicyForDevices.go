package device

import (
	"fmt"
	"os"
	"strings"

	"github.com/opencontainers/selinux/go-selinux"
)

func CreatePolicyFromConfig(devices []map[string]interface{}, strPolicy string) (string, error) {
	for _, item := range devices {
		strPolicy += "    (deny container_t "

		// 檢查 item["path"] 的類型
		switch pathVal := item["path"].(type) {
		case []interface{}: // 如果是陣列，按預期處理
			highRiskFiltered := filterHighRiskDevices(pathVal)
			for _, paths := range highRiskFiltered {
				path, ok := paths.(string)
				if ok {
					strPolicy += getRealLabel(path)
				} else {
					fmt.Println("Warning: path is not a string")
				}
			}
		case string: // 如果是字串，直接處理
			highRiskFiltered := filterHighRiskDevices([]interface{}{pathVal})
			for _, paths := range highRiskFiltered {
				path, ok := paths.(string)
				if ok {
					strPolicy += getRealLabel(path)
				} else {
					fmt.Println("Warning: path is not a string")
				}
			}
		default: // 其他情況，報錯或跳過
			fmt.Println("Error: unexpected type for item['path']")
		}

		strPolicy += (" (chr_file (read write open ioctl getattr)))\n")
	}
	return strPolicy, nil
}

func getRealLabel(directory string) (context string) {
	context, err := selinux.FileLabel(directory)
	if err != nil && !os.IsNotExist(err) {
		return ""
	}
	if context != "" {
		parts := splitContext(context)
		if len(parts) > 2 {
			context = parts[2]
		}
		return context
	}
	return ""
}

var highRiskDevices = map[string]bool{
	"/dev/mem":         true, //"直接訪問物理內存，可能洩露敏感信息。"
	"/dev/kmem":        true, //"訪問內核內存，可能導致內核數據損壞或被利用。",
	"/dev/port":        true, //"訪問 I/O 端口，危害硬件安全。",
	"/dev/sd*":         true, //"SCSI 或 SATA 磁碟，可能導致數據損壞或洩露。",
	"/dev/nvme*":       true, //"NVMe 驅動的存儲設備，具有相同風險。",
	"/dev/loop*":       true, //"回環設備（虛擬磁碟）。",
	"/dev/md*":         true, //"RAID 設備。",
	"/dev/kvm":         true, //"虛擬機管理（KVM）設備。",
	"/dev/vhost-*":     true, //"虛擬化支持設備。",
	"/dev/vfio/*":      true, //"虛擬功能接口設備。",
	"/dev/virtio*":     true, //"VirtIO 驅動設備。",
	"/dev/net/tun":     true, //"虛擬網絡設備。",
	"/dev/ppp":         true, //"點對點協議支持。",
	"/dev/dri/*":       true, //"顯示驅動接口（Direct Rendering Interface）。",
	"/dev/fb*":         true, //"幀緩衝設備（Framebuffer）。",
	"/dev/nvidia*":     true, //"NVIDIA GPU 驅動。",
	"/dev/vga_arbiter": true, //"VGA 控制設備。",
	"/dev/tty*":        true, //"虛擬終端設備，可能導致敏感信息洩露。",
	"/dev/pts/*":       true, //"偽終端設備。",
	"/dev/console":     true, //"系統控制台。",
	"/dev/random":      true, //"隨機數生成設備（必要時允許）。",
	"/dev/urandom":     true, //"隨機數生成設備（必要時允許）。",
	"/dev/rtc*":        true, //"實時時鐘設備。",
	"/dev/cdrom":       true, //"光碟驅動器。",
	"/dev/dvd":         true, //"光碟驅動器。",
	"/dev/sr*":         true, //"SCSI 光碟設備。",
	"/dev/snd/*":       true, //"音頻設備。",
	"/dev/dsp":         true, //"數字信號處理設備。",
	"/dev/full":        true, //"返回 '磁碟已滿' 的虛擬設備，可能被惡意使用。",
	"/dev/null":        true, //"通常安全，但需審核用途。",
	"/dev/zero":        true, //"返回零值的設備。",
	"/dev/core":        true, //"核心轉儲設備。",
	"/dev/bus/*":       true, //"USB 和其他總線設備。",
}

// internal function
func splitContext(context string) []string {
	return strings.Split(context, ":")
}

func filterHighRiskDevices(caps []interface{}) []interface{} {
	var highRiskFiltered []interface{}
	for _, cap := range caps {
		if highRiskDevices[cap.(string)] {
			highRiskFiltered = append(highRiskFiltered, cap)
		} else {
			// fmt.Printf("Allowed capability: %s\n", cap)
		}
	}
	return highRiskFiltered
}
