#!/bin/bash

set -e

rm -rf ./server/docs

cd ./server
./swagger/swag init
cd -

if [[ $1 == 'all' ]]
then

go clean -cache
go clean -modcache
rm -f go.sum go.mod
go mod init github.com/wismed-web/wisite-api
go mod tidy

fi