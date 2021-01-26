-- +goose Up
set schema 'users';

create table user_types
(
  code varchar(64) primary key,
  description varchar not null
);
insert into user_types values('client', 'клиент');
insert into user_types values('consultant', 'консультант');
insert into user_types values('expert', 'эксперт');
insert into user_types values('supervisor', 'супервизор');

create table users
(
  id uuid primary key,
  type varchar references user_types(code) not null,
  status varchar not null,
  username varchar not null,
  mm_id varchar,
  kk_id varchar,
  details jsonb,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null
);

create index idx_users_un on users(username);
create index idx_users_mm_id on users(mm_id);
create index idx_users_kk_id on users(kk_id);

-- +goose Down
drop table users;
drop table user_types;