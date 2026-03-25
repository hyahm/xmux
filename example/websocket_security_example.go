package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hyahm/xmux"
)

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Room     string `json:"room"`
}

func chatWebSocket(w http.ResponseWriter, r *http.Request) {
	wsConfig := &xmux.WebSocketSecurityConfig{
		AllowedOrigins:   []string{`^http://localhost:\d+$`, `^https://trusted\.example\.com$`},
		MaxMessageSize:   1 << 20,
		EnableAuthCheck:  true,
		AuthHeader:       "X-Auth-Token",
		AllowedProtocols: []string{"chat"},
	}
	
	ws, err := xmux.SecureUpgradeWebSocket(w, r, wsConfig)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()
	
	log.Printf("WebSocket connection established from %s", ws.RemoteAddr)
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	done := make(chan struct{})
	
	go func() {
		defer close(done)
		for {
			select {
			case <-ticker.C:
				if err := ws.Ping([]byte("keepalive")); err != nil {
					log.Printf("Ping failed: %v", err)
					return
				}
			}
		}
	}()
	
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		
		if msgType == xmux.TypeClose {
			log.Println("Client requested close")
			break
		}
		
		if msgType == xmux.TypePing {
			ws.SendMessage([]byte("pong"), xmux.TypePong)
			continue
		}
		
		if msgType == xmux.TypeMsg {
			var chatMsg ChatMessage
			if err := json.Unmarshal([]byte(msg), &chatMsg); err != nil {
				log.Printf("Invalid message format: %v", err)
				continue
			}
			
			validator := xmux.NewValidator().
				AddField("username", 
					xmux.Required(),
					xmux.MinLength(2),
					xmux.MaxLength(20),
					xmux.Alphanumeric(),
				).
				AddField("message",
					xmux.Required(),
					xmux.MinLength(1),
					xmux.MaxLength(500),
					xmux.NoXSS(),
				).
				AddField("room",
					xmux.Required(),
					xmux.MaxLength(50),
					xmux.Pattern(`^[a-zA-Z0-9_-]+$`),
				)
			
			errors := validator.Validate(&chatMsg)
			if len(errors) > 0 {
				log.Printf("Validation errors: %v", errors)
				errorMsg := map[string]interface{}{
					"error":   "validation_failed",
					"details": errors,
				}
				errorBytes, _ := json.Marshal(errorMsg)
				ws.SendMessage(errorBytes, xmux.TypeMsg)
				continue
			}
			
			response := map[string]interface{}{
				"type":      "message",
				"username":  chatMsg.Username,
				"message":   chatMsg.Message,
				"room":      chatMsg.Room,
				"timestamp": time.Now().Unix(),
			}
			
			responseBytes, _ := json.Marshal(response)
			if err := ws.SendMessage(responseBytes, xmux.TypeMsg); err != nil {
				log.Printf("Send error: %v", err)
				break
			}
		}
	}
	
	<-done
	log.Printf("WebSocket connection closed from %s", ws.RemoteAddr)
}

func main() {
	router := xmux.NewRouter()
	
	router.SetHeader("Access-Control-Allow-Origin", "*")
	
	securityConfig := &xmux.SecurityConfig{
		EnableRequestSizeLimit: true,
		EnableHeaderCheck:     true,
		MaxRequestSize:        10 << 20,
		MaxHeaderSize:         8192,
	}
	
	securityMiddleware := xmux.NewSecurityMiddleware(securityConfig)
	router.AddModule(securityMiddleware.SecurityCheck)
	
	router.Get("/ws/chat", chatWebSocket)
	
	fmt.Println("WebSocket server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
