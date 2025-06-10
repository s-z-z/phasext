package grpctry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"k8s.io/klog/v2"
)

var ErrNoHealthyConn = fmt.Errorf("no healthy connections available")

// ConnPool 封装 gRPC 连接池
type ConnPool struct {
	conns         []*grpc.ClientConn
	mu            sync.RWMutex
	addrs         []string
	dialOpts      []grpc.DialOption
	maxRetry      int
	timeOutSecond int
}

// NewConnPool 创建连接池
func NewConnPool(addrs []string, maxRetry int, timeOutSecond int, dialOpts ...grpc.DialOption) (*ConnPool, error) {
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no addresses provided")
	}
	dialOpts = append(dialOpts, grpc.WithBlock())
	return &ConnPool{
		addrs:         addrs,
		dialOpts:      dialOpts,
		maxRetry:      maxRetry,
		timeOutSecond: timeOutSecond,
		mu:            sync.RWMutex{},
		conns:         make([]*grpc.ClientConn, len(addrs)),
	}, nil
}

// GetConn 获取一个健康的连接
func (p *ConnPool) GetConn(ctx context.Context) (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, addr := range p.addrs {
		// 如果已有连接，检查是否健康
		if conn := p.conns[i]; conn != nil {
			if p.isHealthy(conn) {
				return conn, nil
			}
			// 关闭不健康的连接
			conn.Close()
			p.conns[i] = nil
		}

		// 尝试拨号
		conn, err := p.dial(ctx, addr)
		if err != nil {
			continue
		}
		p.conns[i] = conn
		return conn, nil
	}

	return nil, ErrNoHealthyConn
}

// isHealthy 检查连接是否健康
func (p *ConnPool) isHealthy(conn *grpc.ClientConn) bool {
	if conn == nil {
		return false
	}
	return conn.GetState() == connectivity.Ready
}

// dial 拨号
func (p *ConnPool) dial(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < p.maxRetry; i++ {

		klog.V(8).Infof("dialing %s\n", addr)

		dialCtx, cancel := context.WithTimeout(ctx, time.Duration(p.timeOutSecond)*time.Second)
		conn, err = grpc.DialContext(dialCtx, addr, p.dialOpts...)
		cancel()
		if err == nil && conn.GetState() == connectivity.Ready {
			return conn, nil
		}
		if conn != nil {
			conn.Close()
		}
	}
	return nil, fmt.Errorf("failed to dial %s after %d attempts: %v", addr, p.maxRetry, err)
}

// Close 关闭所有连接
func (p *ConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conn := range p.conns {
		if conn != nil {
			conn.Close()
		}
	}
	p.conns = nil
}
