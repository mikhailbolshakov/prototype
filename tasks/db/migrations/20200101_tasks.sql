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
  reported_type varchar not null,
  reported_user_id varchar not null,
  reported_username varchar not null,
  reported_at timestamp not null,
  due_date timestamp,
  assignee_type varchar,
  assignee_group varchar,
  assignee_user_id varchar,
  assignee_username varchar,
  assignee_at timestamp,
  description text,
  title varchar,
  details json,
  reminders json,
  channel_id varchar,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null
);

create index idx_tasks_num on tasks(num);
create index idx_tasks_reported_usrid on tasks(reported_user_id);
create index idx_tasks_reported_usrname on tasks(reported_username);
create index idx_tasks_assign_usrid on tasks(assignee_user_id);
create index idx_tasks_assign_usrname on tasks(assignee_username);
create index idx_tasks_channel on tasks(channel_id);

create table histories
(
  id uuid primary key,
  task_id uuid not null,
  status varchar not null,
  substatus varchar not null,
  assignee_type varchar,
  assignee_group varchar,
  assignee_user_id varchar,
  assignee_username varchar,
  assignee_at timestamp,
  changed_by varchar not null,
  changed_at timestamp not null
);

create index idx_histories_task on histories(task_id);

create table assignment_logs (
    id uuid primary key,
    start_time timestamp not null,
    finish_time timestamp,
    status varchar not null,
    rule_code varchar not null,
    rule_description varchar not null,
    users_in_pool integer not null,
    tasks_to_assign integer not null,
    assigned integer not null,
    error varchar
);

create index idx_ass_log_start_time on assignment_logs(start_time);
create index idx_ass_log_fin_time on assignment_logs(finish_time);

-- +goose Down
drop table tasks;
drop table histories;
drop table assignment_logs;
