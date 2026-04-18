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

func (r *router) unmarshalJson(req *http.Request, fd *FlowData) (bool, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}

	req.Body.Close()
	fd.Body = b

	if len(b) > 0 {
		err = json.Unmarshal(b, &fd.Data)
	}

	return false, err
}

func (r *router) unmarshalYaml(req *http.Request, fd *FlowData) (bool, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}
	req.Body.Close()
	fd.Body = b
	if len(b) > 0 {
		err = yaml.Unmarshal(b, &fd.Data)
	}

	return false, err
}

func (r *router) unmarshalXml(req *http.Request, fd *FlowData) (bool, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, err
	}
	req.Body.Close()
	fd.Body = b
	if len(b) > 0 {
		err = xml.Unmarshal(b, &fd.Data)
	}

	return false, err
}

func (r *router) bind(route *rt, w http.ResponseWriter, req *http.Request, fd *FlowData) bool {
	// 数据绑定
	defer req.Body.Close()
	switch route.bindType {
	case jsonT:
		cont, err := r.unmarshalJson(req, fd)
		if err != nil {
			return r.UnmarshalError(err, w, req)
		}
		return cont
	case yamlT:
		cont, err := r.unmarshalYaml(req, fd)
		if err != nil {
			return r.UnmarshalError(err, w, req)
		}
		return cont
	case xmlT:
		cont, err := r.unmarshalXml(req, fd)
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

		for _, head := range headers {
			head = strings.Trim(head, " ")
			if head == MIMEJSON {
				cont, err := r.unmarshalJson(req, fd)
				if err != nil {
					return r.UnmarshalError(err, w, req)
				}
				return cont
			}
			if head == MIMEXML || head == MIMEXML2 {
				cont, err := r.unmarshalXml(req, fd)
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
				return cont
			}

		}

	}
	w.Write([]byte("unsupport content-type"))
	return true
}

func (r *router) unmarsharForm(w http.ResponseWriter, req *http.Request, fd *FlowData) (bool, error) {
	cl := req.Header.Get("Content-Length")
	length, err := strconv.Atoi(cl)
	if err == nil && length > 0 {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return true, err
		}
		if length > r.MaxPrintLength {
			fd.Body = b[:r.MaxPrintLength]
		} else {
			fd.Body = b
		}

		req.Body = io.NopCloser(bytes.NewBuffer(b))
	}
	tt := reflect.TypeOf(fd.Data).Elem()
	vv := reflect.ValueOf(fd.Data).Elem()
	l := tt.NumField()
	for i := 0; i < l; i++ {
		keys := tt.Field(i).Tag.Get("form")
		tagkeys := strings.Split(keys, ",")
		key := tagkeys[0]
		if key == "" {
			continue
		}
		value := req.FormValue(key)

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
		case reflect.Int, reflect.Int64, reflect.Int8, reflect.Int16, reflect.Int32:
			i64, _ := strconv.ParseInt(value, 10, 64)
			vv.Field(i).SetInt(i64)
		case reflect.Uint64, reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32:
			i64, _ := strconv.ParseUint(value, 10, 64)
			vv.Field(i).SetUint(i64)

		default:
			return false, errors.New("not support type, url: " + req.URL.Path)
		}
	}
	return false, nil
}
