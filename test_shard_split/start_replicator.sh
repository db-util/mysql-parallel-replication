#!/usr/bin/env bash

source common.sh

go build -v .
./`basename ${DIR}`
