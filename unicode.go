package xmux

import (
	"strconv"
	"strings"
)

// 将前端传进来的部分中文被编译成unicode编码进行还原
func UnescapeUnicode(raw []byte) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}
