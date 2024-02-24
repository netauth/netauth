# This is outputing either netauth or netauthd to /n and then it
# subsequently gets copied into a completely empty container.  This is
# because docker can't put ARGs into ENTRYPOINT, so that needs to be
# static.  Since in the general case the container isn't
# introsepctable without external tools (no shell!) this isn't
# noticed, but its still something to be aware of if you're here.

FROM golang:1.22-alpine as build
WORKDIR /netauth
COPY . .
ARG TARGET=netauthd
RUN go mod vendor && \
        CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /n cmd/$TARGET/main.go && \
        apk add upx binutils && \
        strip /n && \
        upx /n && \
        ls -alh /n

FROM scratch
LABEL org.opencontainers.image.source https://github.com/netauth/netauth
ENTRYPOINT ["/n"]
COPY --from=build /n /n
