FROM golang:1.14 as builder
ADD . /src
WORKDIR /src
RUN make build-linux

FROM alpine:3.12
LABEL MAINTAINER=mtharrison
COPY --from=builder /src/out /opt/resource

