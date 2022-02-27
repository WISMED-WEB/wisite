#!/bin/bash

set -e

./clean.sh

cd ./server
./swagger/swag init
go build