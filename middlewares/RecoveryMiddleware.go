package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"short/templ"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				mssg := "Caught painc: %v, \nStack trace: \n%s\n"
				log.Printf(mssg, err, string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				templ.PageInternalServerError.Execute(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
