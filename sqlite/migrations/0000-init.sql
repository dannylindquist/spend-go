create table if not exists user(
  id integer primary key autoincrement,
  email text not null unique,
  password text not null,
  createdAt text not null default (strftime('%s', 'now')),
  updatedAt text not null default (strftime('%s', 'now'))
);