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

func PostOpt(s interface{}) []Opt {
	opts := make([]Opt, 0)
	tpy := reflect.TypeOf(s)
	if tpy.Kind() == reflect.Ptr {
		tpy = tpy.Elem()
	}
	if tpy.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < tpy.NumField(); i++ {
		opt := Opt{}
		opt.Default = tpy.Field(i).Tag.Get("default")
		opt.Name = strings.Split(tpy.Field(i).Tag.Get("json"), ",")[0]
		opt.Need = tpy.Field(i).Tag.Get("need")
		opt.Typ = tpy.Field(i).Tag.Get("type")
		opt.Information = tpy.Field(i).Tag.Get("information")
		opts = append(opts, opt)
	}
	return opts
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
