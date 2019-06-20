package urlshortener

import "net/http"

func MapHandler(redirects map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if to, ok := redirects[r.URL.Path]; ok {
			http.Redirect(w, r, to, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}
