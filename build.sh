#!/bin/bash

set -e

./clean.sh

cd ./server
./swagger/swag init

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

rm -rf ./build

GOARCH=amd64
LDFLAGS="-s -w"
OUT=server-`date +%F@%T@%Z`

# For Docker, one build below for linux64 is enough.
OUTPATH_LINUX=./build/linux64/
mkdir -p $OUTPATH_LINUX
CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH_LINUX
cp -r ./www $OUTPATH_LINUX
echo "${G}server(linux64) built${W}"

OUTPATH_WIN=./build/win64/
mkdir -p $OUTPATH_WIN
CGO_ENABLED=0 GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
mv $OUT.exe $OUTPATH_WIN
cp -r ./www $OUTPATH_WIN
echo "${G}server(win64) built${W}"

# OUTPATH_MAC=./build/mac/
# mkdir -p $OUTPATH_MAC
# CGO_ENABLED=0 GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
# mv $OUT $OUTPATH_MAC
# cp -r ./www $OUTPATH_MAC
# echo "${G}server(mac) built${W}"

# GOARCH=arm
# OUTPATH_LARM=./build/linuxarm/
# mkdir -p $OUTPATH_LARM
# CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" GOARM=7 go build -ldflags="$LDFLAGS" -o $OUT
# mv $OUT $OUTPATH_LARM
# cp -r ./www $OUTPATH_LARM
# echo "${G}server(linuxArm) built${W}"

#######################################################################################

RELEASE_NAME=wisite-api.tar.gz 
cd ./build
tar -czvf $RELEASE_NAME ./linux64 ./win64  # ./mac ./linuxarm