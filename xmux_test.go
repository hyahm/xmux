package xmux

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
}

func grouphome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("grouphome")
}

var group *RouteGroup

func init() {
	group = NewRouteGroup().Prefix("test")
	group.Get("/group", grouphome)
}

func TestMain(t *testing.T) {
	router := NewRouter()
	router.AddGroup(Pprof())
	router.EnableConnect = true
	router.Get("/pp", home)
	router.AddGroup(group)
	log.Fatal(router.Run())
}
