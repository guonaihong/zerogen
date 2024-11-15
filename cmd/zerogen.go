package main

import (
	"fmt"

	"github.com/guonaihong/clop"
	"github.com/guonaihong/zerogen"
)

func main() {
	var zeroGen zerogen.ZeroGen
	clop.Bind(&zeroGen)
	err := zeroGen.Run()
	if err != nil {
		fmt.Println(err)
	}
}
