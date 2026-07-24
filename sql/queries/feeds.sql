-- name: CreateFeed :one
INSERT INTO
    feeds (
        id,
        name,
        url,
        user_id,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: ListFeedsByUser :many
SELECT * FROM feeds WHERE user_id = $1 ORDER BY created_at DESC;

-- name: ListAllFeeds :many
SELECT feeds.id, feeds.name, feeds.url, feeds.user_id, feeds.created_at, feeds.updated_at, users.name AS user_name
FROM feeds
    INNER JOIN users ON feeds.user_id = users.id
ORDER BY feeds.created_at DESC;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
    last_fetched_at = $1,
    updated_at = NOW()
WHERE
    id = $2;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
WHERE
    last_fetched_at IS NULL
    OR last_fetched_at < NOW() - INTERVAL '1 hour'
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;