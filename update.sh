#!/bin/bash

set -e

rm -rf ./server/docs

cd ./server
./swagger/swag init

cd -

rm go.sum

go get -u ./...