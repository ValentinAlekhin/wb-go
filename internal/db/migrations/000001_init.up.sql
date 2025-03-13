create table if not exists virtual_controls
(
    topic      varchar  not null primary key,
    value      varchar  not null,
    created_at datetime not null,
    updated_at datetime not null
);