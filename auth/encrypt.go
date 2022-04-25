package auth

import (
	"crypto/md5"
	"fmt"
)

func Md5(base []byte) string {
	md5.New()
	return fmt.Sprintf("%x", md5.Sum(base))
}
