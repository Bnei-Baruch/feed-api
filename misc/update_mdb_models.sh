#!/usr/bin/env bash
# Usage: misc/update_mdb_models.sh
# Copy the models package from the mdb project, remove tests and rename the package.

set -ev

rm -f databases/mdb/models/*
cp  $GOPATH/src/github.com/Bnei-Baruch/archive-backend/mdb/models/*.go mdb/models
sed -i .bak 's/models/mdbmodels/' databases/mdb/models/*
rm databases/mdb/models/*.bak
rm databases/mdb/models/*_test.go

