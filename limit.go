package xmux

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(100, 200)

// 限流中间件（兼容 xmux 的 http.Handler 风格）
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type FixedWindowCounter struct {
	windowSize    time.Duration // 窗口大小
	maxRequests   int           // 每个窗口的最大请求数
	currentWindow time.Time     // 当前窗口的开始时间
	currentCount  int           // 当前窗口的请求计数
	mu            sync.Mutex    // 互斥锁，保证线程安全
}

func NewFixedWindowCounter(windowSize time.Duration, maxRequests int) *FixedWindowCounter {
	return &FixedWindowCounter{
		windowSize:    windowSize,
		maxRequests:   maxRequests,
		currentWindow: time.Now(),
		currentCount:  0,
	}
}

func (fwc *FixedWindowCounter) Allow() bool {
	fwc.mu.Lock()
	defer fwc.mu.Unlock()

	now := time.Now()
	// 检查是否需要切换窗口
	if now.Sub(fwc.currentWindow) >= fwc.windowSize {
		fwc.currentWindow = now
		fwc.currentCount = 0
	}

	// 检查是否超过最大请求数
	if fwc.currentCount < fwc.maxRequests {
		fwc.currentCount++
		return true
	}

	return false
}

var counter *FixedWindowCounter

func init() {
	counter = NewFixedWindowCounter(1*time.Second, 1000) // 每秒最多允许1000个请求
}

func LimitFixedWindowCounterTemplate(w http.ResponseWriter, r *http.Request) (exit bool) {
	// GetConnents() 是全局连接数
	if !counter.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		return true
	}
	return
}

func LimitFixedWindowCounterTemplate1(w http.ResponseWriter, r *http.Request) (exit bool) {
	// GetConnents() 是全局连接数
	if GetConnents() > 1000 {
		w.WriteHeader(http.StatusTooManyRequests)
		return true
	}
	return
}
