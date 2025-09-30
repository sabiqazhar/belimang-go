-- name: CreateMerchant :one
INSERT INTO merchants (id, name, merchant_category, image_url, longitude, latitude, h3_index)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;