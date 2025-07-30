package grpctry_test

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/s-z-z/phasext/pkg/grpctry"
)

func main() {
	// 定义多个 gRPC 服务地址
	addrs := []string{
		"ip1:50051",
		"ip2:50052",
		"ip3:50053",
	}

	// 创建连接池
	pool, err := grpctry.NewConnPool(addrs, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// 获取健康连接
	conn, err := pool.GetConn(context.Background())
	if err != nil {
		log.Fatalf("failed to get healthy connection: %v", err)
	}
	fmt.Println("Got healthy connection:", conn.Target())

	// 使用 conn 进行 gRPC 调用
	// 例如：client := pb.NewYourServiceClient(conn)
}
