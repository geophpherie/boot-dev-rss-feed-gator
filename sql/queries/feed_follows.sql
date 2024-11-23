-- insert a feed follow record, return all fields of feed_follow AND name of linked user and 
-- name: CreateFeedFollows :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
) SELECT 
    iff.*,
    f.name as feed_name,
    u.name as user_name
  FROM inserted_feed_follow as iff
    INNER JOIN feeds as f ON iff.feed_id = f.id
    INNER JOIN users as u ON iff.user_id = u.id;

-- name: DeleteAllFeedFollows :exec
DELETE FROM feed_follows;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;

-- name: GetAllFeedFollowsByUser :many
SELECT
    ff.*,
    f.name as feed_name,
    u.name as user_name
FROM feed_follows as ff
    INNER JOIN feeds as f ON ff.feed_id = f.id
    INNER JOIN users as u ON ff.user_id = u.id
WHERE ff.user_id = $1;
