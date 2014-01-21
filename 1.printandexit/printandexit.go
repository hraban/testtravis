package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	_, err := fmt.Println(strings.Join(os.Args[1:], " "))
	var status int
	if err == nil {
		status = 0
	} else {
		status = 1
	}
	fmt.Fprintln(os.Stderr, status)
	os.Exit(status)
}
