#!/bin/bash
set -euo pipefail

mkdir -p /tmp/commentron/init.d
tee /tmp/commentron/init.d/001_init.sql <<EOF
CREATE DATABASE IF NOT EXISTS commentron;
CREATE DATABASE IF NOT EXISTS social;
CREATE USER 'lbry-ro'@'%' IDENTIFIED BY 'lbry';
CREATE USER 'lbry-rw'@'%' IDENTIFIED BY 'lbry';
GRANT ALL ON commentron.* TO 'lbry-rw'@'%';
GRANT SELECT ON commentron.* TO 'lbry-ro'@'%';
GRANT ALL ON social.* TO 'lbry-rw'@'%';
FLUSH PRIVILEGES;
EOF

docker run --rm -it -p 3306:3306 \
  --name tmp-commentron-mysql \
  -e MYSQL_ROOT_PASSWORD=lbry \
  -v /tmp/commentron/init.d:/docker-entrypoint-initdb.d \
  mysql/mysql-server:8.0


# Verify with:
# mysql -h localhost -u lbry-rw -p commentron --protocol tcp