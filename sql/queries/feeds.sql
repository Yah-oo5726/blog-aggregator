-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedInfo :one
SELECT feeds.name, feeds.url, users.name
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id
WHERE feeds.id = $1;

-- name: GetByUrl :one
SELECT * FROM feeds WHERE url = $1;