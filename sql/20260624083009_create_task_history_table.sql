-- +goose Up
-- +goose StatementBegin
create table task_history
(
    id         bigint auto_increment
        primary key,
    task_id    bigint       not null,
    changed_by bigint       not null,
    team_id    bigint       not null,
    created_by bigint       not null,
    name       varchar(256) not null,
    created_at timestamp    not null,
    updated_at timestamp    not null,
    constraint fk_task_history_changed_by_user
        foreign key (changed_by) references users (id),
    constraint fk_task_history_created_by_user
        foreign key (created_by) references users (id),
    constraint fk_task_history_task
        foreign key (task_id) references tasks (id),
    constraint fk_task_history_team
        foreign key (team_id) references teams (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_history;
-- +goose StatementEnd
