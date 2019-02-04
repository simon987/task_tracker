DROP TABLE IF EXISTS worker_identity, worker, project, task, log_entry,
  worker_has_access_to_project, manager, manager_has_role_on_project, project_monitoring, worker_verifies_task;
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
  id                SERIAL PRIMARY KEY,
  priority          INTEGER DEFAULT 0,
  name              TEXT UNIQUE,
  clone_url         TEXT,
  git_repo          TEXT UNIQUE,
  version           TEXT,
  motd              TEXT,
  public            boolean,
  closed_task_count INT     DEFAULT 0
);

CREATE TABLE worker_has_access_to_project
(
  worker  INTEGER REFERENCES worker (id),
  project INTEGER REFERENCES project (id),
  primary key (worker, project)
);

CREATE TABLE task
(
  hash64             BIGINT   DEFAULT NULL UNIQUE,
  id                 SERIAL PRIMARY KEY,
  project            INTEGER REFERENCES project (id),
  assignee           INTEGER REFERENCES worker (id),
  max_assign_time    INTEGER  DEFAULT 0,
  assign_time        INTEGER  DEFAULT 0,
  verification_count INTEGER  DEFAULT 0,
  priority           SMALLINT DEFAULT 0,
  retries            SMALLINT DEFAULT 0,
  max_retries        SMALLINT,
  status             SMALLINT DEFAULT 1,
  recipe             TEXT
);

CREATE TABLE worker_verifies_task
(
  verification_hash BIGINT,
  task              BIGINT REFERENCES task (id) ON DELETE CASCADE,
  worker            INT REFERENCES worker (id)
);

CREATE TABLE log_entry
(
  level        INTEGER,
  message      TEXT,
  message_data TEXT,
  timestamp    INTEGER
);

CREATE TABLE manager
(
  id            SERIAL PRIMARY KEY,
  username      TEXT,
  password      TEXT,
  website_admin BOOLEAN
);

CREATE TABLE manager_has_role_on_project
(
  manager INTEGER REFERENCES manager (id),
  role    SMALLINT,
  project INTEGER REFERENCES project (id)
);

CREATE TABLE project_monitoring
(
  project           INT REFERENCES project (id),
  new_task_count    INT,
  failed_task_count INT,
  closed_task_count INT
);

CREATE OR REPLACE FUNCTION on_task_delete_proc() RETURNS TRIGGER AS
$$
BEGIN
  UPDATE project SET closed_task_count=closed_task_count + 1 WHERE id = OLD.project;
  RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER on_task_delete
  BEFORE DELETE
  ON task
  FOR EACH ROW
EXECUTE PROCEDURE on_task_delete_proc();
