-- name: InsertOrder :one
INSERT INTO orders (id, customer_id, order_date, status, total_amount, estimated_delivery_time_in_minutes, longitude, latitude, total_distance_in_meters)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;

-- name: InsertOrderItems :copyfrom
INSERT INTO order_items (id, merchant_id, order_id, product_id, quantity, price, starting_point)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateOrderAmount :exec
UPDATE orders SET total_amount = $2 WHERE id = $1;

-- name: UpdateOrderStatus :exec
UPDATE orders SET status = $2 WHERE id = $1;