# To build the Milkyway image, just run:
# > docker build -t milkyway .
#
# Simple usage with a mounted data directory:
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.milkywayd:/root/.milkywayd milkyway milkywayd init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.milkywayd:/root/.milkywayd milkyway milkywayd start
#
# If you want to run this container as a daemon, you can do so by executing
# > docker run -td -p 26657:26657 -p 26656:26656 -v ~/.milkywayd:/root/.milkywayd --name milkyway milkywayd
#
# Once you have done so, you can enter the container shell by executing
# > docker exec -it milkyway bash
#
# To exit the bash, just execute
# > exit
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates

# Install bash
RUN apk add --no-cache bash

# Copy over binaries from the build-env
COPY --from=milkywaylabs/builder:latest /code/build/milkywayd /usr/bin/milkywayd

EXPOSE 26656 26657 1317 9090

# Run milkywayd by default, omit entrypoint to ease using container with milkywayd
CMD ["milkywayd"]