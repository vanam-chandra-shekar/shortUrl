package middleware

import "net/http"

func CreateStack(ms ...MiddleWare) MiddleWare {

	return func(next http.Handler) http.Handler {

		for i := len(ms) - 1; i >= 0; i-- {
			x := ms[i]
			next = x(next)
		}
		return next
	}

}
