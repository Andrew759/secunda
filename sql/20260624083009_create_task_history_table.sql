-- +goose Up
-- +goose StatementBegin
create table task_history
(
    id         bigint auto_increment
        primary key,
    task_id    bigint    not null,
    changed_by bigint    not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_task_history_changed_by_user
        foreign key (changed_by) references users (id),
    constraint fk_task_history_task
        foreign key (task_id) references tasks (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_history;
-- +goose StatementEnd
