-- name: InsertOrder :one
INSERT INTO orders (id, customer_id, order_date, status, total_amount, estimated_delivery_time_in_minutes)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: InsertOrderItem :one
INSERT INTO order_items (id, order_id, product_id, quantity, price, starting_point)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;