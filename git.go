package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mnhkahn/gg/conf"
)

func GitPull() {
	args := append([]string{"checkout", conf.NewGGConfig().GitPullBranch})
	cmd := exec.Command("git", args...)
	log.Println(strings.Join(cmd.Args, " "))
	var err_output bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &err_output

	if err := cmd.Start(); err != nil { //Use start, not run
		log.Println("An error occured: ", err) //replace with logger, or anything you want
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Git checkout error: %v. %s.\n", err, string(err_output.Bytes()))
		return
	}

	args = append([]string{"pull", "origin", conf.NewGGConfig().GitPullBranch})
	cmd = exec.Command("git", args...)
	log.Println(strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = &err_output

	if err := cmd.Start(); err != nil { //Use start, not run
		log.Println("An error occured: ", err) //replace with logger, or anything you want
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Git pull error: %v. %s.\n", err, string(err_output.Bytes()))
		return
	}
	log.Println("Git pull Success.")
}
