-- +goose Up
-- +goose StatementBegin
create table tasks
(
    id          bigint auto_increment
        primary key,
    assignee_id bigint    not null,
    team_id     bigint    not null,
    created_by  bigint    not null,
    created_at  timestamp not null,
    updated_at  timestamp not null,
    constraint fk_tasks_assignee_user
        foreign key (assignee_id) references users (id),
    constraint fk_tasks_created_by_user
        foreign key (created_by) references users (id),
    constraint fk_tasks_team
        foreign key (team_id) references teams (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
