#!/usr/bin/env bash
export DEBUGGING=1
export MYSQL_DSN_RO="commentron-ro:commentron@tcp(localhost:3306)/commentron"
export MYSQL_DSN_RW="commentron-rw:commentron@tcp(localhost:3306)/commentron"
touch -a .env && set -o allexport; source ./.env; set +o allexport
