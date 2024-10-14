FROM golang:1.22-alpine
ARG arch=x86_64

# Set up dependencies
ENV PACKAGES="curl make git libc-dev bash gcc linux-headers eudev-dev python3 ca-certificates build-base"
RUN set -eux; apk add --no-cache $PACKAGES;

# Set working directory for the build
WORKDIR /code
COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Add source files
COPY . .

# See https://github.com/CosmWasm/wasmvm/releases
ENV LIBWASMVM_VERSION=v2.1.2
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.1.2/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.1.2/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 0881c5b463e89e229b06370e9e2961aec0a5c636772d5142c68d351564464a66
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 58e1f6bfa89ee390cb9abc69a5bc126029a497fe09dd399f38a82d0d86fe95ef

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.${arch}.a /usr/local/lib/libwasmvm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    BUILD_TAGS=muslc GOOS=linux GOARCH=amd64 LEDGER_ENABLED=true LINK_STATICALLY=true make build
RUN echo "Ensuring binary is statically linked ..." \
  && (file ./build/milkywayd | grep "statically linked")
