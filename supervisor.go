package main

import (
	"log"
	"os"
	"text/template"

	"github.com/mnhkahn/gg/conf"
)

func Supervisor() {
	f, err := os.OpenFile(conf.NewGGConfig().SupervisorConf, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("Generate supervisor file error", err)
	}

	tmpl, _ := template.New("supervisor").Parse(SupervisorTemplate)
	tmpl.Execute(f, conf.NewGGConfig())
}

var (
	SupervisorTemplate = `[program:{{.AppName}}]
directory = {{.RunDirectory}}
command = {{.RunDirectory}}{{.AppName}}
autostart = true
autorestart = true
startsecs = 5
user = {{.RunUser}}
redirect_stderr = true
stdout_logfile = {{.LogDirectory}}{{.AppName}}.log
`
)
