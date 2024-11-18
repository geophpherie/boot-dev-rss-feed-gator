
-- +goose Up
CREATE TABLE feeds (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_At TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	url TEXT NOT NULL,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	UNIQUE (url)
);

-- +goose Down
DROP TABLE feeds;
