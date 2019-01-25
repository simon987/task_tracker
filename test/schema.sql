DROP TABLE IF EXISTS worker_identity, worker, project, task, log_entry,
  worker_has_access_to_project;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS log_level;

CREATE TYPE status as ENUM (
  'new',
  'failed',
  'closed',
  'timeout'
  );

CREATE TYPE log_level as ENUM (
  'fatal', 'panic', 'error', 'warning', 'info', 'debug', 'trace'
  );

CREATE TABLE worker_identity
(
  id          SERIAL PRIMARY KEY,
  remote_addr TEXT,
  user_agent  TEXT,

  UNIQUE (remote_addr)
);

CREATE TABLE worker
(
  id       TEXT PRIMARY KEY,
  alias    TEXT DEFAULT NULL,
  created  INTEGER,
  identity INTEGER REFERENCES workerIdentity (id)
);

CREATE TABLE project
(
  id        SERIAL PRIMARY KEY,
  priority  INTEGER DEFAULT 0,
  name      TEXT UNIQUE,
  clone_url TEXT,
  git_repo  TEXT UNIQUE,
  version   TEXT,
  motd      TEXT,
  public    boolean
);

CREATE TABLE worker_has_access_to_project
(
  worker  TEXT REFERENCES worker (id),
  project INTEGER REFERENCES project (id),
  primary key (worker, project)
);

CREATE TABLE task
(
  id              SERIAL PRIMARY KEY,
  priority        INTEGER DEFAULT 0,
  project         INTEGER REFERENCES project (id),
  assignee        TEXT REFERENCES worker (id),
  retries         INTEGER DEFAULT 0,
  max_retries     INTEGER,
  status          Status  DEFAULT 'new',
  recipe          TEXT,
  max_assign_time INTEGER DEFAULT 0,
  assign_time     INTEGER DEFAULT 0
);

CREATE TABLE log_entry
(
  level        log_level,
  message      TEXT,
  message_data TEXT,
  timestamp    INT
);


