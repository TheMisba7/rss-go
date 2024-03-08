package handlers

import (
	"net/http"
	"rss-aggregator/internal/database"
	"strconv"
)

var defaultSize int32 = 100

func (config *ApiConfig) getPostsByUser(w http.ResponseWriter, req *http.Request, user database.User) {
	var size = defaultSize
	if req.URL.Query().Get("size") != "" {
		parsed, err := strconv.Atoi(req.URL.Query().Get("size"))
		if err != nil {
			size = defaultSize
		} else {
			size = int32(parsed)
		}
	}
	posts, _ := config.DB.GetPostsByUser(req.Context(), database.GetPostsByUserParams{UserID: user.ID, Limit: size})
	respondWithJSON(w, 200, posts)
}
