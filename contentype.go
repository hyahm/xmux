package xmux

type ContentType string

func (ct ContentType) String() string {
	return string(ct)
}

const (
	URLENCODED ContentType = "application/x-www-form-urlencoded"
	JSON       ContentType = "application/json"
	FORM       ContentType = "multipart/form-data"
	XML        ContentType = "application/xml"
)
