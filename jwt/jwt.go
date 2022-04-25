package xmux

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var ErrToken = errors.New("jwt token error")
var ErrTokenExpired = errors.New("jwt token expired")

// 使用jwt 认证
const header = `{'typ': 'JWT', 
'alg': 'HS256'
}`

// type Jwter interface {
// 	// type Token struct {
// 	// 	Id       int64  `json:"id"`
// 	// 	NickName string `json:"nickname"`
// 	// 	Roles    string `json:"roles"`
// 	// 	UserName string `json:"username"`
// 	// 	Avatar   string `json:"avatar"`
// 	// 	Exp      int64  `json:"exp"`
// 	// }
// 	// func (tk *Token) Marshal() []byte {
// 	// 	payload, err := json.Marshal(tk)
// 	// 	if err != nil {
// 	// 		return nil
// 	// 	}
// 	// func (tk *Token) Unmarshal(b []byte) error {
// 	// 	return json.Unmarshal(b, tk)
// 	// }

// 	// Marshal() []byte
// 	// Unmarshal([]byte) error
// }

// 创建jwt
func MakeJwt(salt string, tk interface{}) (string, error) {
	payload, err := json.Marshal(tk)
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString([]byte(header))
	p := base64.StdEncoding.EncodeToString([]byte(payload))
	pre := s + "." + p
	token := pre + "." + getHc(pre, salt)
	return token, nil
}

// 获取hc
func getHc(b, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	io.WriteString(h, b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 检查jwt, must be a point
func GetJwt(jwt, salt string, token interface{}) error {
	if reflect.TypeOf(token).Kind() != reflect.Ptr {
		return errors.New("token must be a pointer")
	}
	js := strings.Split(jwt, ".")
	if len(js) < 3 {
		return errors.New("invalid jwt")
	}
	b, err := base64.StdEncoding.DecodeString(js[1])
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, token)
	// err = json.NewDecoder(bytes.NewReader(b)).Decode(token)
	if err != nil {
		return err
	}

	// // 检查过期时间
	pre := js[0] + "." + js[1]
	if getHc(pre, salt) == js[2] {
		return nil
	}
	return ErrToken
}
