package main

import (
	"context"
	"fmt"
	"github.com/gofxq/host_monitor_agent_go/monitor"
	pb "github.com/gofxq/host_monitor_agent_go/protos/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math"
	"time"
)

const (
	//target = "sh.gofxq.com:50051"
	//target = "sh.gofxq.com:50051"
	target = "192.168.6.3:50051"
)

func main() {
	// 设置一个连接到 gRPC 服务的连接
	conn, err := grpc.Dial(
		//"localhost:50051",
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建一个新的 MonitorService 客户端

	// 准备请求数据

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(time.Now())

	go ReportHostInfo(conn)

	go ReportMonitorInfo(conn)

	select {}

}

func ReportHostInfo(conn *grpc.ClientConn) {
	c := pb.NewMonitorServiceClient(conn)
	clientInfo, err := c.ReportHostInfoStream(context.Background())
	if err != nil {
		log.Fatalln("init ReportHostInfoStream failed")
	}
	t := time.Tick(time.Minute)
	for i := range t {
		hostInfo := monitor.GetHost(context.TODO())
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
	interval := time.Second
	c := pb.NewMonitorServiceClient(conn)
	clientInfo, err := c.ReportMonitorStream(context.Background())
	if err != nil {
		log.Fatalln("init ReportHostInfoStream failed")
	}
	t := time.Tick(interval)
	for range t {
		info := monitor.GetMonitor(context.TODO())
		log.Printf("[monitor]  host_id:%s\tcpu_load:%3.2f%%\tmem_used:%s\tswap_used:%s\tnet_speed_snt:%s\tnet_speed_rcv:%s\tnet_total_snt:%s\tnet_total_rcv:%s",
			info.HostId, info.CpuLoad, unitUint64(info.MemUsed), unitUint64(info.SwapUsed), unitUint64(info.NetSpeedSnt*8), unitUint64(info.NetSpeedRcv*8), unitUint64(info.NetTotalSnt), unitUint64(info.NetTotalRcv))
		// 调用 Report 方法
		err := clientInfo.Send(info)
		if err != nil {
			log.Printf("could not report: %v", err)
			time.Sleep(interval * 10)
			conn, err = grpc.Dial(
				//"localhost:50051",
				target,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			)
			clientInfo, err = pb.NewMonitorServiceClient(conn).ReportMonitorStream(context.Background())
		}
	}

}

func unitUint64(n uint64) string {
	return unitFloat(float64(n))
}

func unitFloat(number float64) string {
	units := []string{"", "K", "M", "G"}
	base := 1024.0

	if number < base {
		return fmt.Sprintf("%.0f", number)
	}

	exp := int(math.Log(number) / math.Log(base))
	scaledNumber := number / math.Pow(base, float64(exp))
	unit := units[exp]

	return fmt.Sprintf("%.1f%s", scaledNumber, unit)
}
