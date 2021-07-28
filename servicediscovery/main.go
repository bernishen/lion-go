package main

import "fmt"

type testinstance struct {
	a string
	b string
}

func main() {
	var t *testinstance
	t = new(testinstance)
	t.a = "aaa"
	t.b = "bbb"
	fmt.Printf("%v\n", &t.a)
	fmt.Printf("%v\n", t)
	var a, b *string
	a = new(string)
	b = new(string)
	*a = "aaa"
	*b = "bbb"
	fmt.Printf("%v\n", *a)
	fmt.Printf("%v\n", b)
}
