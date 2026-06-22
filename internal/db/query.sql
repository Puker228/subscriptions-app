-- name: GetSubscription :one
SELECT * FROM subscriptions
WHERE id = $1 LIMIT 1;

-- name: CreateSubscription :exec
INSERT INTO subscriptions (
    id,
    service_name,
    price,
    user_id,
    start_date,
    end_date
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateSubscription :exec
UPDATE subscriptions
SET
    service_name = $2,
    price = $3,
    user_id = $4,
    start_date = $5,
    end_date = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE id = $1;
