# Documentation

## Installation (Docker)

1. *(Optional)* Tweak configuration file
    `vim config.yml`

1. 
    `docker-compose up`

    
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

