package xmux

import (
	"log"
	"net/http"
	"testing"
)

func home(w http.ResponseWriter, r *http.Request) {
	name := Var(r)["name"]
	w.Write([]byte("home admin" + name))
}

func grouphome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("grouphome" + Var(r)["name"] + "-" + Var(r)["age"]))
}

func adminGroup() *RouteGroup {
	admin := NewRouteGroup().Prefix("test")
	admin.Get("/admin/{bbb}", home)
	admin.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd: name, age}", grouphome)
	return admin
}

func userGroup() *RouteGroup {
	user := NewRouteGroup().Prefix("test")
	user.Get("/group", grouphome)
	user.AddGroup(adminGroup())
	return user
}

func TestMain(t *testing.T) {
	router := NewRouter()
	router.AddGroup(Pprof())
	router.EnableConnect = true
	router.Get("/pp/{name}", home)
	router.SetAddr(":9000")

	router.AddGroup(userGroup())
	log.Fatal(router.Run())
}
