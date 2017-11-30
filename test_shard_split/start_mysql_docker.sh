#!/usr/bin/env bash

source common.sh

DATETIME=`date +%Y%m%d_%H%M%S`

docker run \
    --name mysql_server_1_${DATETIME} \
    -p "$MYSQL_PORT:3306" \
    -e MYSQL_ROOT_PASSWORD="$MYSQL_PASSWORD" \
    -d mysql:5.7 \
    --log-bin \
    --binlog-format=ROW \
    --server-id=1

# init database
CONTAINER_ID=`docker ps  | grep "mysql_server_1_" | awk '{print $1}'`

set +e
until echo "SELECT 1" | docker exec -i ${CONTAINER_ID} mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD}
do
    echo "preparing mysql server..."
    sleep 1
done
set -e

cat init.sql | docker exec -i ${CONTAINER_ID} mysql -uroot -p${MYSQL_PASSWORD}
