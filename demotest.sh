#!/bin/bash
# Copyright 2020 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -euo pipefail

port=8080
echo "Expects test-vmbackend running on default port: ${port}"

function call {
  method=$1
  url=$2
  echo "${method} ${url}"
  curl -s -X "${method}" "${url}"
  echo
}

function wait {
  max_wait=$1
  state=$2
  id=$3
  port=$4
  i=0
  while true
  do
      sleep 1
      i=$((i+1))
      if call GET "http://localhost:${port}/vms/${id}" | grep "${state}" > /dev/null
      then
        printf "\nVM %s is now %s\n" "${id}" "${state}" 
        break
      fi
      if (( i == max_wait ))
      then
        printf "FAILED TO get to state: %s\n" "${state}"
        exit 1
      fi
      echo -n .
  done
}

call GET http://localhost:${port}/vms
call GET http://localhost:${port}/vms/0
call PUT http://localhost:${port}/vms/0/launch
call GET http://localhost:${port}/vms/0
echo "Wait for started..."
wait 20 Running 0 "${port}"

call GET http://localhost:${port}/vms/0
call PUT http://localhost:${port}/vms/0/stop
call GET http://localhost:${port}/vms/0 
echo "Wait for stopped"
wait 20 Stopped 0 "${port}"

call GET http://localhost:${port}/vms/0
call DELETE http://localhost:${port}/vms/0
call GET http://localhost:${port}/vms

call GET http://localhost:${port}/ui/vms.html

echo "Demotest: OK/PASS"
