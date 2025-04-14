package middleware

import (
	"net/http"
)

type MiddleWare func(http.Handler) http.Handler
