FROM golang:1.14 as builder
ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/check ./cmd/check
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/in ./cmd/in

FROM alpine:3.11 as resource
LABEL MAINTAINER=mtharrison
COPY --from=builder /src/out /opt/resource

