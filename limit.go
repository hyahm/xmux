package xmux

import (
	"net/http"
	"sync"
	"time"
)

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

func LimitFixedWindowCounterTemplate(w http.ResponseWriter, r *http.Request) bool {
	// GetConnents() 是全局连接数
	if !counter.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		return true
	}
	return false
}

func LimitFixedWindowCounterTemplate1(w http.ResponseWriter, r *http.Request) bool {
	// GetConnents() 是全局连接数
	if GetConnents() > 1000 {
		w.WriteHeader(http.StatusTooManyRequests)
		return true
	}
	return false
}

// var swc SlidingWindowCounter

// func init() {
// 	swc = SlidingWindowCounter{}
// }

// type SlidingWindowCounter struct {
// 	mapWindow map[int64]int
// 	mu        sync.RWMutex
// }

// func AddSlidingWindowCounter(t int64) bool {
// 	swc.mu.Lock()
// 	defer swc.mu.Unlock()
// 	if _, ok := swc.mapWindow[t]; !ok {
// 		swc.mapWindow[t] = 0
// 	}
// 	if swc.mapWindow[t] >= 10 {
// 		return false
// 	}
// 	swc.mapWindow[t]++
// 	return true
// }

// func DelSlidingWindowCounter(t int64) {
// 	swc.mu.Lock()
// 	defer swc.mu.Unlock()
// 	if _, ok := swc.mapWindow[t]; !ok {
// 		swc.mapWindow[t] = 0
// 	}
// 	swc.mapWindow[t]--
// 	if swc.mapWindow[t] <= 0 {
// 		delete(swc.mapWindow, t)
// 	}
// }
// func LimitSlidingWindowCounterTemplate(w http.ResponseWriter, r *http.Request) bool {
// 	// GetInstance(r).GetConnectId()   是请求的连接ID， 默认是当前请求的纳秒时间, 除以100000000 后是每0.1秒请求数
// 	return !AddSlidingWindowCounter(GetInstance(r).GetConnectId() / 100000000)
// }

// type LeakyBucket struct {
// 	capacity int           // 桶的容量
// 	rate     int           // 桶的流出速率（每秒流出的数据包数）
// 	bucket   chan struct{} // 桶，用通道模拟
// 	mu       sync.Mutex    // 互斥锁，保证线程安全
// }

// func NewLeakyBucket(capacity, rate int) *LeakyBucket {
// 	bucket := make(chan struct{}, capacity)
// 	for i := 0; i < capacity; i++ {
// 		bucket <- struct{}{}
// 	}
// 	return &LeakyBucket{
// 		capacity: capacity,
// 		rate:     rate,
// 		bucket:   bucket,
// 	}
// }

// func (lb *LeakyBucket) Allow() bool {
// 	lb.mu.Lock()
// 	defer lb.mu.Unlock()

// 	select {
// 	case lb.bucket <- struct{}{}:
// 		return true
// 	default:
// 		return false
// 	}
// }

// var bucket *LeakyBucket

// func init() {
// 	bucket = NewLeakyBucket(10, 5)
// 	go bucket.leak()
// }

// func (lb *LeakyBucket) leak() {
// 	ticker := time.NewTicker(time.Duration(1000/lb.rate) * time.Millisecond)
// 	for range ticker.C {
// 		<-lb.bucket
// 	}
// }

// func LimitLeakyBucketTemplate(w http.ResponseWriter, r *http.Request) bool {
// 	// 不允许就丢弃
// 	return !bucket.Allow()
// }

// type TokenBucket struct {
// 	capacity int           // 桶的容量
// 	rate     int           // 令牌生成速率（每秒生成的令牌数）
// 	bucket   chan struct{} // 桶，用通道模拟
// 	mu       sync.Mutex    // 互斥锁，保证线程安全
// }

// func NewTokenBucket(capacity, rate int) *TokenBucket {
// 	bucket := make(chan struct{}, capacity)
// 	for i := 0; i < capacity; i++ {
// 		bucket <- struct{}{}
// 	}
// 	return &TokenBucket{
// 		capacity: capacity,
// 		rate:     rate,
// 		bucket:   bucket,
// 	}
// }

// func (tb *TokenBucket) Allow() bool {
// 	tb.mu.Lock()
// 	defer tb.mu.Unlock()

// 	select {
// 	case tb.bucket <- struct{}{}:
// 		return true
// 	default:
// 		return false
// 	}
// }

// func (tb *TokenBucket) refill() {
// 	ticker := time.NewTicker(time.Duration(1000/tb.rate) * time.Millisecond)
// 	for range ticker.C {
// 		select {
// 		case tb.bucket <- struct{}{}:
// 		default:
// 		}
// 	}
// }

// var token *TokenBucket

// func init() {
// 	token = NewTokenBucket(10, 5)
// 	go token.refill()
// }

// func LimitTokenBucketTemplate(w http.ResponseWriter, r *http.Request) bool {
// 	// 不允许就丢弃
// 	return !token.Allow()
// }
