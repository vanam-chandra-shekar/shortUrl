package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
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
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
