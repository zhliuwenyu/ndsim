package main

import "fmt"

func main() {
	a := []byte("abcdefg")
	b := a[2:4]
	fmt.Println(a)
	fmt.Println(b)
	b = append(b, '3')
	fmt.Println(a)
	fmt.Println(b)
}
