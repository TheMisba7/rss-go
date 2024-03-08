-- +goose Up
CREATE TABLE posts(
    id UUID primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    title varchar(200) not null,
    url varchar(300) unique not null,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID not null
);

alter table posts add constraint fk_posts_feeds foreign key (feed_id) references feeds(id);
-- +goose Down
alter table posts drop constraint fk_posts_feeds;
DROP TABLE posts;
