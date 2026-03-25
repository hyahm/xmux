package xmux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type SecurityConfig struct {
	EnableIPWhitelist    bool
	EnableIPBlacklist    bool
	EnableOriginCheck    bool
	EnableRateLimiting   bool
	EnableRequestSizeLimit bool
	EnableHeaderCheck    bool
	EnablePathTraversalCheck bool
	
	IPWhitelist         []string
	IPBlacklist         []string
	AllowedOrigins      []string
	MaxRequestSize      int64
	MaxHeaderSize       int64
	AllowedMethods      []string
	AllowedHeaders      []string
	BlockedUserAgents   []string
	BlockedPaths        []string
}

type SecurityMiddleware struct {
	config      *SecurityConfig
	ipWhitelist map[string]bool
	ipBlacklist map[string]bool
	originRegex *regexp.Regexp
	userAgentRegexes []*regexp.Regexp
	pathRegexes      []*regexp.Regexp
	mu          sync.RWMutex
}

func NewSecurityMiddleware(config *SecurityConfig) *SecurityMiddleware {
	if config == nil {
		config = &SecurityConfig{
			MaxRequestSize: 10 << 20, // 10MB
			MaxHeaderSize:  8192,     // 8KB
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		}
	}
	
	sm := &SecurityMiddleware{
		config:      config,
		ipWhitelist: make(map[string]bool),
		ipBlacklist: make(map[string]bool),
	}
	
	for _, ip := range config.IPWhitelist {
		sm.ipWhitelist[ip] = true
	}
	
	for _, ip := range config.IPBlacklist {
		sm.ipBlacklist[ip] = true
	}
	
	if len(config.AllowedOrigins) > 0 {
		pattern := strings.Join(config.AllowedOrigins, "|")
		sm.originRegex = regexp.MustCompile(pattern)
	}
	
	for _, ua := range config.BlockedUserAgents {
		sm.userAgentRegexes = append(sm.userAgentRegexes, regexp.MustCompile(ua))
	}
	
	for _, path := range config.BlockedPaths {
		sm.pathRegexes = append(sm.pathRegexes, regexp.MustCompile(path))
	}
	
	return sm
}

func (sm *SecurityMiddleware) CheckIP(w http.ResponseWriter, r *http.Request) bool {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	if sm.config.EnableIPBlacklist {
		if sm.ipBlacklist[ip] {
			w.WriteHeader(http.StatusForbidden)
			return true
		}
	}
	
	if sm.config.EnableIPWhitelist {
		if !sm.ipWhitelist[ip] {
			w.WriteHeader(http.StatusForbidden)
			return true
		}
	}
	
	return false
}

func (sm *SecurityMiddleware) CheckOrigin(w http.ResponseWriter, r *http.Request) bool {
	if !sm.config.EnableOriginCheck {
		return false
	}
	
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	if sm.originRegex != nil {
		if !sm.originRegex.MatchString(origin) {
			w.WriteHeader(http.StatusForbidden)
			return true
		}
	}
	
	return false
}

func (sm *SecurityMiddleware) CheckUserAgent(w http.ResponseWriter, r *http.Request) bool {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return false
	}
	
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	for _, regex := range sm.userAgentRegexes {
		if regex.MatchString(userAgent) {
			w.WriteHeader(http.StatusForbidden)
			return true
		}
	}
	
	return false
}

func (sm *SecurityMiddleware) CheckPath(w http.ResponseWriter, r *http.Request) bool {
	if !sm.config.EnablePathTraversalCheck {
		return false
	}
	
	path := r.URL.Path
	
	if strings.Contains(path, "../") || strings.Contains(path, "..\\") {
		w.WriteHeader(http.StatusBadRequest)
		return true
	}
	
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	for _, regex := range sm.pathRegexes {
		if regex.MatchString(path) {
			w.WriteHeader(http.StatusForbidden)
			return true
		}
	}
	
	return false
}

func (sm *SecurityMiddleware) CheckRequestSize(w http.ResponseWriter, r *http.Request) bool {
	if !sm.config.EnableRequestSizeLimit {
		return false
	}
	
	if r.ContentLength > sm.config.MaxRequestSize {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return true
	}
	
	return false
}

func (sm *SecurityMiddleware) CheckHeaders(w http.ResponseWriter, r *http.Request) bool {
	if !sm.config.EnableHeaderCheck {
		return false
	}
	
	for name, values := range r.Header {
		for _, value := range values {
			if len(value) > int(sm.config.MaxHeaderSize) {
				w.WriteHeader(http.StatusBadRequest)
				return true
			}
			
			if sm.isSuspiciousHeader(name, value) {
				w.WriteHeader(http.StatusBadRequest)
				return true
			}
		}
	}
	
	return false
}

func (sm *SecurityMiddleware) isSuspiciousHeader(name, value string) bool {
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"onerror=",
		"onload=",
		"eval(",
		"document.cookie",
	}
	
	lowerValue := strings.ToLower(value)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerValue, pattern) {
			return true
		}
	}
	
	return false
}

func (sm *SecurityMiddleware) CheckMethod(w http.ResponseWriter, r *http.Request) bool {
	if len(sm.config.AllowedMethods) == 0 {
		return false
	}
	
	for _, method := range sm.config.AllowedMethods {
		if r.Method == method {
			return false
		}
	}
	
	w.WriteHeader(http.StatusMethodNotAllowed)
	return true
}

func (sm *SecurityMiddleware) CheckContentType(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return false
	}
	
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return false
	}
	
	allowedTypes := []string{
		"application/json",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
		"text/plain",
		"application/xml",
		"text/xml",
	}
	
	for _, allowedType := range allowedTypes {
		if strings.HasPrefix(contentType, allowedType) {
			return false
		}
	}
	
	w.WriteHeader(http.StatusUnsupportedMediaType)
	return true
}

func (sm *SecurityMiddleware) SecurityCheck(w http.ResponseWriter, r *http.Request) bool {
	if sm.CheckMethod(w, r) {
		return true
	}
	
	if sm.CheckPath(w, r) {
		return true
	}
	
	if sm.CheckIP(w, r) {
		return true
	}
	
	if sm.CheckOrigin(w, r) {
		return true
	}
	
	if sm.CheckUserAgent(w, r) {
		return true
	}
	
	if sm.CheckRequestSize(w, r) {
		return true
	}
	
	if sm.CheckHeaders(w, r) {
		return true
	}
	
	if sm.CheckContentType(w, r) {
		return true
	}
	
	return false
}

func SecurityMiddlewareTemplate(w http.ResponseWriter, r *http.Request) bool {
	config := &SecurityConfig{
		EnableIPWhitelist:    false,
		EnableIPBlacklist:    false,
		EnableOriginCheck:    false,
		EnableRateLimiting:   false,
		EnableRequestSizeLimit: true,
		EnableHeaderCheck:    true,
		EnablePathTraversalCheck: true,
		MaxRequestSize:       10 << 20,
		MaxHeaderSize:        8192,
		AllowedMethods:       []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:       []string{"Content-Type", "Authorization", "X-Requested-With"},
		BlockedPaths:         []string{`/admin`, `/\.env`, `/config`},
	}
	
	sm := NewSecurityMiddleware(config)
	return sm.SecurityCheck(w, r)
}

type WebSocketSecurityConfig struct {
	AllowedOrigins   []string
	MaxMessageSize   int64
	EnableRateLimit  bool
	EnableAuthCheck  bool
	AuthHeader       string
	AllowedProtocols []string
}

type WebSocketSecurity struct {
	config      *WebSocketSecurityConfig
	originRegex *regexp.Regexp
	mu          sync.RWMutex
}

func NewWebSocketSecurity(config *WebSocketSecurityConfig) *WebSocketSecurity {
	if config == nil {
		config = &WebSocketSecurityConfig{
			MaxMessageSize:  1 << 20, // 1MB
			AllowedProtocols: []string{"chat", "notification"},
		}
	}
	
	ws := &WebSocketSecurity{
		config: config,
	}
	
	if len(config.AllowedOrigins) > 0 {
		pattern := strings.Join(config.AllowedOrigins, "|")
		ws.originRegex = regexp.MustCompile(pattern)
	}
	
	return ws
}

func (ws *WebSocketSecurity) CheckUpgrade(w http.ResponseWriter, r *http.Request) error {
	origin := r.Header.Get("Origin")
	if origin != "" && ws.originRegex != nil {
		ws.mu.RLock()
		matched := ws.originRegex.MatchString(origin)
		ws.mu.RUnlock()
		
		if !matched {
			return fmt.Errorf("origin not allowed: %s", origin)
		}
	}
	
	if ws.config.EnableAuthCheck {
		auth := r.Header.Get(ws.config.AuthHeader)
		if auth == "" {
			return fmt.Errorf("missing authorization header")
		}
	}
	
	return nil
}

func (ws *WebSocketSecurity) ValidateMessageSize(size int64) error {
	if size > ws.config.MaxMessageSize {
		return fmt.Errorf("message size exceeds limit: %d > %d", size, ws.config.MaxMessageSize)
	}
	return nil
}

func SecureUpgradeWebSocket(w http.ResponseWriter, r *http.Request, wsConfig *WebSocketSecurityConfig) (*BaseWs, error) {
	wsSecurity := NewWebSocketSecurity(wsConfig)
	
	if err := wsSecurity.CheckUpgrade(w, r); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return nil, err
	}
	
	return UpgradeWebSocket(w, r)
}

func SanitizeJSON(input []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, err
	}
	
	sanitized := sanitizeValue(data)
	return json.Marshal(sanitized)
}

func sanitizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		sanitized := make(map[string]interface{})
		for key, value := range val {
			if isSafeKey(key) {
				sanitized[key] = sanitizeValue(value)
			}
		}
		return sanitized
	case []interface{}:
		sanitized := make([]interface{}, len(val))
		for i, value := range val {
			sanitized[i] = sanitizeValue(value)
		}
		return sanitized
	case string:
		return sanitizeString(val)
	default:
		return v
	}
}

func isSafeKey(key string) bool {
	blacklistedKeys := []string{
		"__proto__", "constructor", "prototype",
		"$where", "$ne", "$gt", "$lt", "$regex",
	}
	
	for _, blacklisted := range blacklistedKeys {
		if strings.Contains(key, blacklisted) {
			return false
		}
	}
	
	return true
}

func sanitizeString(s string) string {
	s = strings.TrimSpace(s)
	
	replacer := strings.NewReplacer(
		"\x00", "",
		"\r", "",
		"\n", " ",
		"\t", " ",
	)
	
	s = replacer.Replace(s)
	
	return s
}

func ValidateURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	return err == nil
}

func ValidateEmail(input string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(input)
}

func ValidatePhone(input string) bool {
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(input)
}

func CheckSQLInjection(input string) bool {
	sqlPatterns := []string{
		`['";]`,
		`\b(OR|AND)\s+\d+\s*=\s*\d+`,
		`\b(UNION|SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|EXEC)\b`,
		`--`,
		`/\*.*\*/`,
		`\b(WHERE|HAVING)\s+\w+\s*=\s*\w+`,
	}
	
	combinedPattern := strings.Join(sqlPatterns, "|")
	regex := regexp.MustCompile(`(?i)` + combinedPattern)
	
	return regex.MatchString(input)
}

func CheckXSS(input string) bool {
	xssPatterns := []string{
		`<script[^>]*>.*?</script>`,
		`javascript:`,
		`on\w+\s*=`,
		`<iframe`,
		`<object`,
		`<embed`,
		`eval\s*\(`,
		`expression\s*\(`,
	}
	
	combinedPattern := strings.Join(xssPatterns, "|")
	regex := regexp.MustCompile(`(?i)` + combinedPattern)
	
	return regex.MatchString(input)
}

func LimitRequestBodySize(maxSize int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maxSize {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("Request body too large"))
			return
		}
		
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		}
		
		next.ServeHTTP(w, r)
	})
}

func SanitizeRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			next.ServeHTTP(w, r)
			return
		}
		
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			
			sanitized, err := SanitizeJSON(body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			
			r.Body = io.NopCloser(bytes.NewBuffer(sanitized))
		}
		
		next.ServeHTTP(w, r)
	})
}

func AddSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		next.ServeHTTP(w, r)
	})
}
