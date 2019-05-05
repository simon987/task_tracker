DROP TABLE IF EXISTS worker, project, task, log_entry,
    worker_access, manager, manager_has_role_on_project, project_monitoring_snapshot,
    worker_verifies_task;

CREATE TABLE worker
(
    id                SERIAL PRIMARY KEY NOT NULL,
    alias             TEXT               NOT NULL,
    created           INTEGER            NOT NULL,
    secret            BYTEA              NOT NULL,
    closed_task_count INTEGER                     DEFAULT 0 NOT NULL,
    paused            boolean            NOT NULL DEFAULT false
);

CREATE TABLE project
(
    id                SERIAL PRIMARY KEY NOT NULL,
    priority          INTEGER DEFAULT 0  NOT NULL,
    closed_task_count INT     DEFAULT 0  NOT NULL,
    chain             INT     DEFAULT NULL REFERENCES project (id),
    public            boolean            NOT NULL,
    hidden            boolean            NOT NULL,
    paused            boolean            NOT NULL,
    name              TEXT UNIQUE        NOT NULL,
    clone_url         TEXT               NOT NULL,
    git_repo          TEXT               NOT NULL,
    version           TEXT               NOT NULL,
    motd              TEXT               NOT NULL,
    secret            TEXT               NOT NULL DEFAULT '{}',
    webhook_secret    TEXT               NOT NULL,
    assign_rate       DOUBLE PRECISION   NOT NULL,
    submit_rate       DOUBLE PRECISION   NOT NULL
);

CREATE TABLE worker_access
(
    worker      INTEGER REFERENCES worker (id),
    project     INTEGER REFERENCES project (id),
    role_assign boolean,
    role_submit boolean,
    request     boolean,
    primary key (worker, project)
);

CREATE TABLE task
(
    hash64             BIGINT   DEFAULT NULL,
    id                 SERIAL PRIMARY KEY,
    project            INTEGER REFERENCES project (id),
    assignee           INTEGER REFERENCES worker (id),
    max_assign_time    INTEGER  DEFAULT 0,
    assign_time        INTEGER  DEFAULT NULL,
    verification_count INTEGER  DEFAULT 0,
    priority           SMALLINT DEFAULT 0,
    retries            SMALLINT DEFAULT 0,
    max_retries        SMALLINT,
    status             SMALLINT DEFAULT 1,
    recipe             TEXT,
    UNIQUE (project, hash64)
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
    role    SMALLINT                        NOT NULL,
    project INTEGER REFERENCES project (id) NOT NULL,
    PRIMARY KEY (manager, project)
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

CREATE OR REPLACE FUNCTION on_task_delete_proc() RETURNS TRIGGER AS
$$
DECLARE
    chain INTEGER;
BEGIN
    if OLD.assignee IS NOT NULL THEN
        UPDATE project
        SET closed_task_count=closed_task_count + 1
        WHERE id = OLD.project returning project.chain into chain;
        UPDATE worker SET closed_task_count=closed_task_count + 1 WHERE id = OLD.assignee;
        IF chain != 0 THEN
            INSERT into task (hash64, project, assignee, max_assign_time, assign_time, verification_count,
                              priority, retries, max_retries, status, recipe)
            VALUES (old.hash64, chain, NULL, old.max_assign_time, NULL,
                    old.verification_count, old.priority, 0, old.max_retries, 1,
                    old.recipe)
            ON CONFLICT DO NOTHING;
        end if;
    end if;
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

CREATE OR REPLACE FUNCTION release_task_ok(wid INT, tid INT, ver BIGINT) RETURNS BOOLEAN AS
$$
DECLARE
    res INT = NULL;
BEGIN
    DELETE FROM task WHERE id = tid AND assignee = wid AND verification_count < 2 RETURNING project INTO res;

    IF res IS NULL THEN
        INSERT INTO worker_verifies_task (worker, verification_hash, task)
        SELECT wid, ver, task.id
        FROM task
        WHERE assignee = wid
          AND task.id = tid;

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
