package main

import "fmt"

func main() {
	x := make(map[string]map[string]string)
	y := make(map[string]string)
	y["2"] = "3"
	x["1"] = y
	fmt.Println("y")
	//s := "/name/asdf/a787"
	//g := regexp.MustCompile("^/name/(\\w+)/(\\w+)$")
	//x := g.FindStringSubmatch(s)
	//for _, v := range x {
	//	fmt.Println(v)
	//}

}
