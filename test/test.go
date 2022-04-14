package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"strconv"
	"strings"

	"github.com/hyahm/xmux"
)

func home1(w http.ResponseWriter, r *http.Request) bool {
	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	w.WriteHeader(404)
	// 	return true
	// }
	// fmt.Println(string(b))
	ct := r.Header.Get("Content-Type")
	ct = strings.ToLower(ct)
	if r.Method == http.MethodGet {
		Form(r)
	}
	if r.Method == http.MethodPost {
		fmt.Println("-------------")
		fmt.Println(r.FormValue("username"))
		Form(r)

		// b, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	w.WriteHeader(404)
		// 	return true
		// }
		// fmt.Println(string(b))
		// err = json.Unmarshal(b, xmux.GetInstance(r).Data)
		// if err != nil {
		// 	w.WriteHeader(404)
		// 	return true
		// }

	}

	return false

}

func Form(r *http.Request) {
	tt := reflect.TypeOf(xmux.GetInstance(r).Data).Elem()
	vv := reflect.ValueOf(xmux.GetInstance(r).Data).Elem()
	len := tt.NumField()
	for i := 0; i < len; i++ {
		key := tt.Field(i).Tag.Get("form")
		value := r.FormValue(key)
		switch tt.Field(i).Type.String() {
		case "string":
			vv.Field(i).Set(reflect.ValueOf(value))
		case "float32":
			b32, _ := strconv.ParseFloat(value, 32)
			vv.Field(i).Set(reflect.ValueOf(b32))
		case "float64":
			b64, _ := strconv.ParseFloat(value, 64)
			vv.Field(i).Set(reflect.ValueOf(b64))
		case "bool":
			ok, _ := strconv.ParseBool(value)
			vv.Field(i).Set(reflect.ValueOf(ok))
		case "int":
			i, _ := strconv.Atoi(value)
			vv.Field(i).Set(reflect.ValueOf(i))
		case "int64":
			i64, _ := strconv.ParseInt(value, 10, 64)
			vv.Field(i).Set(reflect.ValueOf(i64))
		case "uint64":
			i64, _ := strconv.ParseUint(value, 10, 64)
			vv.Field(i).Set(reflect.ValueOf(i64))

		case "*xmux.FormFile":
			fmt.Println("file")
			f, h, err := r.FormFile(key)
			file := reflect.New(reflect.TypeOf(xmux.FormFile{}))
			fmt.Println(file.Elem().NumField())
			file.Elem().Field(0).Set(reflect.ValueOf(f))
			file.Elem().Field(1).Set(reflect.ValueOf(h))
			if err != nil {
				file.Elem().Field(2).Set(reflect.ValueOf(err))
			}

			vv.Field(i).Set(file)
		// case reflect.Ptr:
		// 	json.Unmarshal([]byte(value), vv.Field(i).Interface())
		// case reflect.Struct, reflect.Slice:
		// 	json.Unmarshal([]byte(value), vv.Field(i).Interface())
		default:
			log.Println("not format ", tt.Field(i).Type.Kind().String())
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Set("aaaa", "bbb")
	user := xmux.GetInstance(r).Data.(*User)
	fmt.Printf("%#v\n", *user)
	fmt.Println(user.File.Header.Filename)
}

type Global struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type User struct {
	UserName string         `json:"username" form:"username"`
	PassWord string         `json:"password" form:"password"`
	Gender   bool           `json:"form" form:"gender"`
	File     *xmux.FormFile `json:"file" form:"file"`
}

func main() {
	global := &Global{
		Code: 200,
	}

	router := xmux.NewRouter()
	group := xmux.NewGroupRoute().BindResponse(global)
	group.Post("/post", home)
	router.Get("/get", home)
	router.Any("/", home).AddModule(home1).Bind(User{})
	router.Run(":8888")
}
