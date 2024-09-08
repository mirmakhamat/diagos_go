package main

import (
	"log"
	"os"

	"github.com/mirmakhamat/diagos_go/services"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "DiagOS",
		Usage: "A CLI for DiagOS, a tool for diagnosing your system.",
		Commands: []*cli.Command{
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "options for task status",
				Action:  services.Status,
			},
			{
				Name:    "cpu",
				Aliases: []string{"c"},
				Usage:   "options for task cpu",
				Action:  services.Cpu,
			},
			{
				Name:    "find",
				Aliases: []string{"f"},
				Usage:   "options for task find",
				Action:  services.Find,
			},
			{
				Name:    "internet",
				Aliases: []string{"i"},
				Usage:   "options for task internet",
				Action:  services.Internet,
			},
			{
				Name:    "memory",
				Aliases: []string{"m"},
				Usage:   "options for task memory",
				Action:  services.Memory,
			},
			{
				Name:    "storage",
				Aliases: []string{"st"},
				Usage:   "options for task storage",
				Action:  services.Storage,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
