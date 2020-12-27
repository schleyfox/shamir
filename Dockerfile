FROM ubuntu:20.04 AS build

ENV DEBIAN_FRONTEND=noninteractive

# 1.14 is the latest official go in 20.04. Installing via apt to avoid arch
# issues.
ARG GO_VERSION=1.14

RUN apt-get update && \
    apt-get install -y make golang-${GO_VERSION}-go curl

ENV PATH="/usr/lib/go-${GO_VERSION}/bin:$PATH"

WORKDIR /build

COPY . .

RUN make all

FROM ubuntu:20.04

COPY --from=build /build/bin/shamir /bin/shamir
