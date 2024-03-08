-- name: CreatePost :one
INSERT INTO posts (id, title, created_at, updated_at, url, feed_id, description, published_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING *;
-- name: GetPostsByUser :many
select * from posts p
         where p.feed_id in
               (select feed_id from feed_follow ff where ff.user_id = $1)
         order by p.published_at desc limit $2;