-- +goose Up
alter table users add column api_key varchar(64) not null default encode(sha256(random()::text::bytea), 'hex');

update users set api_key = encode(sha256(random()::text::bytea), 'hex');

alter table users add constraint users_api_key unique (api_key);


-- +goose Down
ALTER TABLE users DROP CONSTRAINT users_api_key;

ALTER TABLE users DROP COLUMN api_key;