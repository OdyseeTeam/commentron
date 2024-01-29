#!/bin/bash
set -euo pipefail
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"
cd ".."
DIR="$PWD"
(
  cd "$DIR"
  go install github.com/volatiletech/sqlboiler/v4@v4.16.1
  go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.16.1
  sqlboiler --no-rows-affected --no-auto-timestamps --no-hooks --no-tests --no-context --wipe --add-soft-deletes mysql
)
