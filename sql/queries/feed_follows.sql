-- name: CreateFeedFollow :one
INSERT INTO
    feed_follows (
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    feed_follows.id,
    feed_follows.created_at,
    feed_follows.updated_at,
    feed_follows.user_id,
    feed_follows.feed_id,
    (
        SELECT name
        FROM feeds
        WHERE
            feeds.id = feed_follows.feed_id
    ) AS feed_name,
    (
        SELECT name
        FROM users
        WHERE
            users.id = feed_follows.user_id
    ) AS user_name;


-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.id,
    feed_follows.created_at,
    feed_follows.updated_at,
    feed_follows.user_id,
    feed_follows.feed_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    feed_follows.user_id = $1;