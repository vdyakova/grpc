-- +goose Up

create table user_logs_table(
    id serial primary key,
    user_id int,
    log text,
    action Text not null,
    timestamp TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_action CHECK (action IN ('create', 'update', 'delete'))
);
-- +goose Down

drop table user_logs_table;