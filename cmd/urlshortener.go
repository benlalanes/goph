package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	if *filepath == "" {
		return errors.New("-filepath argument must be specified")
	}

	f, err := os.Open(*filepath)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	fb := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintln(w, "Redirect for specified path was not found.")
	}

	handler, err := urlshortener.YAMLHandler(b, http.HandlerFunc(fb))

	log.Fatal(http.ListenAndServe("localhost:8080", handler))

	return nil

}
