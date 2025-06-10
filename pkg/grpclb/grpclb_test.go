package grpclb_test

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/suzi1037/phasext/pkg/grpclb"
)

func main() {
	ctx := context.Background()

	config := grpclb.BalancerConfig{
		Endpoints:     []string{"localhost:50051", "localhost:50052", "localhost:50053"},
		HealthTimeout: 5 * time.Second,
		CheckInterval: 10 * time.Second,
		DialOptions:   []grpc.DialOption{
			// 添加额外的拨号选项
		},
	}

	balancer, err := grpclb.NewBalancer(ctx, config)
	if err != nil {
		fmt.Printf("Failed to create balancer: %v\n", err)
		return
	}
	defer balancer.Close()

	// 获取连接示例
	//conn, err := balancer.GetConn()
	//if err != nil {
	//	fmt.Printf("Failed to get connection: %v\n", err)
	//	return
	//}

	// 使用 conn 进行 gRPC 调用
	// client := your_service.NewYourServiceClient(conn)
	// ...
}
