package cmdb

import (
	"os"
	"runtime"
	"syscall"
)

var (
	kernel = syscall.NewLazyDLL("Kernel32.dll")
)

func (h *Host) GetHostname() string {
	name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return name
}
func (h *Host) GetOS() string {
	return runtime.GOOS
}
func (h *Host) GetGID() uint {
	//later will get judgements by process or file, or manually checked; and new will be written
	return 1
}
func (h *Host) GetExtras() map[string]interface{} {
	extras := make(map[string]interface{})
	extras["kernel"] = "3.15.1"
	return extras
}
func (h *Host) GetUptime() int64 {
	GetTickCount := kernel.NewProc("GetTickCount")
	r, _, _ := GetTickCount.Call()
	if r == 0 {
		return 0
	}
	//return s instead
	return int64(r / 1000)
}
