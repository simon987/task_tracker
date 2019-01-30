DROP TABLE IF EXISTS worker_identity, worker, project, task, log_entry,
  worker_has_access_to_project;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS log_level;

CREATE TABLE worker_identity
(
  id          SERIAL PRIMARY KEY,
  remote_addr TEXT,
  user_agent  TEXT,

  UNIQUE (remote_addr)
);

CREATE TABLE worker
(
  id       SERIAL PRIMARY KEY,
  alias    TEXT,
  created  INTEGER,
  identity INTEGER REFERENCES worker_identity (id),
  secret   BYTEA
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
  worker  INTEGER REFERENCES worker (id),
  project INTEGER REFERENCES project (id),
  primary key (worker, project)
);

CREATE TABLE task
(
  hash64          BIGINT   DEFAULT NULL UNIQUE,
  id              SERIAL PRIMARY KEY,
  project         INTEGER REFERENCES project (id),
  assignee        INTEGER REFERENCES worker (id),
  max_assign_time INTEGER  DEFAULT 0,
  assign_time     INTEGER  DEFAULT 0,
  priority        SMALLINT DEFAULT 0,
  retries         SMALLINT DEFAULT 0,
  max_retries     SMALLINT,
  status          SMALLINT DEFAULT 1,
  recipe          TEXT
);

CREATE TABLE log_entry
(
  level        INTEGER,
  message      TEXT,
  message_data TEXT,
  timestamp    INTEGER
);

