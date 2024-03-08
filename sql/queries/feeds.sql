-- name: CreateFeed :one
INSERT INTO feeds (id, name, created_at, updated_at, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedsByUser :many
SELECT  * FROM feeds where user_id = $1;

-- name: GetALlFeeds :many
SELECT  * FROM feeds;

-- name: CreateFeedFollow :one
INSERT INTO feed_follow (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllFeedFollowByUser :many
SELECT * FROM feed_follow where user_id = $1;

-- name: GetByFeedFollowById :one
SELECT * FROM feed_follow where id = $1;
-- name: DeleteFeedFollow :exec
DELETE FROM feed_follow WHERE id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY
    CASE
        WHEN last_fetched_at IS NULL THEN 0
        ELSE 1
        END,
    last_fetched_at ASC
    LIMIT $1;


-- name: MarkFeedFetched :exec
update feeds set last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP where id = $1;