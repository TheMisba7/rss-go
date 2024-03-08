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

type feed struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Url       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func parseJson(req *http.Request, target interface{}) {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(target)
	if err != nil {
		panic(err)
	}
}

func (config *apiConfig) newUser() http.HandlerFunc {
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
		respondWithJSON(writer, 200, createUser)
	}
}

func (config *apiConfig) createFeed(w http.ResponseWriter, req *http.Request, user database.User) {
	now := time.Now()
	newFeed := feed{
		CreatedAt: now,
		UpdatedAt: now,
	}
	parseJson(req, &newFeed)
	newFeed.UserId = user.ID
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      newFeed.Name,
		CreatedAt: newFeed.CreatedAt,
		UpdatedAt: newFeed.UpdatedAt,
		UserID:    newFeed.UserId,
		Url:       newFeed.Url,
	}
	createFeed, err := config.DB.CreateFeed(req.Context(), feedParams)
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, 201, createFeed)
}
func (config *apiConfig) getByAPIKey() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.Header.Get("Authorization")
		userByAPIkey, err := config.DB.GetUserByAPIkey(request.Context(), apiKey)
		if err != nil {
			panic(err)
		}
		respondWithJSON(writer, 200, userByAPIkey)
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
func (config *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			respondWithError(writer, 401, "API Key not found.")
		}

		userByAPIkey, err := config.DB.GetUserByAPIkey(request.Context(), authHeader)
		if err != nil {
			return
		}
		handler(writer, request, userByAPIkey)
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

	//feeds
	v1.Post("/feeds", config.middlewareAuth(config.createFeed))

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
