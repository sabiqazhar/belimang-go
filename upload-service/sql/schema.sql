CREATE TABLE images (
    id BIGINT PRIMARY KEY,
    image_url varchar(255),
    upload_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);