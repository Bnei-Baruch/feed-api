#!/usr/bin/env bash
# Usage: cd databases/mdb;./update_mdb_models.sh
# Builds models package for the mdb, remove tests.

set -ev

rm -f ./models/*
sqlboiler -c ./sqlboiler.toml -o ./models psql
