package main

import (
	"context"
	"github.com/gofxq/host_monitor_agent_go/monitor"
	pb "github.com/gofxq/host_monitor_agent_go/protos/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	// 设置一个连接到 gRPC 服务的连接
	conn, err := grpc.Dial(
		//"localhost:50051",
		"sh.gofxq.com:50051",
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

	go ReportMinitorInfo(conn)

	select {}

}

func ReportHostInfo(conn *grpc.ClientConn) {
	c := pb.NewHostServiceClient(conn)
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

func ReportMinitorInfo(conn *grpc.ClientConn) {
	interval := time.Second
	c := pb.NewMonitorServiceClient(conn)
	clientInfo, err := c.ReportMonitorStream(context.Background())
	if err != nil {
		log.Fatalln("init ReportHostInfoStream failed")
	}
	t := time.Tick(interval)
	for i := range t {
		info := monitor.GetMonitor(context.TODO())
		log.Println(i, time.Now())
		// 调用 Report 方法
		err := clientInfo.Send(info)
		if err != nil {
			log.Printf("could not report: %v", err)
			time.Sleep(interval * 10)
		}
	}

}
