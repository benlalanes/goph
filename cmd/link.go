package cmd

import (
	"errors"
	"fmt"
	"github.com/benlalanes/goph/internal/link"
	"github.com/urfave/cli"
	"golang.org/x/net/html"
	"os"
)

var Link = cli.Command{
	Name:   "link",
	Action: runLink,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "HTML file to parse for links",
			Value: "",
		},
	},
}

func runLink(ctx *cli.Context) error {

	fp := ctx.String("file")
	if fp == "" {
		return errors.New("file argument must be specified")
	}

	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	defer f.Close()

	root, err := html.Parse(f)
	if err != nil {
		return err
	}

	links := link.ParseLinks(root)

	fmt.Printf("Found %d links.\n", len(links))

	for _, lnk := range links {
		fmt.Printf("<a href='%s'>%s</a>\n", lnk.Href, lnk.Text)
	}

	return nil

}
