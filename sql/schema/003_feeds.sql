-- +goose Up
CREATE TABLE feeds(
    id UUID primary key,
    name varchar(60) not null,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    url varchar(100) unique NOT NULL,
    user_id UUID not null
);

alter table feeds add constraint fk_feeds_user foreign key (user_id) references users(id) on delete cascade ;



-- +goose Down
alter table feeds drop constraint fk_feeds_user;
drop table feeds;