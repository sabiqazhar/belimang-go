-- name: InsertImage :one
INSERT INTO images (
    id,
    image_url,
    upload_by
) VALUES (
    $1, $2, $3
)
RETURNING *;