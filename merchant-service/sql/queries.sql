-- noinspection SqlResolve FOR WHOLE FILE
-- noinspection SqlNoDataSourceInspection

-- name: CreateMerchant :one
INSERT INTO merchants (id, name, merchant_category, image_url, longitude, latitude, h3_index)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetMerchantList :many
-- name: FindMerchants :many
SELECT * FROM merchants
WHERE
  -- If @merchant_id is NULL, this condition is ignored
    (@merchant_id::bigint IS NULL OR id = @merchant_id)
  AND
  -- If @name is NULL, this condition is ignored
    (@name::text IS NULL OR name ILIKE '%' || @name || '%')
  AND
  -- If @merchant_category is NULL, this condition is ignored
    (@merchant_category::text IS NULL OR merchant_category = @merchant_category)
ORDER BY
    -- This CASE statement handles dynamic sorting
    CASE WHEN @created_at_sort_asc::boolean THEN created_at END ASC,
    CASE WHEN @created_at_sort_desc::boolean THEN created_at END DESC
LIMIT sqlc.arg('limit')
    OFFSET sqlc.arg('offset');