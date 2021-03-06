#!/usr/bin/env bash
# Usage: misc/update_mdb_models.sh
# Copy the models package from the mdb project, remove tests and rename the package.

set -ev

rm -f mdb/models/*
cp  $GOPATH/src/github.com/Bnei-Baruch/archive-backend/mdb/models/*.go mdb/models
sed -i .bak 's/models/mdbmodels/' mdb/models/*
rm mdb/models/*.bak
rm mdb/models/*_test.go

