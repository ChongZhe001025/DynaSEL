package automation

import (
	"DynaSEL-latest/monitor"
)

func AutoApplyPolicyToContainer(strArrConfigParentDirPath []string) {
	monitor.MonitorConfigJson(strArrConfigParentDirPath)
}
