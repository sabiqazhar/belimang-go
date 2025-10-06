CREATE TABLE IF NOT EXISTS merchants
(
    id                BIGINT PRIMARY KEY,
    name              VARCHAR(30) NOT NULL,
    merchant_category VARCHAR(20) NOT NULL,
    image_url         TEXT        NOT NULL,
    longitude         FLOAT       NOT NULL,
    latitude          FLOAT       NOT NULL,
    h3_index          BIGINT,
    created_at        TIMESTAMP DEFAULT NOW()
);