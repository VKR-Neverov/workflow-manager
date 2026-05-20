-- +goose Up
-- +goose StatementBegin
CREATE TABLE flows(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    state TEXT NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    completed_at TIMESTAMP,
    data JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE flows;
-- +goose StatementEnd
