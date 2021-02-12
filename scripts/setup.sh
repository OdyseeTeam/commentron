#!/usr/bin/env bash
export DEBUGGING=1
export MYSQL_DSN="lbry:lbry@tcp(localhost:3306)/commentron"
touch -a .env && set -o allexport; source ./.env; set +o allexport
