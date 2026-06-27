-- +goose Up
-- +goose StatementBegin
create table team_members
(
    user_id    bigint    not null,
    team_id    bigint    not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_team_members_team
        foreign key (team_id) references teams (id),
    constraint fk_team_members_user
        foreign key (user_id) references users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team_members;
-- +goose StatementEnd
