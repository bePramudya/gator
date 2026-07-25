-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ( $1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.*, users.name AS user_name FROM feeds
LEFT JOIN users
    ON feeds.user_id = users.id;

-- name: GetFeed :one
SELECT feeds.*, users.name AS user_name FROM feeds
LEFT JOIN users
    ON feeds.user_id = users.id
WHERE feeds.url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $1, updated_at = $2
WHERE feeds.id = $3
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
