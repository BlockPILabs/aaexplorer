#!/usr/bin/env bash

cwd=$(dirname $0)
cd $cwd
cd ../../../
pwd
set -x
make build
ssh root@18.211.227.200 systemctl stop aim
ssh root@18.211.227.200 systemctl stop aim-task
scp ./dist/aim root@18.211.227.200:/blockpi/aa-scan/aim
ssh root@18.211.227.200 systemctl start aim
ssh root@18.211.227.200 systemctl start aim-task

#./aim --home ./.aim start