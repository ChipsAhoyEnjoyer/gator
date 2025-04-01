-- name: CreateFeedFollow :many
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
    inserted_feed_follow.*, feed.name, users.name
    FROM inserted_feed_follow
    INNER JOIN feed
    ON feed.id = inserted_feed_follow.feed_id
    INNER JOIN users
    ON users.id = inserted_feed_follow.user_id;