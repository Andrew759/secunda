-- +goose Up
-- +goose StatementBegin
create table users
(
    id         bigint auto_increment
        primary key,
    phone      varchar(30)  null,
    name       varchar(256) not null,
    surname    varchar(256) not null,
    login      varchar(256) not null,
    password   varchar(256) null,
    created_at timestamp    not null,
    updated_at timestamp    not null,
    constraint uni_users_id
        unique (id),
    constraint uni_users_login
        unique (login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
