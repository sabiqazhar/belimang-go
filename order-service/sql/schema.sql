CREATE TABLE  IF NOT EXISTS orders (
    id BIGINT PRIMARY KEY,
    customer_id INT NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    estimated_delivery_time_in_minutes INT NOT NULL,
    total_distance_in_meters INT NOT NULL,
    longitude FLOAT,
    latitude FLOAT
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT PRIMARY KEY,
    merchant_id VARCHAR(50) NOT NULL,
    order_id BIGINT REFERENCES orders(id) ON DELETE CASCADE,
    product_id VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    starting_point BOOLEAN DEFAULT FALSE
);