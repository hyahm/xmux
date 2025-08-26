package xmux

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"
)

const (
	zero  = byte('0')
	one   = byte('1')
	lsb   = byte('[') // left square brackets
	rsb   = byte(']') // right square brackets
	space = byte(' ')
)

var uint8arr [8]uint8

// ErrBadStringFormat represents a error of input string's format is illegal .
var ErrBadStringFormat = errors.New("bad string format")

// ErrEmptyString represents a error of empty input string.
var ErrEmptyString = errors.New("empty string")

var ErrTypeUnsupport = errors.New("data type is unsupported")

var enableJsonV2 bool

func getEnableJsonV2() {
	cmd := exec.Command("go", "env")

	// 创建一个 bytes.Buffer 来存储命令的输出
	var out bytes.Buffer
	cmd.Stdout = &out

	// 执行命令
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		value := scanner.Text()
		if strings.Contains(scanner.Text(), "GOEXPERIMENT") {

			v := strings.Split(value, "=")[1]
			if v == "jsonv2" {
				enableJsonV2 = true
			}
			break
		}
	}
	// 打印
	// 命令的输出

}

func init() {
	getEnableJsonV2()
	// for _, value := range os.Environ() {
	// 	if strings.Contains(value, "GOEXPERIMENT") {
	// 		fmt.Println(value)
	// 	}

	// }

	// jsonv2 "encoding/json/v2"
	uint8arr[0] = 128
	uint8arr[1] = 64
	uint8arr[2] = 32
	uint8arr[3] = 16
	uint8arr[4] = 8
	uint8arr[5] = 4
	uint8arr[6] = 2
	uint8arr[7] = 1
}

// 保存url里面的参数
type params map[string]string // url 参数对应的值

var allparams map[string]params // 保存的url 参数
var paramsLocker sync.RWMutex

func init() {
	allparams = make(map[string]params)
	paramsLocker = sync.RWMutex{}
}
func Var(r *http.Request) params {
	return getParams(r.URL.Path)
}

func getParams(key string) params {
	paramsLocker.RLock()
	defer paramsLocker.RUnlock()
	return allparams[key]
}

func setParams(key string, params params) {
	paramsLocker.Lock()
	allparams[key] = params
	paramsLocker.Unlock()
}
