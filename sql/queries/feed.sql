-- name: PostFeed :one
INSERT INTO feed (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feed;

-- name: GetFeedByURL :one
SELECT * FROM feed
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feed SET last_fetched_at = $1, updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT * FROM feed
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;