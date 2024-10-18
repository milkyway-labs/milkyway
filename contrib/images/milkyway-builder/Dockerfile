FROM golang:1.23-alpine
ARG arch=x86_64

# Set up dependencies
ENV PACKAGES="curl make git libc-dev bash gcc linux-headers eudev-dev python3 ca-certificates build-base"
RUN set -eux; apk add --no-cache $PACKAGES;

# Set working directory for the build
WORKDIR /code
COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Add source files
COPY . .

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.1.3/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep faea4e15390e046d2ca8441c21a88dba56f9a0363f92c5d94015df0ac6da1f2d

ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.1.3/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 8dab08434a5fe57a6fbbcb8041794bc3c31846d31f8ff5fb353ee74e0fcd3093

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.${arch}.a /usr/local/lib/libwasmvm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    BUILD_TAGS=muslc GOOS=linux GOARCH=amd64 LEDGER_ENABLED=true LINK_STATICALLY=true make build
RUN echo "Ensuring binary is statically linked ..." \
  && (file ./build/milkywayd | grep "statically linked")
