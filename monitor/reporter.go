package monitor

import (
	"context"

	pb "github.com/gofxq/host_monitor_agent_go/protos/protos"
	"github.com/gofxq/host_monitor_agent_go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	ServerID       = "test_server"
	GrpcServerAddr = "host-cd.gofxq.com:8009"
	//GrpcServerAddr = "localhost:50051"
	//GrpcServerAddr = "192.168.6.3:50051"
)

func ReportHostInfo(conn *grpc.ClientConn) {
	c := pb.NewMonitorServiceClient(conn)
	clientInfo, err := c.ReportHostInfoStream(context.Background())
	if err != nil {
		log.Fatalln("init ReportHostInfoStream failed")
	}
	t := time.Tick(time.Minute)
	for i := range t {
		hostInfo := GetHost(context.TODO())
		log.Println(i, time.Now())
		// 调用 Report 方法
		err := clientInfo.Send(hostInfo)
		if err != nil {
			log.Printf("could not report: %v", err)
			time.Sleep(time.Minute * 10)
		}
	}
}

func ReportMonitorInfo(conn *grpc.ClientConn) {
	interval := time.Second //*3
	c := pb.NewMonitorServiceClient(conn)
	clientInfo, err := c.ReportMonitorStream(context.Background())
	if err != nil {
		log.Fatalln("init ReportHostInfoStream failed")
	}
	for range time.Tick(interval) {
		info := GetMonitor(context.TODO())
		log.Printf("[monitor]  id:%s\tcpu:%3.2f%%\tmem:%s\tswap:%s\tnet_s:%s\tnet_r:%s\tnet_st:%s\tnet_rt:%s",
			info.HostId, info.CpuLoad, utils.UnitUinUt64(info.MemUsed), utils.UnitUinUt64(info.SwapUsed), utils.UnitUinUt64(info.NetSpeedSnt*8), utils.UnitUinUt64(info.NetSpeedRcv*8), utils.UnitUinUt64(info.NetTotalSnt), utils.UnitUinUt64(info.NetTotalRcv))
		// 调用 Report 方法
		err := clientInfo.Send(info)
		if err != nil {
			log.Printf("could not report: %v", err)
			time.Sleep(interval * 3)
			conn, err = grpc.Dial(
				//"localhost:50051",
				GrpcServerAddr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			)
			clientInfo, err = pb.NewMonitorServiceClient(conn).ReportMonitorStream(context.Background())
		}
	}

}
