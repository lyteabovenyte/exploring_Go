package main

import (
	"fmt"
	"os/exec"
)

func main() {
	const (
		// just the binary's name.
		kubectl = "kubectl"
		git     = "git"
	)

	_, err1 := exec.LookPath(kubectl)
	if err1 != nil {
		fmt.Println("cannot find kubectl in our PATH")
	}

	_, err2 := exec.LookPath(git)
	if err2 != nil {
		fmt.Println("cannot find git in our PATH")
	}
	fmt.Println("Done searching.")
}
