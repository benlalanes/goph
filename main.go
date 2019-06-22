package main

import (
	"github.com/urfave/cli"
	"log"
	"os"

	"github.com/benlalanes/goph/cmd"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cmd.URLShortener,
		cmd.Link,
		cmd.Task,
		cmd.RecoverChromaDemo,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
