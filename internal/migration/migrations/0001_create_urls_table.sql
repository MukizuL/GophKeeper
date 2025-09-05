-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
       
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login TEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    kdf_salt BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- +goose StatementEnd
    
-- +goose Down
DROP TABLE IF EXISTS users;