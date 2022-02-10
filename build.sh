#!/bin/bash

set -e

cd ./server
./swagger/swag init
go build