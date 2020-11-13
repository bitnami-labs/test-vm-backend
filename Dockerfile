# Copyright 2020 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

FROM golang:1.15.2-buster AS builder
WORKDIR /test-vmbackend
COPY *.go /test-vmbackend/
RUN go test && CGO_ENABLED=0 go build

FROM scratch
COPY --from=builder /test-vmbackend/test-vmbackend /
ENTRYPOINT ["/test-vmbackend"]
