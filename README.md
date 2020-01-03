
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-68%25-brightgreen.svg?longCache=true&style=flat-square)</a>
[![CodeFactor](https://www.codefactor.io/repository/github/simon987/task_tracker/badge)](https://www.codefactor.io/repository/github/simon987/task_tracker)

Fast task tracker (job queue) with authentication, statistics and web frontend

### Documentation

* Go client: [client](https://github.com/simon987/task_tracker/tree/master/client)
* Python client: [task_tracker_drone](https://github.com/simon987/task_tracker_drone)
* API specs: [API_DOCS](API_DOCS.md)
* Installation/Usage: [DOCS](DOCS.md)

### Features

* Stateless/Fault tolerent
* Integrate projects (or queue, tube) with Github/Gogs/Gitea - make workers aware of new commits
* Granular user permissions for administration tasks
* Prioritisable (project-level and task-level)
* Optionnal unique task constraint
* Per-project rate-limitting
* Per-project and per-worker stats monitoring

![image](https://user-images.githubusercontent.com/7120851/55676940-714cf980-58ac-11e9-8f5d-0d76a7afa80d.png)

### Terminology


**task_tracker** | Beanstalkd | Amazon SQS | IronMQ
:---|:---|:---|:---  
Project | Tube | Queue | Queue 
Task | Job | Message | Message
Recipe | Job data | Message body | Message body 
Submit | Put | Send message | POST
Assign | Reserve | Receive message | GET
Release | Delete | Delete message | DELETE
max_assign_time | TTR (time-to-run) | Visibility timeout | Timeout
\- | Delay | Delivery delay | Delay
\- | - | Retention Period | Expires in


### Running tests
```bash
cd test/
go test
```
