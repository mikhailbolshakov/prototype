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

create table groups
(
  code varchar(64) primary key,
  type varchar(64) references user_types(code) not null,
  description varchar not null
);

insert into groups values('client', 'client', 'клиент');
insert into groups values('consultant-lawyer', 'consultant', 'консультант юрист');
insert into groups values('consultant-med', 'consultant', 'медконсультант');
insert into groups values('consultant', 'consultant', 'консультант по общим вопросам');
insert into groups values('doctor-dentist','expert', 'врач стоматолог');
insert into groups values('doctor-orthopedist','expert', 'врач ортопед');
insert into groups values('supervisor-rgs', 'supervisor', 'супервизор РГС');

create table user_groups
(
  id uuid primary key,
  user_id uuid references users(id) not null,
  type varchar(64) references user_types(code) not null,
  group_code varchar(64) references groups(code) not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null
);

create index idx_usr_grp_user_id on user_groups(user_id);

-- +goose Down
drop table user_groups;
drop table groups;
drop table users;
drop table user_types;