package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/mnhkahn/gg/conf"
)

func main() {
	app := cli.NewApp()
	app.Name = "gg"
	app.Usage = "A Deploy tool written in Golang. It will works with Supervisor."
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "start",
			Usage:   "Rebuild & run in local path.",
			Aliases: []string{"s"},
			Action: func(c *cli.Context) {
				// Build
				Build()
				// Run
				Start()
			},
		},
		{
			Name:    "build",
			Usage:   "Build.",
			Aliases: []string{"b"},
			Action: func(c *cli.Context) {
				// Build
				Build()
			},
		},
		{
			Name:    "deploy",
			Usage:   "Build & restart.",
			Aliases: []string{"d"},
			Action: func(c *cli.Context) {
				// Build()
				if conf.NewGGConfig().IsGitPull {
					GitPull()
				}
				Backup()
				Deploy()
				if conf.NewGGConfig().IsNgrok {
					Ngrok()
				}
			},
		},
		{
			Name:    "pack",
			Usage:   "Pack & generate supervisor configuration file",
			Aliases: []string{"p"},
			Action: func(c *cli.Context) {
				Build()
				// Supervisor
				Supervisor()

				// Pack
				if err := Pack(); err != nil {
					log.Println("Generate package error", err)
				} else {
					log.Printf("Pack success in %s.\n", conf.NewGGConfig().AppPath)
				}
			},
		},
	}

	app.Run(os.Args)
}
