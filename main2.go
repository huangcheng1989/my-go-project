package main

import (
	"fmt"
	"strings"
)

func main() {
	var a string
	a = "image/png" // 随便写的例子，因为字符串变量中的单双引号是我们不能提前知道的
	b := strings.Index(a, "/")
	c := a[0:strings.Index(a, "/")]
	fmt.Println(a, b, c)
}
