package monitor

import (
	"context"
	pb "github.com/gofxq/host_monitor_agent_go/protos/protos"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"log"
	"runtime"
	"strings"
	"time"
)

var (
	Version           string
	expectDiskFsTypes = []string{
		"apfs", "ext4", "ext3", "ext2", "f2fs", "reiserfs", "jfs", "btrfs",
		"fuseblk", "zfs", "simfs", "ntfs", "fat32", "exfat", "xfs", "fuse.rclone",
	}
	excludeNetInterfaces = []string{
		"lo", "tun", "docker", "veth", "br-", "vmbr", "vnet", "kube",
	}
)

var (
	NetSpeedRcv, NetSpeedSnt, netInTransfer, netOutTransfer, lastUpdateNetStats uint64
	cachedBootTime                                                              time.Time
)

var hostID string

func init() {
	info, err := host.Info()
	if err != nil {
		log.Fatal(err)
	}
	hostID = info.HostID
}

func GetMonitor(c context.Context) *pb.MonitorTick {
	TrackNetworkSpeed()
	loadStat, err := load.Avg()
	if err != nil {
		loadStat = &load.AvgStat{}
	}

	ret := &pb.MonitorTick{
		HostId:             hostID,
		Timestamp:          time.Now().Unix(),
		NetSpeedSnt:        NetSpeedSnt,
		NetSpeedRcv:        NetSpeedRcv,
		NetTotalSnt:        netOutTransfer,
		NetTotalRcv:        netInTransfer,
		NetTotalPacketsSnt: 0,
		NetTotalPacketsRcv: 0,
		Load1:              loadStat.Load1,
		Load5:              loadStat.Load5,
		Load15:             loadStat.Load15,
	}

	cp, err := cpu.PercentWithContext(c, 0, false)
	if err != nil {
		println("cpu.Percent error:", err)
	} else {
		ret.CpuLoad = cp[0]
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		println("mem.VirtualMemory error:", err)
	} else {
		ret.MemUsed = vm.Total - vm.Available
		if runtime.GOOS != "windows" {
			ret.SwapUsed = vm.SwapTotal - vm.SwapFree
		}
	}
	if runtime.GOOS == "windows" {
		// gopsutil 在 Windows 下不能正确取 swap
		ms, err := mem.SwapMemory()
		if err != nil {
			println("mem.SwapMemory error:", err)
		} else {
			ret.SwapUsed = ms.Used
		}
	}

	return ret
}

// TrackNetworkSpeed NIC监控，统计流量与速度
func TrackNetworkSpeed() {
	var innerNetInTransfer, innerNetOutTransfer uint64
	nc, err := net.IOCounters(true)
	if err == nil {
		for _, v := range nc {
			if isListContainsStr(excludeNetInterfaces, v.Name) {
				continue
			}

			innerNetInTransfer += v.BytesRecv
			innerNetOutTransfer += v.BytesSent
		}
		now := uint64(time.Now().Unix())
		diff := now - lastUpdateNetStats
		if diff > 0 {
			NetSpeedRcv = (innerNetInTransfer - netInTransfer) / diff
			NetSpeedSnt = (innerNetOutTransfer - netOutTransfer) / diff
		}
		netInTransfer = innerNetInTransfer
		netOutTransfer = innerNetOutTransfer
		lastUpdateNetStats = now
	}
}

func isListContainsStr(list []string, str string) bool {
	for i := 0; i < len(list); i++ {
		if strings.Contains(str, list[i]) {
			return true
		}
	}
	return false
}
