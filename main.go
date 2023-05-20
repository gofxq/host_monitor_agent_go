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
	c := pb.NewMonitorServiceClient(conn)

	// 准备请求数据
	hostInfo := monitor.GetHost(context.TODO())

	client, err := c.ReportHostInfoStream(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(time.Now())
	t := time.Tick(time.Second)
	for i := range t {
		log.Println(i, time.Now())

		// 调用 Report 方法
		err := client.Send(hostInfo)
		if err != nil {
			log.Printf("could not report: %v", err)
			time.Sleep(time.Second * 10)
		}
	}

}
