package xmux

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func home(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("home admin" + Var(r)["bbb"]))
}

func grouphome(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("grouphome" + Var(r)["name"] + "-" + Var(r)["age"]))
	fmt.Println(1111)
}

func adminGroup() *RouteGroup {
	admin := NewRouteGroup().Prefix("test")
	admin.Get("/admin/{bbb}", home)
	admin.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd: name, age}", grouphome)
	return admin
}

func userGroup() *RouteGroup {
	user := NewRouteGroup().Prefix("test")
	user.Get("/group", grouphome).BindResponse(nil)
	user.Request("/group/add", nil, http.MethodGet, http.MethodDelete, http.MethodPost)
	user.AddGroup(adminGroup())
	return user
}

type ResponseParameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
}

func TestMain(t *testing.T) {
	response := &ResponseParameter{
		Name:     "name",
		In:       "query",
		Required: true,
		Type:     "string",
	}
	router := NewRouter().BindResponse(response)

	router.AddGroup(Pprof())
	router.EnableConnect = true
	router.Get("/pp", home)
	router.SetAddr(":9000")

	router.AddGroup(userGroup())
	router.DebugRoute()
	log.Fatal(router.Run())
}
