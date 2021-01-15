-- +goose Up
set schema 'services';

create table service_types (
  id varchar primary key,
  description varchar not null,
  delivery_wf_id varchar
);

insert into service_types values ('expert-online-consultation', 'Онлайн-консультация с экспертом', 'expert_online_consultation');
insert into service_types values ('medical-checkup', 'Медицинское обследование', null);

create table balances (
  id uuid primary key,
  client_id uuid not null,
  service_type_id varchar not null,
  total integer not null,
  delivered integer not null,
  locked integer not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null
);

create index idx_srv_balance_client on balances(client_id);

create table deliveries (
  id uuid primary key,
  client_id uuid not null,
  service_type_id varchar not null,
  status varchar not null,
  start_time timestamp not null,
  finish_time timestamp,
  details json,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null
);

create index idx_srv_deliveries_client on deliveries(client_id);

create table deliveries_tasks (
  delivery_id uuid not null,
  task_id uuid not null,
  primary key(delivery_id, task_id)
);

create index idx_del_tsk_delivery on deliveries_tasks(delivery_id);
create index idx_del_tsk_task on deliveries_tasks(task_id);

-- +goose Down
drop table service_types;
drop table balances;
drop table deliveries;
drop table deliveries_tasks;
