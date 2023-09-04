#!/usr/bin/env bash

cwd=$(dirname $0)
cd $cwd
cd ../../../
pwd

make build

scp ./dist/aim root@18.211.227.200:/blockpi/aa-scan/aim

#./aim --home ./.aim start