package xmux

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/time/rate"
)

// 基于module增加超时， 与全局的不一样，这里是会对传入的 module 有效
func SetTimeout(d time.Duration, m ...func(http.ResponseWriter, *http.Request) bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		// time.Sleep(time.Second * 4)
		for _, v := range m {
			select {
			case <-ctx.Done():
				return
			default:
				if v(w, r) {
					return
				}

			}
		}
	}
}

var (
	ipLimiters = make(map[string]*rate.Limiter)
	limitMu    sync.Mutex
)

// RateLimit 限流：每秒生成 r 个令牌，最大突发 b
func RateLimit(r rate.Limit, b int) func(w http.ResponseWriter, r *http.Request) bool {
	return func(w http.ResponseWriter, req *http.Request) bool {
		ip := req.RemoteAddr

		limitMu.Lock()
		l, ok := ipLimiters[ip]
		if !ok {
			l = rate.NewLimiter(r, b)
			ipLimiters[ip] = l
		}
		limitMu.Unlock()

		if !l.Allow() {
			http.Error(w, "too many requests", 429)
			return true
		}
		return false
	}
}

func CircuitBreakerTemplate(name string, timeout time.Duration) func(w http.ResponseWriter, r *http.Request) bool {
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:               int(timeout),
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 50,
	})

	return func(w http.ResponseWriter, r *http.Request) bool {
		err := hystrix.Do(name, func() error {
			// 什么都不做，让后续 module/handler 继续执行
			return nil
		}, nil)

		if err != nil {
			http.Error(w, "service unavailable", 503)
			return true
		}
		return false
	}
}

// FlowData 内置数据实例
// type FlowData struct {
// 	Data map[string]interface{}
// }

// func GetInstance(r *http.Request) *FlowData {
// 	return &FlowData{Data: make(map[string]interface{})}
// }
