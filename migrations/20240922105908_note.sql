-- +goose Up
create table note (
                      id serial primary key,
                      name text not null,
                      email text not null,
                      role int not null,
                      created_at timestamp not null default now(),
                      updated_at timestamp
);


-- +goose Down
drop table note;