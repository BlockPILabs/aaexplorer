#!/usr/bin/env bash

cwd=$(dirname $0)
cd $cwd
cd ../../../
pwd

make build
ssh root@18.211.227.200 systemctl stop aim
scp ./dist/aim root@18.211.227.200:/blockpi/aa-scan/aim
ssh root@18.211.227.200 systemctl start aim

#./aim --home ./.aim start