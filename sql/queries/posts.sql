-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUserID :many
SELECT p.*
FROM posts p
INNER JOIN (
  select * from feed_follows 
  where user_id = $1
) ff on ff.feed_id = p.feed_id
ORDER BY p.published_at DESC
LIMIT $2;