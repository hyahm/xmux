package xmux

import (
	"fmt"
	"net/http"
	"testing"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home admin" + Var(r)["bbb"]))
}

func grouphome(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("grouphome" + Var(r)["name"] + "-" + Var(r)["age"]))
	fmt.Println(1111)
}

func adminGroup() *RouteGroup {
	admin := NewRouteGroup().Prefix("test").DelPrefix("aa")
	admin.Get("/admin/{bbb}", home).Prefix("aa").DelPrefix("bb")
	admin.Get("/aaa/adf{re:([a-z]{1,4})sf([0-9]{0,10})sd: name, age}", grouphome)
	return admin
}

func userGroup() *RouteGroup {
	user := NewRouteGroup().Prefix("test")
	user.Get("/group", grouphome)
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
	router.Get("/", home).BindResponse(nil)
	// router.Get("/static", tt)
	router.AddGroup(FileBrowse("/static", "D:\\ProgramData", true, false))
	router.AddGroup(userGroup())
	router.DebugTpl()
	GenerateCertificate("cert.pem", "key.pem", "localhost")
	err := router.RunQuic("cert.pem", "key.pem")
	fmt.Println(err)

}
