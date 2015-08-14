package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mnhkahn/gg/conf"
)

func Build() {
	cmd := exec.Command("go", "build", "-o", conf.NewGGConfig().AppName+conf.NewGGConfig().AppSuffix, conf.NewGGConfig().MainApplication)
	log.Println(strings.Join(cmd.Args, " "))
	var err_output bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &err_output

	if err := cmd.Start(); err != nil { //Use start, not run
		log.Println("An error occured: ", err) //replace with logger, or anything you want
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Build error: %v. %s.\n", err, string(err_output.Bytes()))
		return
	}
	log.Println("Build Success.")
}
