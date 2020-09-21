#!/bin/bash
docker build . -t test-vmbackend && docker run --rm --network=host test-vmbackend ${@}

