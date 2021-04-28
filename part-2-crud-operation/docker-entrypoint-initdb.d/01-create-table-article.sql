CREATE TABLE IF NOT EXISTS article (
    id  serial PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content text,
    created_at TIMESTAMP
);