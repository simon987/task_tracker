#!/usr/bin/env bash

export INSTALL_DIR="/home/drone/task_tracker"

mkdir ${INSTALL_DIR} 2> /dev/null

# Gogs
if [[ ! -d "${INSTALL_DIR}/gogs" ]]; then
    wget "https://dl.gogs.io/0.11.79/gogs_0.11.79_linux_amd64.tar.gz"
    tar -xzf "gogs_0.11.79_linux_amd64.tar.gz" -C ${INSTALL_DIR}
    rm "gogs_0.11.79_linux_amd64.tar.gz"
fi

# Postgres
su - postgres -c "createuser task_tracker"
su - postgres -c "dropdb gogs"
su - postgres -c "createdb gogs"
