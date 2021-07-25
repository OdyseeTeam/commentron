#!/usr/bin/env bash
export DEBUGGING=1
export MYSQL_DSN_RO="lbry-ro:lbry@tcp(localhost:3306)/commentron"
export MYSQL_DSN_RW="lbry-rw:lbry@tcp(localhost:3306)/commentron"
touch -a .env && set -o allexport; source ./.env; set +o allexport
