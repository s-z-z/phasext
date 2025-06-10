package grpclb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Balancer struct {
	endpoints     []string           // 后端地址列表
	conns         []*grpc.ClientConn // 连接池
	healthy       []bool             // 健康状态
	currentIndex  uint32             // 当前连接索引
	mutex         sync.Mutex         // 同步锁
	healthTimeout time.Duration      // 健康检查超时
	checkInterval time.Duration      // 健康检查间隔
}

// BalancerConfig 负载均衡器配置
type BalancerConfig struct {
	Endpoints     []string
	HealthTimeout time.Duration
	CheckInterval time.Duration
	DialOptions   []grpc.DialOption
}

// NewBalancer 创建新的负载均衡器
func NewBalancer(ctx context.Context, config BalancerConfig) (*Balancer, error) {
	if len(config.Endpoints) == 0 {
		return nil, fmt.Errorf("at least one endpoint is required")
	}

	if config.HealthTimeout == 0 {
		config.HealthTimeout = 5 * time.Second
	}
	if config.CheckInterval == 0 {
		config.CheckInterval = 10 * time.Second
	}

	b := &Balancer{
		endpoints:     config.Endpoints,
		conns:         make([]*grpc.ClientConn, len(config.Endpoints)),
		healthy:       make([]bool, len(config.Endpoints)),
		healthTimeout: config.HealthTimeout,
		checkInterval: config.CheckInterval,
	}

	// 默认的拨号选项
	defaultOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if len(config.DialOptions) > 0 {
		defaultOptions = append(defaultOptions, config.DialOptions...)
	}

	// 初始化所有连接
	for i, endpoint := range config.Endpoints {
		conn, err := grpc.DialContext(ctx, endpoint, defaultOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to dial %s: %v", endpoint, err)
		}
		b.conns[i] = conn
		b.healthy[i] = true // 初始都假设健康
	}

	// 启动健康检查
	go b.startHealthCheck(ctx)

	return b, nil
}

// GetConn 获取一个健康的连接（轮询策略）
func (b *Balancer) GetConn() (*grpc.ClientConn, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// 最多尝试所有连接
	for i := 0; i < len(b.conns); i++ {
		// 原子操作获取并更新索引
		index := b.currentIndex % uint32(len(b.conns))
		b.currentIndex++

		if b.healthy[index] {
			return b.conns[index], nil
		}
	}

	return nil, fmt.Errorf("no healthy connections available")
}

// startHealthCheck 启动健康检查
func (b *Balancer) startHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(b.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.checkAllEndpoints(ctx)
		}
	}
}

// checkAllEndpoints 检查所有端点的健康状态
func (b *Balancer) checkAllEndpoints(ctx context.Context) {
	for i, conn := range b.conns {
		healthy := b.checkEndpoint(ctx, conn)
		b.mutex.Lock()
		b.healthy[i] = healthy
		b.mutex.Unlock()
	}
}

// checkEndpoint 检查单个端点的健康状态
func (b *Balancer) checkEndpoint(ctx context.Context, conn *grpc.ClientConn) bool {
	healthClient := grpc_health_v1.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(ctx, b.healthTimeout)
	defer cancel()

	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return false
	}

	return resp.Status == grpc_health_v1.HealthCheckResponse_SERVING
}

// Close 关闭所有连接
func (b *Balancer) Close() error {
	var firstErr error
	for _, conn := range b.conns {
		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
