package middleware

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"short/templ"
	"strings"
	"time"
)

type MiddleWare func(http.Handler) http.Handler

func CreateStack(ms ...MiddleWare) MiddleWare {

	return func(next http.Handler) http.Handler {

		for i := len(ms) - 1; i >= 0; i-- {
			x := ms[i]
			next = x(next)
		}
		return next
	}

}

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type wrapedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func BasicLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wraped := &wrapedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(w, r)
		log.Println(wraped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				mssg := "Caught painc: %v, Stack trace: %s"
				log.Printf(mssg, err, string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				templ.PageInternalServerError.Execute(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

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
