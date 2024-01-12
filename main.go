package main

import (
	"flag"
	"fmt"
	"github.com/gofxq/host_monitor_agent_go/monitor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {

	// 后面三个参数分别是命令行参数名、默认值、帮助信息
	flag.StringVar(&monitor.ServerID, "id", "test_server", "Report Server ID")
	flag.StringVar(&monitor.GrpcServerAddr, "addr", "192.168.9.13:8008", "GRPC Server Addr")
	// 解析命令行参数
	flag.Parse()

	fmt.Printf("report %s to %s. \n", monitor.ServerID, monitor.GrpcServerAddr)

	// 设置一个连接到 gRPC 服务的连接
	conn, err := grpc.Dial(
		monitor.GrpcServerAddr,
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

	go monitor.ReportHostInfo(conn)

	go monitor.ReportMonitorInfo(conn)

	select {}
}
