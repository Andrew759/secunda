-- +goose Up
-- +goose StatementBegin
create table teams
(
    id         bigint auto_increment
        primary key,
    name       varchar(256) not null,
    created_by bigint       not null,
    created_at timestamp    not null,
    updated_at timestamp    not null,
    constraint fk_teams_created_by_user
        foreign key (created_by) references users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
