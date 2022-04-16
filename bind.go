package xmux

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	yaml "gopkg.in/yaml.v2"
)

type bindType int

const (
	jsonT bindType = 1
	yamlT bindType = 2
	xmlT  bindType = 3
	formT bindType = 4
	headT bindType = 5
)

const (
	MIMEJSON              = "application/json"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

func (r *Router) unmarsharJson(w http.ResponseWriter, req *http.Request, fd *FlowData) (bool, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}
	if r.PrintRequestStr {
		r.RequestBytes(b, req)
	}
	tt := reflect.TypeOf(fd.Data).Elem()
	l := tt.NumField()
	for i := 0; i < l; i++ {
		keys := tt.Field(i).Tag.Get("json")
		tagkeys := strings.Split(keys, ",")
		if tagkeys[0] == "" {
			continue
		}
		key := tagkeys[0]
		if len(tagkeys) > 1 {
			if strings.Contains(keys[len(key):], "require") && !gjson.Get(string(b), key).Exists() {
				if r.NotFoundRequireField(key, w, req) {
					return true, nil
				}
			}
		}
	}

	err = json.Unmarshal(b, &fd.Data)
	return false, err
}

func (r *Router) unmarsharYaml(w http.ResponseWriter, req *http.Request, fd *FlowData) (bool, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}
	if r.PrintRequestStr {
		r.RequestBytes(b, req)
	}
	err = yaml.Unmarshal(b, &fd.Data)
	return false, err
}

func (r *Router) unmarsharXml(w http.ResponseWriter, req *http.Request, fd *FlowData) (bool, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}
	if r.PrintRequestStr {
		r.RequestBytes(b, req)
	}
	err = xml.Unmarshal(b, &fd.Data)
	return false, err
}

func (r *Router) bind(route *rt, w http.ResponseWriter, req *http.Request, fd *FlowData) bool {
	// 数据绑定
	defer req.Body.Close()
	switch route.bindType {
	case jsonT:
		cont, err := r.unmarsharJson(w, req, fd)
		if err != nil {
			return r.UnmarshalError(err, w, req)
		}
		return cont
	case yamlT:
		cont, err := r.unmarsharYaml(w, req, fd)
		if err != nil {
			return r.UnmarshalError(err, w, req)
		}
		return cont
	case xmlT:
		cont, err := r.unmarsharXml(w, req, fd)
		if err != nil {
			return r.UnmarshalError(err, w, req)
		}
		return cont
	case formT:
		cont, err := r.unmarsharForm(w, req, fd)
		if err != nil {
			return r.UnmarshalError(err, w, req)
		}
		return cont
	case headT:
		// 根据请求头自动解析
		ct := req.Header.Get("content-type")
		headers := strings.Split(ct, ";")
		if len(headers) == 1 && headers[0] == "" {
			cont, err := r.unmarsharJson(w, req, fd)
			if err != nil {
				return r.UnmarshalError(err, w, req)
			}
			return cont
		}
		for _, head := range headers {
			if head == MIMEJSON {
				cont, err := r.unmarsharJson(w, req, fd)
				if err != nil {
					return r.UnmarshalError(err, w, req)
				}
				if cont {
					return true
				}
			}
			if head == MIMEXML || head == MIMEXML2 {
				cont, err := r.unmarsharXml(w, req, fd)
				if err != nil {
					return r.UnmarshalError(err, w, req)
				}
				return cont

			}
			if head == MIMEPOSTForm || head == MIMEMultipartPOSTForm {
				cont, err := r.unmarsharForm(w, req, fd)
				if err != nil {
					return r.UnmarshalError(err, w, req)
				}
				if cont {
					return true
				}
			}

		}

	}
	return false
}

var MaxPrintLength uint64 = 2 << 10

func (r *Router) unmarsharForm(w http.ResponseWriter, req *http.Request, fd *FlowData) (bool, error) {
	cl := req.Header.Get("Content-Length")
	length, err := strconv.ParseUint(cl, 10, 64)
	if r.PrintRequestStr && err != nil && length >= MaxPrintLength {
		b, _ := io.ReadAll(req.Body)
		if r.PrintRequestStr {
			r.RequestBytes(b, req)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(b))
	}

	tt := reflect.TypeOf(fd.Data).Elem()
	vv := reflect.ValueOf(fd.Data).Elem()
	l := tt.NumField()
	for i := 0; i < l; i++ {
		keys := tt.Field(i).Tag.Get("form")
		tagkeys := strings.Split(keys, ",")
		if len(tagkeys) == 0 && tagkeys[0] == "" {
			continue
		}
		key := tagkeys[0]
		value := req.FormValue(key)
		if len(tagkeys) > 1 {
			if (strings.Contains(keys[len(key):], "require")) && value == "" {
				if r.NotFoundRequireField(key, w, req) {
					return true, nil
				}
			}
		}

		switch tt.Field(i).Type.Kind() {

		case reflect.String:
			vv.Field(i).Set(reflect.ValueOf(value))
		case reflect.Float32:
			b32, _ := strconv.ParseFloat(value, 32)
			vv.Field(i).Set(reflect.ValueOf(b32))
		case reflect.Float64:
			b64, _ := strconv.ParseFloat(value, 64)
			vv.Field(i).Set(reflect.ValueOf(b64))
		case reflect.Bool:
			ok, _ := strconv.ParseBool(value)
			vv.Field(i).Set(reflect.ValueOf(ok))
		case reflect.Int:
			i, _ := strconv.Atoi(value)
			vv.Field(i).Set(reflect.ValueOf(i))
		case reflect.Int64:
			i64, _ := strconv.ParseInt(value, 10, 64)
			vv.Field(i).Set(reflect.ValueOf(i64))
		case reflect.Uint64:
			i64, _ := strconv.ParseUint(value, 10, 64)
			vv.Field(i).Set(reflect.ValueOf(i64))

		default:
			return false, errors.New("not support type")
		}
	}
	return false, nil
}
