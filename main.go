package main

import (
	"fmt"
	"net/http"
	middleware "short/middlewares"
	"short/server"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func main() {
	mainSrv := server.NewServer("0.0.0.0", 5000)

	middlewareStack := middleware.CreateStack(
		middleware.RecoveryMiddleware,
		middleware.BasicLogger,
	)

	mainSrv.Use(middlewareStack)

	mainSrv.Register("/", handler)

	mainSrv.Run()

}
