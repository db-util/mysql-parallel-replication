#!/usr/bin/env bash

source common.sh

docker ps  | grep mysql | awk '{print $1}' | xargs -I{} docker stop {}