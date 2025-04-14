package middleware

import (
	"log"
	"net/http"
	"time"
)

// ANSI color codes
// const (
// 	colorReset  = "\033[0m"
// 	colorRed    = "\033[31m"
// 	colorGreen  = "\033[32m"
// 	colorYellow = "\033[33m"
// 	colorBlue   = "\033[34m"
// 	colorWhite  = "\033[37m"
// 	colorPurple = "\033[35m"
// 	colorCyan   = "\033[36m"
// )

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
