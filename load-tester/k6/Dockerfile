
# syntax=docker/dockerfile:1

FROM alpine:latest

ARG TARGETARCH
ARG K6_VERSION=0.45.0

# Install dependencies required to fetch and unpack k6
RUN apk add --no-cache curl tar

# Download and install k6 based on the architecture
RUN curl -L "https://github.com/grafana/k6/releases/download/v${K6_VERSION}/k6-v${K6_VERSION}-linux-${TARGETARCH}.tar.gz" -o k6.tar.gz \
    && tar xzf k6.tar.gz \
    && mv k6-v${K6_VERSION}-linux-${TARGETARCH}/k6 /usr/bin/ \
    && chmod +x /usr/bin/k6 \
    && rm -rf k6-v${K6_VERSION}-linux-${TARGETARCH} \
    && rm k6.tar.gz

# Continue with any other commands you need for your setup

