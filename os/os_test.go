package os

import (
	"os"
	"runtime"
	"testing"

	"github.com/shirou/gopsutil/disk"

	"github.com/shirou/gopsutil/host"

	"github.com/shirou/gopsutil/mem"
)

func TestOS(t *testing.T) {
	t.Log(runtime.GOARCH)
	t.Log(runtime.GOOS)

	t.Log(runtime.NumCPU())
	t.Log(runtime.NumGoroutine())
	t.Log(os.Environ())
	m, _ := mem.VirtualMemory()
	t.Log(m.Total)
	hStat, _ := host.Info()
	t.Log(hStat.Platform, hStat.PlatformFamily, hStat.PlatformVersion)
	t.Log()
	pStats, _ := disk.Partitions(true)
	t.Log(pStats)

}
