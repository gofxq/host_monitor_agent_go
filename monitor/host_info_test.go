package monitor

import (
	"context"
	pb "github.com/gofxq/host_monitor_agent_go/protos/protos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHost(t *testing.T) {
	type args struct {
		c context.Context
	}
	tests := []struct {
		name string
		args args
		want *pb.TickHostInfo
	}{
		{
			name: "",
			args: args{
				c: context.TODO(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetHost(tt.args.c), "GetHost(%v)", tt.args.c)
			t.Logf("%#v", GetHost(tt.args.c))
		})
	}
}
