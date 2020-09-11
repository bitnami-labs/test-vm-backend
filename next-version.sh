#!/bin/bash

set -euo pipefail

number=0

previous=$(cat VERSION)
ptimestamp=$(echo "${previous}" | awk -F. '{print $1"."$2"."$3}')
pnumber=$(echo "${previous}" | awk -F. '{print $4}')
timestamp=$(date +"%Y.%m.%d")
if [ "${timestamp}" = "${ptimestamp}" ]; then
  number=$(( pnumber + 1))
fi

date +"%Y.%m.%d.${number}" > VERSION
echo "Next version is now:"
cat VERSION
