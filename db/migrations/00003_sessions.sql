-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
                          id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
                          token_hash TEXT        NOT NULL UNIQUE,
                          user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                          expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd