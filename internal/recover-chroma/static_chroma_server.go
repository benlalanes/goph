package recover_chroma

import (
	"errors"
	"fmt"
	"io/ioutil"
	//"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

func GetFilepathServer(fp string) (http.HandlerFunc, error) {

	if !filepath.IsAbs(fp) {
		return nil, errors.New("filepath must be absolute")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		b, err := ioutil.ReadFile(fp)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, "File %s not found.", fp)
			return
		}

		_ = quick.Highlight(w, string(b), "go", "html", "monokai")

	}), nil



}


type StaticChromaServer string

const DebugPrefix = "/debug/"

func (s StaticChromaServer) GetFilepathServer(fp string) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile(filepath.Join(string(s), fp))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, "File %s not found.", filepath.Join(string(s), fp))
			return
		}

		_ = quick.Highlight(w, string(b), "go", "html", "monokai")
	})
}

func (s StaticChromaServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("serving")

	if strings.Index(r.URL.Path, DebugPrefix) == -1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Incorrect URL format.")
		return
	}

	fp := strings.TrimPrefix(r.URL.Path, DebugPrefix)

	b, err := ioutil.ReadFile(filepath.Join(string(s), fp))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "File %s not found.", filepath.Join(string(s), fp))
		return
	}

	_ = quick.Highlight(w, string(b), "go", "html", "monokai")
}