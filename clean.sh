#!/bin/bash

set -e

rm -rf ./data

rm -rf ./server/__debug_bin
rm -rf ./server/server
rm -rf ./server/tmp*
rm -rf ./server/data

rm -rf ./prelease/prelease

if [[ $1 == 'all' ]] 
then

    rm -rf ./server/build

else

    rm -rf ./server/build/linux64/tmp*
    rm -rf ./server/build/linux64/server*
    rm -rf ./server/build/win64/tmp*
    rm -rf ./server/build/win64/server*
    
fi

rm -rf ./server/build/*.gz
