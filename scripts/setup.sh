#!/usr/bin/env bash
export DEBUGGING=1
export MYSQL_DSN_RO="odysee:thomas@tcp(localhost:3306)/commentron"
export MYSQL_DSN_RW="odysee:thomas@tcp(localhost:3306)/commentron"
touch -a .env && set -o allexport; source ./.env; set +o allexport
