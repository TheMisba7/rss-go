package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"rss-aggregator/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func parseJson(req *http.Request, target interface{}) {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(target)
	if err != nil {
		panic(err)
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

func (config *ApiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
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

func (config *ApiConfig) RegisterRoutes() http.Handler {
	v1 := chi.NewRouter()
	v1.Post("/users", config.newUser())
	v1.Get("/users", config.getByAPIKey())

	//feeds
	v1.Post("/feeds", config.middlewareAuth(config.createFeed))
	v1.Get("/feeds", config.getAllFeeds())

	//feed follow
	v1.Post("/feed_follow", config.middlewareAuth(config.followFeed))
	v1.Get("/feed_follow", config.middlewareAuth(config.getFeedFollowByUser))
	v1.Delete("/feed_follow/{feedFollowId}", config.middlewareAuth(config.unfollow))

	return v1
}
