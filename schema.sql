DROP TABLE IF EXISTS workerIdentity, Worker, Project, Task;
DROP TYPE IF EXISTS Status;

CREATE TYPE status as ENUM (
  'new',
  'failed',
  'closed'
  );

CREATE TABLE workerIdentity
(
  id          SERIAL PRIMARY KEY,
  remote_addr TEXT,

  UNIQUE (remote_addr)
);

CREATE TABLE worker
(
  id       TEXT PRIMARY KEY,
  created  INTEGER,
  identity INTEGER REFERENCES workerIdentity(id)
);

CREATE TABLE project
(
  id      SERIAL PRIMARY KEY,
  name    TEXT UNIQUE,
  git_url TEXT,
  version TEXT
);

CREATE TABLE task
(
  id          TEXT PRIMARY KEY,
  project     INTEGER REFERENCES project (id),
  assignee    TEXT REFERENCES worker (id),
  retries     INTEGER DEFAULT 0,
  max_retries INTEGER,
  status      Status  DEFAULT 'new'
);


