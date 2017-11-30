#!/usr/bin/env bash

source common.sh

source stop_mysql_docker.sh
source start_mysql_docker.sh

source start_replicator.sh

source stop_mysql_docker.sh