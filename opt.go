package xmux

import (
	"fmt"
	"reflect"
	"strings"
)

// 自动生成接口文档
type Opt struct {
	Name        string
	Typ         string
	Need        string
	Default     string
	Information string
}

func PostOpt(s interface{}) (o Opt) {
	tpy := reflect.TypeOf(s)
	if tpy.Kind() == reflect.Ptr {
		tpy = tpy.Elem()
	}
	if tpy.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < tpy.NumField(); i++ {
		o.Default = tpy.Field(i).Tag.Get("default")
		o.Name = strings.Split(tpy.Field(i).Tag.Get("json"), ",")[0]
		o.Need = tpy.Field(i).Tag.Get("need")
		o.Typ = tpy.Field(i).Tag.Get("type")
		o.Information = tpy.Field(i).Tag.Get("information")
	}
	return
}

func GetOpt(s map[string]string) string {
	if len(s) == 0 {
		return ""
	}
	pms := make([]string, 0)
	for k, v := range s {
		pm := fmt.Sprintf("%s=%s", k, v)
		pms = append(pms, pm)
	}
	return fmt.Sprintf("?%s", strings.Join(pms, "&"))
}
