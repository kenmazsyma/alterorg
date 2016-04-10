package main

import "fmt"

func main() {
	s := "12345"
	b := []byte(s)
	fmt.Printf("%02x", b)
}
