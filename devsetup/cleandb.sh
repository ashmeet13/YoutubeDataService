#!/bin/sh

mongo $MONGO_DATABASE_NAME --eval "db.dropDatabase()"