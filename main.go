package main

import (
	"database/sql"
	"github.com/japangermany1998/go-server/controller"
	"github.com/japangermany1998/go-server/internal/database"
	"github.com/japangermany1998/go-server/router"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	controller.ApiCfg.NewDB(dbQueries)
	controller.ApiCfg.NewPlatform(os.Getenv("PLATFORM"))
	controller.ApiCfg.NewJWTSecret(os.Getenv("JWT_SECRET"))
	controller.ApiCfg.NewPolkaKey(os.Getenv("POLKA_KEY"))

	router.Router(mux)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
