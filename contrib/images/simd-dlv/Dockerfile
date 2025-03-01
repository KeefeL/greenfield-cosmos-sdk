FROM golang:1.18-alpine AS build

ARG GH_TOKEN=""

RUN apk add build-base git linux-headers libc-dev
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /work
COPY go.mod go.sum /work/
COPY db/go.mod db/go.sum /work/db/
COPY errors/go.mod errors/go.sum /work/errors/
COPY math/go.mod math/go.sum /work/math/

RUN go env -w GOPRIVATE="github.com/bnb-chain/*"
RUN git config --global url."https://${GH_TOKEN}@github.com".insteadOf "https://github.com"

RUN go mod download
COPY ./ /work
RUN LEDGER_ENABLED=false make COSMOS_BUILD_OPTIONS="debug,nostrip" clean build

FROM alpine AS run
RUN apk add bash curl jq libstdc++
COPY contrib/images/simd-dlv/wrapper.sh /usr/bin/wrapper.sh

VOLUME /simd
COPY --from=build /work/build/simd /simd/
COPY --from=build /go/bin/dlv /usr/local/bin
WORKDIR /simd

EXPOSE 26656 26657 2345
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start", "--log_format", "plain"]
STOPSIGNAL SIGTERM
