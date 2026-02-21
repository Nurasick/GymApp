-- +goose Up
-- +goose StatementBegin
CREATE TABLE training_plans (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_json TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE INDEX idx_training_plans_user_id ON training_plans(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS training_plans;
-- +goose StatementEnd
