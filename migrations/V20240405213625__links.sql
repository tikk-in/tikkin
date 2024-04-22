CREATE TABLE links
(
    id          SERIAL PRIMARY KEY,
    user_id     BIGINT        NOT NULL REFERENCES users (id),
    slug        VARCHAR(255)  NOT NULL
        CONSTRAINT slug_not_empty CHECK (slug <> ''),
    description VARCHAR(1000)          DEFAULT NULL,
    banned      BOOLEAN       NOT NULL DEFAULT FALSE,
    expire_at   TIMESTAMP              DEFAULT NULL,
    target_url  VARCHAR(1000) NOT NULL
        CONSTRAINT target_url_not_empty CHECK (target_url <> ''),
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX links_slug_idx ON links (slug);

