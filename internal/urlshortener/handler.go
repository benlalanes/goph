package urlshortener

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// YAMLHandler returns an HTTP handler that shortens URLs by redirecting
// requests according to the specified YAML configuration.
func YAMLHandler(y []byte, fallback http.Handler) (http.HandlerFunc, error) {
	m := make(map[string]string)
	if err := yaml.Unmarshal(y, &m); err != nil {
		return nil, err
	}

	return MapHandler(m, fallback), nil
}

// JSONHandler returns an HTTP handler that shortens URLs by redirecting
// requests according to the specified JSON configuration.
func JSONHandler(j []byte, fallback http.Handler) (http.HandlerFunc, error) {
	m := make(map[string]string)
	if err := json.Unmarshal(j, &m); err != nil {
		return nil, err
	}

	return MapHandler(m, fallback), nil
}

func MapHandler(redirects map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if to, ok := redirects[r.URL.Path]; ok {
			http.Redirect(w, r, to, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}
