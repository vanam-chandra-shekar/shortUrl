package main

import (
	"short/handlers"
	middleware "short/middlewares"
	"short/server"
	"short/templ"
)

func main() {

	mainSrv := server.NewServer("0.0.0.0", 5000)
	myHandler := handlers.NewHandler(templ.NewTemplBlob("./web/*html"))

	middlewareStack := middleware.CreateStack(
		middleware.RecoveryMiddleware,
		middleware.BasicLogger,
		middleware.StaticFileMiddleware("./web/css", "/css/"),
	)

	mainSrv.Use(middlewareStack)

	mainSrv.Register("/", myHandler.RootHandler)

	mainSrv.Register("POST /onUrlFormSubmit", myHandler.HxOnUrlFormSubmit)

	mainSrv.Run()

}
