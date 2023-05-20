package main

import (
	"context"
	pb "github.com/gofxq/host_monitor_agent_go/protos/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	// 设置一个连接到 gRPC 服务的连接
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建一个新的 MonitorService 客户端
	c := pb.NewMonitorServiceClient(conn)

	// 准备请求数据
	hostInfo := &pb.HostInfo{
		DeviceId:     111,
		HostName:     "abc",
		Os:           "",
		OsVersion:    "",
		CpuName:      "",
		CpuMaxFreq:   "",
		CpuCoreCount: "",
		RamTotal:     "",
		RamFreq:      "",
		SwapTotal:    "",
		DiskInfos:    nil,
		GpuInfos:     nil,
	}

	log.Println(time.Now())
	for i := int64(0); i < 1e6; i++ {

		hostInfo.DeviceId = i
		// 调用 Report 方法
		response, err := c.ReportHostInfo(context.Background(), hostInfo)
		if err != nil {
			log.Printf("Report Response: %v, errCode:%d", response, response.ErrCode)
			log.Fatalf("could not report: %v", err)
		}
	}
	log.Println(time.Now())

}
