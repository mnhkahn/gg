package main

import (
	"os"
	"os/exec"
)

func Start() {
	cmd := exec.Command("./" + NewGGConfig().AppName + NewGGConfig().AppSuffix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
