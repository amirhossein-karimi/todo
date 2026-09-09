CREATE TABLE IF NOT EXISTS todos (
    id BIGSERIAL PRIMARY KEY,

    uuid UUID NOT NULL UNIQUE,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    priority smallint NOT NULL DEFAULT 0,

    status smallint NOT NULL DEFAULT 0,

    to_do_start_at TIMESTAMPTZ,

    do_start_at TIMESTAMPTZ,

    done_start_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

);

CREATE INDEX idx_todos_uuid
    ON todos(uuid);

CREATE INDEX idx_todos_status
    ON todos(status);

CREATE INDEX idx_todos_priority
    ON todos(priority);