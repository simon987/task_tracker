
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-68%25-brightgreen.svg?longCache=true&style=flat-square)</a>
[![CodeFactor](https://www.codefactor.io/repository/github/simon987/task_tracker/badge)](https://www.codefactor.io/repository/github/simon987/task_tracker)
[![Build Status](https://ci.simon987.net/buildStatus/icon?job=task_tracker)](https://ci.simon987.net/job/task_tracker/)

Fast task tracker with authentication, statistics and web frontend



### Postgres setup
```bash
sudo su postgres
createuser task_tracker
createdb task_tracker
psql task_tracker
> 
```

### Nginx Setup

```nginx
index index.html;

root /path/to/webroot;

location / {
        try_files $uri $uri/ /index.html;
}
location ~ /api(.*)$ {
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_pass http://127.0.0.1:3010$1?$args; # Change host/port if necessary
}
```

### Running tests
```bash
cd test/
go test
```
