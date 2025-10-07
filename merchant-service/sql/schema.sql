CREATE TABLE IF NOT EXISTS merchants
(
    id                BIGINT PRIMARY KEY,
    name              VARCHAR(30) NOT NULL,
    merchant_category VARCHAR(50) NOT NULL,
    image_url         TEXT        NOT NULL,
    longitude         FLOAT       NOT NULL,
    latitude          FLOAT       NOT NULL,
    h3_index          BIGINT,
    created_at        TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS items
(
    id          BIGINT PRIMARY KEY,
    merchant_id BIGINT REFERENCES merchants(id) ON DELETE CASCADE,
    name        VARCHAR(50) NOT NULL,
    product_category VARCHAR(50) NOT NULL,
    image_url TEXT        NOT NULL,
    price       NUMERIC(10, 2) NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);