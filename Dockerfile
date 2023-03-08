# Build the manager binary
FROM --platform=$BUILDPLATFORM ghcr.io/kedacore/build-tools:1.19.7 AS builder

ARG BUILD_VERSION=main
ARG GIT_COMMIT=HEAD
ARG GIT_VERSION=main

WORKDIR /workspace

COPY Makefile Makefile

# Copy the go source
COPY hack/ hack/
COPY version/ version/
COPY cmd/ cmd/
COPY apis/ apis/
COPY controllers/ controllers/
COPY pkg/ pkg/
COPY vendor/ vendor/
COPY go.mod go.mod
COPY go.sum go.sum

# Build
# https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide/
ARG TARGETOS
ARG TARGETARCH
RUN VERSION=${BUILD_VERSION} GIT_COMMIT=${GIT_COMMIT} GIT_VERSION=${GIT_VERSION} TARGET_OS=$TARGETOS ARCH=$TARGETARCH make manager

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/keda .
# 65532 is numeric for nonroot
USER 65532:65532

ENTRYPOINT ["/keda", "--zap-log-level=info", "--zap-encoder=console"]
