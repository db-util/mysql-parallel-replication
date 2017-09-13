#!/usr/bin/env bash

set -e
set -x

PS4='\n$ Line ${LINENO}: '
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}
#APP_NAME=`basename ${DIR}`
DATETIME=`date +%Y%m%d_%H%M%S`

# load configurations
source ./env.sh

# restart mysql docker
docker ps  | grep mysql | awk '{print $1}' | xargs -I{} docker stop {}

docker run \
    --name mysql_server_1_${DATETIME} \
    -p "$REPL_SRC_PORT:3306" \
    -e MYSQL_ROOT_PASSWORD="$MYSQL_ROOT_PASSWORD" \
    -d mysql:5.7 \
    --log-bin \
    --server-id=1

# init database
CONTAINER_ID=`docker ps  | grep "mysql_server_1_" | awk '{print $1}'`

set +e
while [ "`docker logs ${CONTAINER_ID} 2>&1  | grep -o "ready for connections" | wc -l`" -lt "2" ] ; do
    echo "waiting for mysql server..."
    sleep 1
done
set -e

cat sql/init.sql | docker exec -i ${CONTAINER_ID} mysql -uroot -p${MYSQL_ROOT_PASSWORD}


go build -v .
./`basename ${DIR}`
