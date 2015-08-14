package main

import (
	"os"
	"os/exec"

	"github.com/mnhkahn/gg/conf"
)

func Start() {
	cmd := exec.Command("./" + conf.NewGGConfig().AppName + conf.NewGGConfig().AppSuffix)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
