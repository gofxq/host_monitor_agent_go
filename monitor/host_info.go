package monitor

import (
	"context"
	pb "github.com/gofxq/host_monitor_agent_go/protos/protos"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// GetHost 获取主机硬件信息
func GetHost(c context.Context) *pb.TickHostInfo {

	hostInfo := &pb.TickHostInfo{
		ServerId:  ServerID,
		Header:    GetHostHeader(c),
		CpuInfos:  GetHostCpus(c),
		GpuInfos:  nil,
		MemInfo:   GetHostMem(c),
		DiskInfos: GetHostDisks(c),
	}
	return hostInfo
}

func GetHostHeader(c context.Context) *pb.Header {
	info, err := host.Info()
	if err != nil {
		return nil
	}
	return &pb.Header{
		DeviceId:             1,
		Hostname:             info.Hostname,
		Uptime:               info.Uptime,
		BootTime:             info.BootTime,
		Procs:                info.Procs,
		Os:                   info.OS,
		Platform:             info.Platform,
		PlatformFamily:       info.PlatformFamily,
		PlatformVersion:      info.PlatformVersion,
		KernelVersion:        info.KernelVersion,
		KernelArch:           info.KernelArch,
		VirtualizationSystem: info.VirtualizationSystem,
		VirtualizationRole:   info.VirtualizationRole,
		HostId:               info.HostID,
	}
}

func GetHostCpus(c context.Context) []*pb.CpuInfo {
	info, err := cpu.Info()
	if err != nil {
		return nil
	}
	res := make([]*pb.CpuInfo, 0, len(info))
	for i := range info {
		res = append(res, &pb.CpuInfo{
			Cpu:        info[i].CPU,
			VendorId:   info[i].VendorID,
			Family:     info[i].Family,
			Model:      info[i].Model,
			Stepping:   info[i].Stepping,
			PhysicalId: info[i].PhysicalID,
			CoreId:     info[i].CoreID,
			Cores:      info[i].Cores,
			ModelName:  info[i].ModelName,
			Mhz:        info[i].Mhz,
			CacheSize:  info[i].CacheSize,
			Flags:      info[i].Flags,
			Microcode:  info[i].Microcode,
		})
	}
	return res
}

func GetHostDisks(c context.Context) []*pb.DiskInfo {
	parts, err := disk.PartitionsWithContext(c, true)
	if err != nil {
		return nil
	}
	res := make([]*pb.DiskInfo, 0, len(parts))
	for _, p := range parts {
		info, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		res = append(res, &pb.DiskInfo{
			Path:              info.Path,
			Fstype:            info.Fstype,
			Total:             info.Total,
			Free:              info.Free,
			Used:              info.Used,
			UsedPercent:       info.UsedPercent,
			InodesTotal:       info.InodesTotal,
			InodesUsed:        info.InodesUsed,
			InodesFree:        info.InodesFree,
			InodesUsedPercent: info.InodesUsedPercent,
		})
	}
	return res
}
func GetHostMem(c context.Context) *pb.MemInfo {
	info, err := mem.VirtualMemoryWithContext(c)
	if err != nil {
		return nil
	}
	return &pb.MemInfo{
		Total: info.Total,
		Used:  info.Used,
		Free:  info.Free,
	}
}
func GetHostGPUs(c context.Context) []*pb.GpuInfo {
	return nil
}
