-- noinspection SqlResolve FOR WHOLE FILE
-- noinspection SqlNoDataSourceInspection

-- name: CreateMerchant :one
INSERT INTO merchants (id, name, merchant_category, image_url, longitude, latitude, h3_index)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetMerchantList :many
SELECT
    id,
    name,
    merchant_category,
    image_url,
    latitude,
    longitude,
    created_at
FROM merchants
WHERE
  -- Filter by merchantId (optional)
    (sqlc.narg('merchant_id')::bigint IS NULL OR id = sqlc.narg('merchant_id')::bigint)

  -- Filter by name with wildcard, case insensitive (optional)
  AND (sqlc.narg('name')::text IS NULL OR LOWER(name) ILIKE LOWER(sqlc.narg('name')::text))

  -- Filter by category (optional)
  AND (sqlc.narg('merchant_category')::text IS NULL OR merchant_category = sqlc.narg('merchant_category')::text)

ORDER BY
    -- Dynamic sorting
    CASE
        WHEN sqlc.arg('sort_asc')::boolean = true THEN created_at
        END ASC,
    CASE
        WHEN sqlc.arg('sort_desc')::boolean = true THEN created_at
        END DESC,
    -- Default sort if neither asc nor desc
    CASE
        WHEN sqlc.arg('sort_asc')::boolean = false AND sqlc.arg('sort_desc')::boolean = false THEN created_at
        END DESC

LIMIT sqlc.arg('limit')::int
    OFFSET sqlc.arg('offset')::int;


-- name: AddItem :one
INSERT INTO items (id, merchant_id, name, price, image_url, product_category)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: GetItemList :many
SELECT id, name, product_category, price, image_url, created_at
FROM items
WHERE
        (sqlc.narg('itemId')::bigint IS NULL OR id = sqlc.narg('itemId')::bigint)
    AND
        (sqlc.narg('name')::text IS NULL OR LOWER(name) ILIKE LOWER(sqlc.narg('name')::text))
    AND
        (sqlc.narg('product_category')::text IS NULL OR product_category = sqlc.narg('product_category')::text)
ORDER BY
    CASE
        WHEN sqlc.arg('sort_asc')::boolean = true THEN created_at
        END ASC,
    CASE
        WHEN sqlc.arg('sort_desc')::boolean = true THEN created_at
        END DESC,
    CASE
        WHEN sqlc.arg('sort_asc')::boolean = false AND sqlc.arg('sort_desc')::boolean = false THEN created_at
        END DESC
LIMIT sqlc.arg('limit')::int
    OFFSET sqlc.arg('offset')::int;