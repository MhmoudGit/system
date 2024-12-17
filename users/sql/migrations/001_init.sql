-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,                     -- Unique identifier for the user
    username VARCHAR(50) UNIQUE NOT NULL,      -- Username, must be unique
    email VARCHAR(255) UNIQUE NOT NULL,        -- Email, must be unique
    password TEXT NOT NULL,               -- Hashed password
    first_name VARCHAR(50),                    -- Optional first name
    last_name VARCHAR(50),                     -- Optional last name
    phone_number VARCHAR(15),                  -- Optional phone number
    is_active BOOLEAN DEFAULT FALSE,            -- Indicates if the account is active
    is_verified BOOLEAN DEFAULT FALSE,        -- Indicates if the account is verified
    role BIGINT NOT NULL,            -- User role
    created_at TIMESTAMP DEFAULT NOW(),        -- Timestamp of when the user was created
    updated_at TIMESTAMP DEFAULT NOW(),        -- Timestamp of the last update
    deleted_at TIMESTAMP DEFAULT NULL,        -- Timestamp of the last update
    CONSTRAINT chk_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$') -- Email format validation
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,                     -- Unique identifier for the role
    role_name VARCHAR(100) UNIQUE NOT NULL,    -- Name of the role
    permissions TEXT[] NOT NULL,               -- List of permissions (array of strings)
    created_at TIMESTAMP DEFAULT NOW(),        -- Timestamp of creation
    updated_at TIMESTAMP DEFAULT NOW(),        -- Timestamp of the last update
    deleted_at TIMESTAMP DEFAULT NULL          -- Timestamp for soft deletion
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,                     -- Unique identifier for the permission
    name VARCHAR(100) UNIQUE NOT NULL,         -- Name of the permission
    created_at TIMESTAMP DEFAULT NOW(),        -- Timestamp of creation
    updated_at TIMESTAMP DEFAULT NOW(),        -- Timestamp of the last update
    deleted_at TIMESTAMP DEFAULT NULL          -- Timestamp for soft deletion
);

-- CMD goose postgres "postgres://postgres:postgres@postgres:5432/auth" status -dir ./sql/migrations
-- CMD goose postgres "postgres://postgres:postgres@localhost:5432/auth" up -dir ./sql/migrations