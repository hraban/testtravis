package main

import (
	"log"
	"os/exec"
)

func main() {
	c1 := exec.Command("echo", "foo")
	c2 := exec.Command("doesnotexist")
	pipe, err := c1.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	c2.Stdin = pipe
	err = c1.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = c2.Start()
	if err == nil {
		log.Fatal("Expected error from nonexisting command")
	}
	err = c1.Wait()
	if err == nil {
		log.Fatal("Expected error writing to nonexisting command")
	}
}
