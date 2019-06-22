package cmd

import (
	"errors"
	"fmt"
	"github.com/alecthomas/chroma/quick"
	recover "github.com/benlalanes/goph/internal/recover-chroma"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strings"
)

const DebugPrefix = "/debug/"

var RecoverChromaDemo = cli.Command{
	Name: "recover-chroma-demo",
	Action: runRecoverChromaDemo,
}

func panickingHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "About to panic.")

	panic(errors.New("panicking now"))
}

func runRecoverChromaDemo(ctx *cli.Context) error {

	log.Println("Running recover-chroma-demo.")

	root := ctx.Args().First()

	if root == "" {
		root = "."
	}

	root, err := filepath.Abs(root)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Serving files from %s.\n", root)

	rootStaticServer := recover.StaticChromaServer(root)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hello, world!")
	})

	mux.HandleFunc(DebugPrefix, func(w http.ResponseWriter, r *http.Request) {

		//if strings.Index(r.URL.Path, DebugPrefix) == -1 {
		//	w.WriteHeader(http.StatusBadRequest)
		//	_, _ = fmt.Fprint(w, "Incorrect URL format.")
		//	return
		//}

		fp := strings.TrimPrefix(r.URL.Path, DebugPrefix)

		rootStaticServer.GetFilepathServer(fp).ServeHTTP(w, r)
	})

	mux.Handle("/panic", recover.NewRecoverable(http.HandlerFunc(panickingHandler), func(w http.ResponseWriter, r *http.Request) {

		stack := debug.Stack()

		lines := strings.Split(string(stack), "\t")

		for _, line := range lines {
			parts := strings.Split(line, " ")

			lastColon := strings.LastIndex(parts[0], ":")

			if lastColon != -1 {
				absPath, _ := parts[0][:lastColon], parts[0][lastColon+1:]
				writeSourceFile(w, absPath)
				break
			}


		}

		fmt.Fprintln(w, "done")

		//
		//_, _ = w.Write(stack)

	}))

	return http.ListenAndServe("localhost:8080", mux)
}

func writeSourceFile(w io.Writer, fp string) {

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Println(err)
	}

	_ = quick.Highlight(w, string(b), "go", "html", "monokai")

}
