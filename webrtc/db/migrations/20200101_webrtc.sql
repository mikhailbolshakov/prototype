-- +goose Up
set schema 'webrtc';

create table rooms
(
  id uuid primary key,
  details jsonb,
  opened_at timestamp,
  closed_at timestamp,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null
);

-- +goose Down
drop table rooms;

