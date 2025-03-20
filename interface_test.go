package xmux

import (
	"fmt"
	"net/http"
	"testing"
)

type Dog struct {
}

type Res struct {
	Id int64
}

type Animal interface {
	Dog
	Mofa()
}

func (dog *Dog) Get() *Res {
	res := &Res{}
	res.Id = 1
	return res
}

func (dog *Dog) GetContent(w http.ResponseWriter, r *http.Request) {
	// GetInstance(r).Data.Mofa()
}

func (*Router) BindStruct() {
	//
}

func TestInherit(t *testing.T) {
	router := NewRouter()
	router.BindStruct()
	router.Run()
}

type Cat interface {
	On()
}

type JiaFeiMao struct {
	Name   string
	Weight float64
}

func (jfm *JiaFeiMao) On() {
	fmt.Println("name is jiafeimao, Weight is 105")
}

type HomeMao struct {
	Name   string
	Weight float64
}

func (jfm *HomeMao) On() {
	fmt.Println("name is home cat, Weight is 15")
}

type TT interface {
	On()
}

func Online(x TT) {
	x.On()
}
