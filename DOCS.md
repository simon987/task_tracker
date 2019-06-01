# Documentation

## Installation (Docker)

Prerequisites:
* You have a postgres container using the network `tt`, listening
for connections on `172.26.0.2:5432`.
Example:
`docker run -d --name tt_pg --network tt postgres:alpine`


1. Initialize the database

    ```bash
    psql -h 172.26.0.2 -U postgres
    > CREATE USER task_tracker;
    > CREATE database task_tracker;
    ```
1. Write configuration file

    `vim config.yml`
    ```yaml
    server:
      address: "localhost:3010"
    database:
      conn_str: "postgres://task_tracker:task_tracker@172.26.0.2/task_tracker?sslmode=disable"
      log_levels: ["error", "info", "warn"]
    git:
      webhook_hash: "sha256"
      webhook_sig_header: "X-Gogs-Signature"
    log:
      level: "trace"
    session:
      cookie_name: "tt"
      expiration: "8h"
    monitoring:
      snapshot_interval: "120s"
      history_length: "400h"
    maintenance:
      reset_timed_out_tasks_interval: "5m"
    ```

1. Create task_tracker container:

    ```bash
    docker run --rm\
    	-v $PWD/config.yml:/root/config.yml\
    	--network tt\
    	-p 0.0.0.0:12345:80\
    	simon987/task_tracker
    ```
    
## Installation (Linux)

* You have a postgres daemon listening for connections on `localhost:5432`.
* You have a working installation of **go**, **nodejs** and **nginx**

1. Initialize the database

    ```bash
    sudo su postgres
    createuser task_tracker
    createdb task_tracker
    psql task_tracker
    > ALTER USER "task_tracker" WITH PASSWORD 'task_tracker';
    ```

1. Acquire binaries

    API
    ```bash
    go get -d github.com/simon987/task_tracker/...
    cd $GOPATH/src/github.com/simon987/task_tracker/main
    go build -o tt_api .
    ```

    UI
    ```bash
    git clone https://github.com/simon987/task_tracker
    cd task_tracker/web/angular
    npm install
    ./node_modules/\@angular/cli/bin/ng build --prod --optimization
    ```
   
1. Setup web server

    Move ./dist/ to /path/to/webroot/, start ./tt_api
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
            proxy_pass http://127.0.0.1:3010$1?$args;
    }
    ```

## Getting started

Register a *Manager* account from the /login page. The first account
will automatically be given tracker admin permissions.

Follow the instructions in the index page to create a project.

## API documentation

See [API_DOCS.md](API_DOCS.md)

