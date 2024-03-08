
-- +goose Up
Create table feed_follow (
    id UUID primary key,
    feed_id UUID not null,
    user_id UUID not null,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null
);

alter table feed_follow add constraint fk_feed_feed_follow foreign key (feed_id) references feeds(id) on delete cascade;
alter table feed_follow add constraint fk_user_feed_follow foreign key (user_id) references users(id) on delete cascade;

-- +goose Down
alter table feed_follow drop constraint fk_feed_feed_follow;
alter table feed_follow drop constraint fk_user_feed_follow;
drop table feed_follow;
