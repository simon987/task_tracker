DROP TABLE IF EXISTS worker, project, task, log_entry,
  worker_has_access_to_project, manager, manager_has_role_on_project, project_monitoring_snapshot,
  worker_verifies_task, worker_requests_access_to_project;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS log_level;

CREATE TABLE worker
(
  id                SERIAL PRIMARY KEY NOT NULL,
  alias             TEXT               NOT NULL,
  created           INTEGER            NOT NULL,
  secret            BYTEA              NOT NULL,
  closed_task_count INTEGER DEFAULT 0  NOT NULL
);

CREATE TABLE project
(
  id                SERIAL PRIMARY KEY NOT NULL,
  priority          INTEGER DEFAULT 0  NOT NULL,
  closed_task_count INT     DEFAULT 0  NOT NULL,
  public            boolean            NOT NULL,
  hidden            boolean            NOT NULL,
  name              TEXT UNIQUE        NOT NULL,
  clone_url         TEXT               NOT NULL,
  git_repo          TEXT UNIQUE        NOT NULL,
  version           TEXT               NOT NULL,
  motd              TEXT               NOT NULL
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
  verification_hash BIGINT                                        NOT NULL,
  task              BIGINT REFERENCES task (id) ON DELETE CASCADE NOT NULL,
  worker            INT REFERENCES worker (id)                    NOT NULL
);

CREATE TABLE log_entry
(
  level        INTEGER NOT NULL,
  message      TEXT    NOT NULL,
  message_data TEXT    NOT NULL,
  timestamp    INTEGER NOT NULL
);

CREATE TABLE manager
(
  id            SERIAL PRIMARY KEY,
  register_time INTEGER     NOT NULL,
  tracker_admin BOOLEAN     NOT NULL,
  username      TEXT UNIQUE NOT NULL,
  password      BYTEA       NOT NULL
);

CREATE TABLE manager_has_role_on_project
(
  manager INTEGER REFERENCES manager (id) NOT NULL,
  role    SMALLINT                        NOT NULl,
  project INTEGER REFERENCES project (id) NOT NULL
);

CREATE TABLE project_monitoring_snapshot
(
  project                          INT REFERENCES project (id) NOT NULL,
  new_task_count                   INT                         NOT NULL,
  failed_task_count                INT                         NOT NULL,
  closed_task_count                INT                         NOT NULL,
  awaiting_verification_task_count INT                         NOT NULL,
  worker_access_count              INT                         NOT NULL,
  timestamp                        INT                         NOT NULL
);

CREATE TABLE worker_requests_access_to_project
(
  worker  INT REFERENCES worker (id)  NOT NULL,
  project INT REFERENCES project (id) NOT NULL
);

CREATE OR REPLACE FUNCTION on_task_delete_proc() RETURNS TRIGGER AS
$$
BEGIN
  UPDATE project SET closed_task_count=closed_task_count + 1 WHERE id = OLD.project;
  UPDATE worker SET closed_task_count=closed_task_count + 1 WHERE id = OLD.assignee;
  RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER on_task_delete
  BEFORE DELETE
  ON task
  FOR EACH ROW
EXECUTE PROCEDURE on_task_delete_proc();

CREATE OR REPLACE FUNCTION on_manager_insert() RETURNS TRIGGER AS
$$
BEGIN
  IF NEW.id = 1 THEN
    UPDATE manager SET tracker_admin= TRUE WHERE id = 1;
  end if;
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER on_manager_insert
  AFTER INSERT
  ON manager
  FOR EACH ROW
EXECUTE PROCEDURE on_manager_insert();

CREATE OR REPLACE FUNCTION release_task_ok(wid INT, tid INT, ver INT) RETURNS BOOLEAN AS
$$
DECLARE
  res INT = NULL;
BEGIN
  DELETE FROM task WHERE id = tid AND assignee = wid AND verification_count < 2 RETURNING project INTO res;

  IF res IS NULL THEN
    INSERT INTO worker_verifies_task (worker, verification_hash, task)
    SELECT wid, ver, task.id
    FROM task
    WHERE assignee = wid;

    DELETE
    FROM task
    WHERE id = tid
      AND assignee = wid
      AND (SELECT COUNT(*) as vcnt
           FROM worker_verifies_task wvt
           WHERE task = tid
           GROUP BY wvt.verification_hash
           ORDER BY vcnt DESC
           LIMIT 1) >= task.verification_count RETURNING task.id INTO res;

    IF res IS NULL THEN
      UPDATE task SET assignee= NULL WHERE id = tid AND assignee = wid;
    end if;
  end if;

  RETURN res IS NOT NULL;
END;
$$ LANGUAGE 'plpgsql';
