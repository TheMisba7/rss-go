package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"rss-aggregator/internal/database"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

type readiness struct {
	Status string
}

type error struct {
	Error string
}

type user struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (config apiConfig) newUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		user := user{
			Id:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		}
		decoder := json.NewDecoder(request.Body)
		decoder.Decode(&user)
		userParams := database.CreateUserParams{Name: user.Name, ID: user.Id, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}

		createUser, err := config.DB.CreateUser(request.Context(), userParams)
		if err != nil {
			fmt.Println("my err")
			panic(err)
		}
		marshal, _ := json.Marshal(createUser)
		writer.Write(marshal)
	}
}

func (config apiConfig) getByAPIKey() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.Header.Get("Authorization")
		userByAPIkey, err := config.DB.GetUserByAPIkey(request.Context(), apiKey)
		if err != nil {
			panic(err)
		}
		marshal, err := json.Marshal(userByAPIkey)
		if err != nil {
			panic(err)
		}
		writer.Write(marshal)
	}
}
func respondWithJSON(w http.ResponseWriter, status int, payload any) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func handlerReadiness() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		success := readiness{
			Status: "kulshi nadi",
		}
		respondWithJSON(w, 200, success)
	}
}

func handleError() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := error{Error: "Internal server error"}
		respondWithJSON(writer, 500, err)
	}
}

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
	config := apiConfig{DB: database.New(db)}
	defer db.Close()

	mainRouter := chi.NewRouter()
	v1 := chi.NewRouter()
	mainRouter.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))
	v1.Get("/readiness", handlerReadiness())
	v1.Get("/error", handleError())
	v1.Post("/users", config.newUser())
	v1.Get("/users", config.getByAPIKey())
	mainRouter.Mount("/v1", v1)
	server := http.Server{
		Handler: mainRouter,
		Addr:    fmt.Sprintf(":%v", port),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
