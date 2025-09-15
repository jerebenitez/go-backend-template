BEGIN;
INSERT INTO migrations(ref) VALUES (1);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Authentication fields
    email CITEXT UNIQUE NOT NULL,
    email_verified_at TIMESTAMPTZ,
    password_hash TEXT NOT NULL,
    password_salt TEXT,

    -- Account status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_superuser BOOLEAN NOT NULL DEFAULT FALSE,
    last_login TIMESTAMPTZ,

    -- Security features
    mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    phone_number VARCHAR(20),
    phone_verified_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_users_email ON users (email);

COMMIT;
