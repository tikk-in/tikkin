CREATE TABLE users
(
    id                 BIGSERIAL PRIMARY KEY,
    email              VARCHAR(255) NOT NULL,
    password           VARCHAR(255) NOT NULL,
    verified           BOOLEAN     DEFAULT FALSE,
    verification_token VARCHAR(255),
    created_at         TIMESTAMPTZ DEFAULT NOW(),
    updated_at         TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX users_email_uindex ON users (email);
CREATE UNIQUE INDEX users_verification_token_uindex ON users (verification_token);
