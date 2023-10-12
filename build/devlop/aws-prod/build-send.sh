#!/usr/bin/env bash

cwd=$(dirname $0)
cd $cwd || exit
cd ../../../ || exit
pwd
set -x

version_dir=$(date +"%Y%m")
version_dir="/blockpi/aa-scan/bin/${version_dir}"

version_num=$(date +"%Y%m%d%H%M%S")

COMMIT_HASH=$(git rev-parse --short HEAD)

version="${version_num}.${COMMIT_HASH}"



make build

ssh -o StrictHostKeyChecking=no  root@ec2-54-85-44-73.compute-1.amazonaws.com mkdir -p /blockpi/aa-scan/log/ || exit 1
ssh -o StrictHostKeyChecking=no  root@ec2-54-85-44-73.compute-1.amazonaws.com mkdir -p ${version_dir} || exit 1
scp -o StrictHostKeyChecking=no  ./dist/aim root@ec2-54-85-44-73.compute-1.amazonaws.com:"${version_dir}/${version}" || exit 1
ssh -o StrictHostKeyChecking=no  root@ec2-54-85-44-73.compute-1.amazonaws.com ln -f -s ${version_dir}/${version} /blockpi/aa-scan/aim || exit 1

ssh -o StrictHostKeyChecking=no  root@ec2-54-85-44-73.compute-1.amazonaws.com supervisorctl restart aim
#ssh -o StrictHostKeyChecking=no  root@ec2-54-85-44-73.compute-1.amazonaws.com supervisorctl start aim-task
#ssh -o StrictHostKeyChecking=no  root@ec2-54-85-44-73.compute-1.amazonaws.com systemctl start aim
#ssh -o StrictHostKeyChecking=no  root@ec2-54-85-44-73.compute-1.amazonaws.com systemctl start aim-task

#./aim --home ./.aim start