package monitor

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestReportHostInfo(t *testing.T) {
	type args struct {
		conn *grpc.ClientConn
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReportHostInfo(tt.args.conn)
		})
	}
}

func BenchmarkReportHostInfo(b *testing.B) {
	// 设置一个连接到 gRPC 服务的连接
	conn, err := grpc.Dial(
		GrpcServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		ReportHostInfo(conn)
	}
}
