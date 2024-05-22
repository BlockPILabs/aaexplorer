#!/usr/bin/env bash

cwd=$(dirname $0)
cd $cwd || exit
cd ../../../ || exit
pwd
set -x
make build
#ssh root@18.211.227.200 systemctl stop aim
#ssh root@18.211.227.200 systemctl stop aim-task

ssh -o StrictHostKeyChecking=no root@18.211.227.200 supervisorctl stop aim
ssh -o StrictHostKeyChecking=no root@18.211.227.200 supervisorctl stop aim-task
ssh -o StrictHostKeyChecking=no root@18.211.227.200 mkdir -p /blockpi/aaexplorer/log/
scp ./dist/aim root@18.211.227.200:/blockpi/aaexplorer/aim
ssh -o StrictHostKeyChecking=no root@18.211.227.200 supervisorctl start aim
ssh -o StrictHostKeyChecking=no root@18.211.227.200 supervisorctl start aim-task

#ssh root@18.211.227.200 systemctl start aim
#ssh root@18.211.227.200 systemctl start aim-task
ssh -o StrictHostKeyChecking=no root@162.55.130.133 supervisorctl stop aim
ssh -o StrictHostKeyChecking=no root@162.55.130.133 supervisorctl stop aim-task
ssh -o StrictHostKeyChecking=no root@162.55.130.133 mkdir -p /blockpi/aaexplorer/log/
scp ./dist/aim root@162.55.130.133:/blockpi/aaexplorer/aim
ssh -o StrictHostKeyChecking=no root@162.55.130.133 supervisorctl start aim
ssh -o StrictHostKeyChecking=no root@162.55.130.133 supervisorctl start aim-task

#./aim --home ./.aim start