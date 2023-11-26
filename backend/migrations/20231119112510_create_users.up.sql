CREATE TABLE users (
    id uuid not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null,
    created_at timestamp not null,
    updated_at timestamp not null
);