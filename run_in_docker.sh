#!/bin/bash
# Copyright 2020 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

docker build . -t test-vmbackend && docker run --rm --network=host test-vmbackend ${@}

