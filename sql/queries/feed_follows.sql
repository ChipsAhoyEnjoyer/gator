-- name: CreateFeedFollow :one
WITH inserted_feed_follow  AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING *
)
SELECT
    inserted_feed_follow.*, feed.name as feed_name, users.name as user_name
    FROM inserted_feed_follow
    INNER JOIN feed
    ON feed.id = inserted_feed_follow.feed_id
    INNER JOIN users
    ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.id, users.name AS user_name, feed.name AS feed_name
FROM feed_follows
INNER JOIN users
ON feed_follows.user_id = users.id
INNER JOIN feed
ON feed_follows.feed_id = feed.id
WHERE feed_follows.user_id = $1;