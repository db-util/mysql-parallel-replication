#!/usr/bin/env bash

set -e
set -x

PS4='$ Line ${LINENO}: '

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}
APP_NAME=`basename ${DIR}`

source ./env.sh

docker ps  | grep mysql_repl_src | awk '{print $1}' | xargs -I{} docker stop {}
docker build -f mysql_repl_src.dockerfile -t mysql_repl_src .
CONTAINER_ID=`docker run -p "$REPL_SRC_PORT:3306" -e MYSQL_ROOT_PASSWORD="$MYSQL_ROOT_PASSWORD" -d mysql_repl_src:latest`

CONTAINER_ID=`docker ps  | grep mysql_repl_src | awk '{print $1}'`
#sleep 10
#echo ${CONTAINER_ID}

go build -v .
./${APP_NAME}

echo "show databases" | docker exec -i ${CONTAINER_ID} mysql -uroot -psecret
echo "show tables from adt_test" | docker exec -i ${CONTAINER_ID} mysql -uroot -psecret
