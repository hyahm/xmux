package xmux

import "mime/multipart"

type bindType int

const (
	jsonT bindType = 1
	yamlT bindType = 2
	xmlT  bindType = 3
	formT bindType = 4
	fileT bindType = 5
)

type FormFile struct {
	File   multipart.File
	Header *multipart.FileHeader
	Err    error
}

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
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

// type Binding interface {
// 	Name() string
// 	Bind(*http.Request, any) error
// }

// func Default(method, contentType string) Binding {
// 	if method == "GET" {
// 		return Form
// 	}

// 	switch contentType {
// 	case MIMEJSON:
// 		return JSON
// 	case MIMEXML, MIMEXML2:
// 		return XML
// 	case MIMEPROTOBUF:
// 		return ProtoBuf
// 	case MIMEYAML:
// 		return YAML
// 	case MIMEMultipartPOSTForm:
// 		return FormMultipart
// 	default: // case MIMEPOSTForm:
// 		return Form
// 	}
// }
