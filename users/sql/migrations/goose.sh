#! /bin/sh

set -a
. ./.env
set +a

goose "$@"