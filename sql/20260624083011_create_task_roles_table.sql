-- +goose Up
-- +goose StatementBegin
create table roles
(
    id         bigint auto_increment
        primary key,
    user_id    bigint    not null,
    role       bigint    not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint uni_roles_id
        unique (id),
    constraint fk_roles_user
        foreign key (user_id) references users (id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
