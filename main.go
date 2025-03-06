package main

import (
	"context"
	"log"
	"os"
	"short/handlers"
	middleware "short/middlewares"
	"short/server"
	"short/templ"

	"github.com/jackc/pgx/v5"
)

func GetEnvOrDefault(name string, def string) string {
	if val, con := os.LookupEnv(name); con {
		return val
	}

	return def
}

func main() {

	dbstring := GetEnvOrDefault("DBSTRING", "")

	port := GetEnvOrDefault("PORT", "5000")

	conn, err := pgx.Connect(context.Background(), dbstring)

	if err != nil {
		log.Fatalf("[Error] %v\n" , err)
	}
	log.Println("[Sucess] : DataBase connected")
	defer conn.Close(context.Background())

	mainSrv := server.NewServer("0.0.0.0", port)
	myHandler := handlers.NewHandler(templ.NewTemplBlob("./web/*.html"), conn)

	middlewareStack := middleware.CreateStack(
		middleware.RecoveryMiddleware,
		middleware.GlobalRateLimiter(10, 15),
		middleware.BasicLogger,
		middleware.StaticFileMiddleware("./web/css", "/css/"),
	)

	mainSrv.Use(middlewareStack)

	mainSrv.Register("/", myHandler.RootHandler)

	mainSrv.Register("POST /onUrlFormSubmit", myHandler.HxOnUrlFormSubmit)

	mainSrv.Run()

}
