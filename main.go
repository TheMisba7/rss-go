package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"rss-aggregator/handlers"
	"rss-aggregator/internal/database"
	"rss-aggregator/scraper"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	dbAddr := os.Getenv("DATASOURCE")
	db, err := sql.Open("postgres", dbAddr)
	if db == nil {
		panic("db connection cannot be nil")
	}
	config := handlers.ApiConfig{DB: database.New(db)}
	defer db.Close()

	mainRouter := chi.NewRouter()
	mainRouter.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))

	mainRouter.Mount("/v1", config.RegisterRoutes())
	server := http.Server{
		Handler: mainRouter,
		Addr:    fmt.Sprintf(":%v", port),
	}

	//start scraper in separate goroutine
	go scraper.StartScrapper(2, 10, config.DB)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
