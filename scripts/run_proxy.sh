#!/bin/bash

CURDIR=$(cd $(dirname $0); pwd)

SVC_NAME=ducker-proxy
BinaryName=ducker-proxy

if [[ -z "${DEPLOY_ENV}" ]]; then
  RUN_ENV="dev"
else
  RUN_ENV="${DEPLOY_ENV}"
fi

echo "$CURDIR/bin/${BinaryName} -conf=conf/config.${RUN_ENV}.yaml"
exec "$CURDIR"/bin/${BinaryName} -conf="conf/config.${RUN_ENV}.yaml"