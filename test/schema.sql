DROP TABLE IF EXISTS workeridentity, Worker, Project, Task, log_entry;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS loglevel;

CREATE TYPE status as ENUM (
  'new',
  'failed',
  'closed'
  );

CREATE TYPE loglevel as ENUM (
  'fatal', 'panic', 'error', 'warning', 'info', 'debug', 'trace'
  );

CREATE TABLE workerIdentity
(
  id          SERIAL PRIMARY KEY,
  remote_addr TEXT,
  user_agent  TEXT,

  UNIQUE (remote_addr)
);

CREATE TABLE worker
(
  id       TEXT PRIMARY KEY,
  created  INTEGER,
  identity INTEGER REFERENCES workerIdentity (id)
);

CREATE TABLE project
(
  id        SERIAL PRIMARY KEY,
  priority  INTEGER DEFAULT 0,
  motd      TEXT    DEFAULT '',
  name      TEXT UNIQUE,
  clone_url TEXT,
  git_repo  TEXT UNIQUE,
  version   TEXT
);

CREATE TABLE task
(
  id          SERIAL PRIMARY KEY,
  priority    INTEGER DEFAULT 0,
  project     INTEGER REFERENCES project (id),
  assignee    TEXT REFERENCES worker (id),
  retries     INTEGER DEFAULT 0,
  max_retries INTEGER,
  status      Status  DEFAULT 'new',
  recipe      TEXT
);

CREATE TABLE log_entry
(
  level        loglevel,
  message      TEXT,
  message_data TEXT,
  timestamp    INT
);


