package main

import (
	"errors"
	"log"
	"os/exec"
)

func stdoutPipe() error {
	c1 := exec.Command("echo", "foo")
	c2 := exec.Command("doesnotexist")
	pipe, err := c1.StdoutPipe()
	if err != nil {
		return err
	}
	c2.Stdin = pipe
	err = c1.Start()
	if err != nil {
		return err
	}
	err = c2.Start()
	if err == nil {
		return errors.New("Expected error from nonexisting command")
	}
	err = c1.Wait()
	if err == nil {
		return errors.New("Expected error writing to nonexisting command")
	}
	return nil
}

func stdinPipe() error {
	c1 := exec.Command("echo", "foo")
	c2 := exec.Command("doesnotexist")
	pipe, err := c2.StdinPipe()
	if err != nil {
		return err
	}
	c1.Stdout = pipe
	err = c1.Start()
	if err != nil {
		return err
	}
	err = c2.Start()
	if err == nil {
		return errors.New("Expected error from nonexisting command")
	}
	err = c1.Wait()
	if err == nil {
		// what about buffering, hmm?
		return errors.New("Expected error writing to nonexisting command")
	}
	return nil
}

func main() {
	var err error
	// on my machine (debian VM), err is consistently "expected error writing"
	err = stdoutPipe()
	if err != nil {
		log.Print("stdout pipe: ", err)
	} else {
		log.Print("stdout piping behaved as expected")
	}
	// on my machine, this is 50/50: err = nil, or the same as above
	err = stdinPipe()
	if err != nil {
		log.Print("stdin pipe: ", err)
	} else {
		log.Print("stdin piping behaved as expected")
	}
	// so, buffering.
}
