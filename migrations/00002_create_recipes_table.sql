-- +goose Up
-- +goose StatementBegin
CREATE TABLE recipes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ingredients TEXT NOT NULL,
    ai_response TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE INDEX idx_recipes_user_id ON recipes(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS recipes;
-- +goose StatementEnd
