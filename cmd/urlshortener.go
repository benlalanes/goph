package cmd

import (
	"flag"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/http"

	"github.com/benlalanes/goph/internal/urlshortener"
)

var URLShortener = cli.Command{
	Name:   "urlshortener",
	Usage:  "Run a server to shorten URLs",
	Action: runURLShortener,
}

func runURLShortener(ctx *cli.Context) error {

	flags := flag.NewFlagSet("urlshortener", flag.PanicOnError)

	filepath := flags.String("filepath", "", "path to the file holding redirect configuration")

	err := flags.Parse(ctx.Args())
	if err != nil {
		return err
	}

	fmt.Println(*filepath)

	m := map[string]string{"/google": "https://google.com"}
	fb := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Redirect for specified path was not found.")
	}

	log.Fatal(http.ListenAndServe("localhost:8080", urlshortener.MapHandler(m,
		http.HandlerFunc(fb))))

	return nil

}
