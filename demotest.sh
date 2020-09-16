#!/bin/bash

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

call GET http://localhost:${port}/vms
call GET http://localhost:${port}/vms/0
call PUT http://localhost:${port}/vms/0/launch
call GET http://localhost:${port}/vms/0
echo "Wait for started..."
sleep 11

call GET http://localhost:${port}/vms/0
call PUT http://localhost:${port}/vms/0/stop
call GET http://localhost:${port}/vms/0 
echo "Wait for stopped"
sleep 6

call GET http://localhost:${port}/vms/0
call DELETE http://localhost:${port}/vms/0
call GET http://localhost:${port}/vms

echo "Demotest: OK/PASS"
