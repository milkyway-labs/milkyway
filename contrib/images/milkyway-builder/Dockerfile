FROM golang:1.22-alpine
ARG arch=x86_64

# Set up dependencies
ENV PACKAGES="curl make git libc-dev bash gcc linux-headers eudev-dev python3 ca-certificates build-base"
RUN set -eux; apk add --no-cache $PACKAGES;

# Set working directory for the build
WORKDIR /code

# Add source files
COPY . /code/

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.2/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep e78b224c15964817a3b75a40e59882b4d0e06fd055b39514d61646689cef8c6e

ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.2/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep e660a38efb2930b34ee6f6b0bb12730adccb040b6ab701b8f82f34453a426ae7

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.${arch}.a /usr/local/lib/libwasmvm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN BUILD_TAGS=muslc GOOS=linux GOARCH=amd64 LEDGER_ENABLED=true LINK_STATICALLY=true make build
RUN echo "Ensuring binary is statically linked ..." \
  && (file ./build/milkywayd | grep "statically linked")