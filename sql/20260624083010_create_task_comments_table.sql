-- +goose Up
-- +goose StatementBegin
create table task_comments
(
    id         bigint auto_increment
        primary key,
    task_id    bigint    not null,
    user_id    bigint    not null,
    comment    text      not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint uni_task_comments_user_id
        unique (user_id),
    constraint fk_task_comments_task
        foreign key (task_id) references tasks (id),
    constraint fk_task_comments_user
        foreign key (user_id) references users (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_comments;
-- +goose StatementEnd
