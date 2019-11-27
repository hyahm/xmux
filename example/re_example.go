package main

import (
	"fmt"
	"regexp"
)

func main() {
	s := "/name/asdf/a787"
	g := regexp.MustCompile("^/name/(\\w+)/(\\w+)$")
	x := g.FindStringSubmatch(s)
	for _, v := range x {
		fmt.Println(v)
	}

}
