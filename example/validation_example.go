package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyahm/xmux"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Phone    string `json:"phone"`
	Bio      string `json:"bio"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	
	validator := xmux.NewValidator().
		AddField("username", 
			xmux.Required(),
			xmux.MinLength(3),
			xmux.MaxLength(20),
			xmux.Alphanumeric(),
		).
		AddField("email",
			xmux.Required(),
			xmux.Email(),
		).
		AddField("age",
			xmux.Required(),
			xmux.Min(18.0),
			xmux.Max(120.0),
		).
		AddField("phone",
			xmux.Required(),
			xmux.Phone(),
		).
		AddField("bio",
			xmux.MaxLength(500),
			xmux.NoXSS(),
		)
	
	if xmux.ValidateAndResponse(w, r, user, validator) {
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

func createPost(w http.ResponseWriter, r *http.Request) {
	type Post struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
		Status  string   `json:"status"`
	}
	
	post := &Post{}
	
	validator := xmux.NewValidator().
		AddField("title",
			xmux.Required(),
			xmux.MinLength(5),
			xmux.MaxLength(100),
			xmux.NoXSS(),
		).
		AddField("content",
			xmux.Required(),
			xmux.MinLength(10),
			xmux.MaxLength(10000),
			xmux.NoXSS(),
			xmux.NoSQLInjection(),
		).
		AddField("tags",
			xmux.Required(),
		).
		AddField("status",
			xmux.Required(),
			xmux.OneOf("draft", "published", "archived"),
		)
	
	if xmux.ValidateAndResponse(w, r, post, validator) {
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created successfully",
		"post":    post,
	})
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	form := &LoginForm{}
	
	validator := xmux.NewValidator().
		AddField("username",
			xmux.Required(),
			xmux.MinLength(3),
			xmux.MaxLength(50),
		).
		AddField("password",
			xmux.Required(),
			xmux.MinLength(8),
			xmux.MaxLength(100),
		)
	
	if xmux.ValidateAndResponse(w, r, form, validator) {
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   "fake-jwt-token",
	})
}

type SearchQuery struct {
	Keyword  string `json:"keyword"`
	Category string `json:"category"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func search(w http.ResponseWriter, r *http.Request) {
	query := &SearchQuery{}
	
	validator := xmux.NewValidator().
		AddField("keyword",
			xmux.Required(),
			xmux.MaxLength(100),
			xmux.NoXSS(),
			xmux.NoSQLInjection(),
		).
		AddField("category",
			xmux.MaxLength(50),
		).
		AddField("page",
			xmux.Min(1.0),
		).
		AddField("page_size",
			xmux.Min(1.0),
			xmux.Max(100.0),
		)
	
	if xmux.ValidateAndResponse(w, r, query, validator) {
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": []string{},
		"page":    query.Page,
		"total":   0,
	})
}

func main() {
	router := xmux.NewRouter()
	
	router.SetHeader("Access-Control-Allow-Origin", "*")
	router.SetHeader("Access-Control-Allow-Headers", "Content-Type,Authorization")
	
	securityConfig := &xmux.SecurityConfig{
		EnableRequestSizeLimit: true,
		EnableHeaderCheck:     true,
		EnablePathTraversalCheck: true,
		MaxRequestSize:        10 << 20,
		MaxHeaderSize:         8192,
		AllowedMethods:        []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		BlockedPaths:          []string{`/admin`, `/\.env`, `/config`},
	}
	
	securityMiddleware := xmux.NewSecurityMiddleware(securityConfig)
	
	router.AddModule(securityMiddleware.SecurityCheck)
	
	router.Post("/users", createUser).
		BindJson(&User{})
	
	router.Post("/login", login).
		BindJson(&LoginForm{})
	
	router.Post("/posts", createPost)
	
	router.Get("/search", search)
	
	fmt.Println("Server starting on :8080")
	router.Run()
}
