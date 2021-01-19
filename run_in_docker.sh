#!/bin/bash
# Copyright 2020 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

docker build . -t test-vm-backend && docker run --rm -p 8080:8080 test-vm-backend ${@}

