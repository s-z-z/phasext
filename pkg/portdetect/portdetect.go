package portdetect

import (
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

func DefaultRunAfterPortOpenS(addrs string) error {
	interval := 3 * time.Second
	maxRetries := 20
	return RunAfterPortOpenS(interval, maxRetries, addrs)
}

func RunAfterPortOpenS(interval time.Duration, maxRetries int, addrs string) error {
	addrs = strings.TrimSpace(addrs)
	if addrs == "" {
		return nil
	}
	addrList := strings.Split(addrs, ",")
	return RunAfterPortOpen(interval, maxRetries, addrList...)
}

func RunAfterPortOpen(interval time.Duration, maxRetries int, addrs ...string) error {
	eg := errgroup.Group{}

	for _, addr := range addrs {
		eg.Go(func() error {
			for i := 0; i < maxRetries; i++ {
				conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
				if err == nil {
					_ = conn.Close()
					fmt.Printf("Port %s is open\n", addr)
					return nil
				}
				fmt.Printf("Attempt %d: Port %s is not open, retrying in %v...\n", i+1, addr, interval)
				time.Sleep(interval)
			}
			return fmt.Errorf("port %s is not open", addr)
		})
	}
	return eg.Wait()
}
