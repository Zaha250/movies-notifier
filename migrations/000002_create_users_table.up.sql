CREATE TABLE users
(
    id                  UUID PRIMARY KEY     DEFAULT gen_random_uuid(),

    telegram_user_id    BIGINT      NOT NULL,
    telegram_chat_id    BIGINT      NOT NULL,

    username            TEXT,
    first_name          TEXT,
    last_name           TEXT,
    language_code       TEXT,

    is_bot              BOOLEAN     NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN     NOT NULL DEFAULT TRUE,

    started_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_interaction_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deactivated_at      TIMESTAMPTZ,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_telegram_user_id_unique
        UNIQUE (telegram_user_id),

    CONSTRAINT users_telegram_chat_id_unique
        UNIQUE (telegram_chat_id)
);