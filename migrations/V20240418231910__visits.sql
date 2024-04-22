CREATE TABLE visits
(
    id           VARCHAR(255) PRIMARY KEY,
    link_id      BIGINT      NOT NULL REFERENCES links (id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_agent   TEXT,
    referrer     TEXT,
    country_code VARCHAR(2)
);
