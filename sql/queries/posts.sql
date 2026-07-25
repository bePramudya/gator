-- name: CreatePost :one
INSERT INTO posts (
id,
created_at,
updated_at,
title,
url,
description,
published_at,
feed_id
)
values ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostForUser :many
SELECT * FROM posts
ORDER BY created_at DESC
LIMIT $1;
