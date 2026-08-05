CREATE TABLE IF NOT EXISTS movies
(
    id
    UUID
    PRIMARY
    KEY
    DEFAULT
    gen_random_uuid
(
),
    source TEXT NOT NULL,
    external_id TEXT,
    fingerprint TEXT NOT NULL,

    title TEXT NOT NULL,
    description TEXT,
    poster_url TEXT,
    source_url TEXT NOT NULL,
    release_date DATE,

    notification_status TEXT NOT NULL DEFAULT 'pending',
    notification_attempts INTEGER NOT NULL DEFAULT 0,

    discovered_at TIMESTAMPTZ NOT NULL DEFAULT NOW
(
),
    notified_at TIMESTAMPTZ,

    raw_data JSONB,
    CONSTRAINT movies_source_fingerprint_unique
    UNIQUE
(
    source,
    fingerprint
)
    );

CREATE UNIQUE INDEX IF NOT EXISTS movies_source_external_id_unique
    ON movies (source, external_id)
    WHERE external_id IS NOT NULL;