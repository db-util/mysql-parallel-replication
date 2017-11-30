#!/usr/bin/env bash

set -e
set -x

PS4='$ Line ${LINENO}: '
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}

export MYSQL_USER=root
export MYSQL_PASSWORD=secret
export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=13306

