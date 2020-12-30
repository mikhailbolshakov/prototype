-- +goose Up
set schema 'tasks';

create table tasks
(
  id uuid primary key,
  num varchar not null,
  type varchar not null,
  subtype varchar not null,
  status varchar not null,
  substatus varchar not null,
  reported_by varchar not null,
  reported_at timestamp not null,
  due_date timestamp,
  assignee_group varchar,
  assignee_user varchar,
  assignee_at timestamp,
  description text,
  title varchar,
  details json,
  channel_id varchar,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null
);

create index idx_tasks_num on tasks(num);
create index idx_tasks_reported_by on tasks(reported_by);
create index idx_tasks_assign_user on tasks(assignee_user);
create index idx_tasks_channel on tasks(channel_id);

create table histories
(
  id uuid primary key,
  task_id uuid not null,
  status varchar not null,
  substatus varchar not null,
  assignee_group varchar,
  assignee_user varchar,
  assignee_at timestamp,
  changed_by varchar not null,
  changed_at timestamp not null
);

create index idx_histories_task on histories(task_id);

-- +goose Down
drop table tasks;
drop table histories;
