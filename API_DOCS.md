
## Worker API

### Worker
`/worker/create`

Request
```bash
curl -X POST 'http://localhost:3010/worker/create' -d '
{
    "alias": "some alias"
}'
```

Response
```json
{
  "ok": true,
  "content": {
    "worker": {
      "id": 1,
      "created": 1559396382,
      "alias": "some alias",
      "secret": "ftZVO4w9Fc7bDuOISRaJL9P92ijkfvNah1Ldgc0a9f8=",
      "paused": false
    }
  }
}
```

----

`/worker/update`

Request
```bash
curl -X POST 'http://localhost:3010/worker/update' \
-H 'X-Worker-ID: 1' -H 'X-Secret: ftZVO4w9Fc7bDuOISRaJL9P92ijkfvNah1Ldgc0a9f8=' -d '
{
    "alias": "another alias"
}'
```

Response
```json
{
  "ok": true,
  "message": "(Error message, if applicable)"
}
```

### Tasks

`/task/submit`
Requires SUBMIT permissions on the project. max_assign_time is in seconds. Hash64 is
a 64-bit number - the submit will fail if another task has the same hash. If UniqueString
is specified, it will be hashed and put in place of Hash64.
 
VerificationCount is the number
of times the task has to be released with the same verification hash by *different workers* 
before the task is marked as closed. For example, a VerificationCount of 2 means that two
different workers have to assign and release the same task with the same verification hash before the
task can be marked as closed.
 

Request
```bash
curl -X POST 'http://localhost:3010/task/submit' \
-H 'X-Worker-ID: 1' -H 'X-Secret: ftZVO4w9Fc7bDuOISRaJL9P92ijkfvNah1Ldgc0a9f8=' -d '
{
    "project": 1,
    "max_retries": 3,
    "recipe": "test recipe",
    "priority": 1,
    "max_assign_time": 3600,
    "hash64": 0,
    "unique_string": "",
    "verification_count": 0
}'
```

Response
```json
{
  "ok": true,
  "message": "(Error message, if applicable)"
}
```

----
`/task/bulk_submit`

Same as `/task/submit`, but instead submit an array of submit requests. Tasks must
be for the same project. Follows the project's submit rate limit.


----
`/task/get/:project`

Request
```bash
curl -X GET 'http://localhost:3010/task/get/1'\
 -H 'X-Worker-ID: 1' -H 'X-Secret: ftZVO4w9Fc7bDuOISRaJL9P92ijkfvNah1Ldgc0a9f8='
```

Response
```json
{
  "ok": true,
  "content": {
    "task": {
      "id": 7,
      "priority": 1,
      "assignee": 1,
      "retries": 0,
      "max_retries": 3,
      "status": 1,
      "recipe": "test recipe",
      "max_assign_time": 0,
      "assign_time": 1559397394,
      "verification_count": 0,
      "project": {
        "id": 1,
        "priority": 999,
        "name": "My project",
        "clone_url": "http://github.com/test/test",
        "git_repo": "myrepo",
        "version": "1.0",
        "motd": "",
        "public": true,
        "hidden": false,
        "chain": 0,
        "paused": false,
        "assign_rate": 2,
        "submit_rate": 2
      }
    }
  }
}
```


----
`/task/release`

Result can be either of:

* *TR_OK*=0: Task was completed    
*  *TR_FAIL*=1: The worker failed to complete the task (tracker will mark the task as
FAILED after `max_retries` retries.   
* *TR_SKIP*=2: Act as if the worker never touched this task

Request
```bash
curl -X POST 'http://localhost:3010/task/release'\
 -H 'X-Worker-ID: 1' -H 'X-Secret: ftZVO4w9Fc7bDuOISRaJL9P92ijkfvNah1Ldgc0a9f8=' -d '
{ 
    "task_id": 7,
    "result": 0,
    "verification": -32315129
}'
```

Response

Updated will be set to false if their was an error (see Message) or if
the task has not yet reached the required number of verifications.
```json
{
  "ok": true,
  "message": "(Message, if applicable)",
  "content": {
    "updated": true
  }
}
```

### Logs

`/log/{trace|info|warn|error}`

Logs will be publicly available through the web UI 
if enabled in the config and will always appear in the log file.

Request
```bash
curl -X POST 'http://localhost:3010/log/info'\
 -H 'X-Worker-ID: 1' -H 'X-Secret: ftZVO4w9Fc7bDuOISRaJL9P92ijkfvNah1Ldgc0a9f8=' -d '
{ 
    "scope": "Arbitrary log category",
    "message": "message",
    "timestamp": 1559393394
}'
```

Response
```json
{
  "ok": true,
  "message": "(Error message, if applicable)"
}
```

## Projects

A=Assign (`/task/get/:project`)  
S=Submit (`/task/submit`)   
R=Read project from API (`project/get/:project`)

### Permissions
| Permission | Public | Private | Hidden |
| --- | --- | --- | --- |
| (none) | A R | R |   |
| `ASSIGN` | A R | A R | A  |
| `SUBMIT` | A S R | A S R | A S  |
| `ROLE_READ` (Manager) |  |  | R |
 
 
----
`/project/request_access`

A Manager with ROLE_MANAGE_ACCESS will need to manually approve the request from the UI.

Request
```bash
curl -X POST 'http://localhost:3010/project/request_access'\
 -H 'X-Worker-ID: 1' -H 'X-Secret: ftZVO4w9Fc7bDuOISRaJL9P92ijkfvNah1Ldgc0a9f8=' -d '
{
  "assign": true,
  "submit": true,
  "project": 23
}'
```

Response
```json
{
  "ok": true,
  "message": "(Error message, if applicable)"
}
```


----
`/project/get/"id`

Request

You must be authenticated as a manager and have ROLE_READ
 on the project to see hidden projects
```bash
curl -X GET 'http://localhost:3010/project/get/1'
```

Response
```json
{
  "ok": true,
  "message": "",
  "content": {
    "id": 1,
    "priority": 999,
    "name": "My project",
    "clone_url": "http://github.com/test/test",
    "git_repo": "myrepo",
    "version": "1.0",
    "motd": "",
    "public": true,
    "hidden": false,
    "chain": 0,
    "paused": false,
    "assign_rate": 2,
    "submit_rate": 2
  }
}
```

-----
`/project/list`

Hidden projects are returned if only if authenticated as a Manager and have the
permission to read them.

Request
```bash
curl -X GET 'http://localhost:3010/project/list'
```

Response
```json
{
  "ok": true,
  "message": "",
  "content": {
    "projects": []
  }
}
```

## UI / Manager API

Currently undocumented

```
/worker/set_paused
/worker/stats
/project/create
/project/monitoring-between/:id
/project/monitoring/:id
/project/assignees/:id
/project/access_list/:id
/project/request_access
/project/accept_request/:id/:wid
/project/reject_request/:id/:wid
/project/secret/:id
/project/secret/:id
/project/webhook_secret/:id
/project/webhook_secret/:id
/project/reset_failed_tasks/:id
/project/hard_reset/:id
/project/reclaim_assigned_tasks/:id
/git/receivehook
/logs
/register
/login
/logout
/account
/manager/list
/manager/list_for_project/:id
/manager/promote/:id
/manager/demote/:id
/manager/set_role_for_project/:id
```
