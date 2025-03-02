create table if not exists users_ (
  id_ integer primary key autoimcrement,
  public_key_ text not null
)

create table if not exists store_ (
  id_ integer primary key autoimcrement,
  key_ text not null unique,
  value_ text not null unique
)
