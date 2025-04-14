package middleware

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"short/templ"
	"strings"
)

func StaticFileMiddleware(dir string, endpoint string) MiddleWare {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !strings.HasPrefix(r.URL.Path, endpoint) {
				next.ServeHTTP(w, r)
				return
			}

			reqpath := strings.TrimPrefix(r.URL.Path, endpoint)
			cleanPath := filepath.Clean(reqpath)

			fullpath := filepath.Join(dir, cleanPath)

			absDir, _ := filepath.Abs(dir)
			absPath, _ := filepath.Abs(fullpath)

			if !strings.HasPrefix(absPath, absDir) {
				log.Println("Blocked path traversal attempt:", absPath)
				w.WriteHeader(http.StatusForbidden)
				templ.PageForbidden.Execute(w, nil)
				return
			}

			if _, err := os.Stat(fullpath); os.IsNotExist(err) {
				w.WriteHeader(http.StatusNotFound)
				templ.Page404.Execute(w, nil)
				return
			}

			http.ServeFile(w, r, fullpath)

		})
	}

}
