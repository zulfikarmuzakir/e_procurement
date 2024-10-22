CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

insert into users ("created_at", "email", "id", "name", "password", "role", "status", "updated_at", "username") values ('2024-10-17 03:05:48.632238+07', 'admin@example.com', 3, 'Admin', '$2a$14$2gybeaD9SNm/rE58jSewyOV2ewjH9A4oRrnWtkkuS/h9ArRWGYHUC', 'admin', 'active', '2024-10-17 03:05:48.632238+07', 'admin')