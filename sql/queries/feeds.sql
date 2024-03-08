-- name: CreateFeed :one
INSERT INTO feeds (id, name, created_at, updated_at, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedsByUser :many
SELECT  * FROM feeds where user_id = $1;

-- name: GetALlFeeds :many
SELECT  * FROM feeds;