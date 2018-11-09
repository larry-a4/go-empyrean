FROM golang:1.10.3-alpine AS build-env

# To copy complete directory
COPY ./ /go/src/github.com/ShyftNetwork/go-empyrean
WORKDIR /go/src/github.com/ShyftNetwork/go-empyrean

RUN \
  apk add --update git make gcc musl-dev linux-headers ca-certificates && \
  (cd /go/src/github.com/ShyftNetwork/go-empyrean && make geth && make bootnode) && \
  cp -v /go/src/github.com/ShyftNetwork/go-empyrean/build/bin/geth /bin && \
  cp -v /go/src/github.com/ShyftNetwork/go-empyrean/build/bin/bootnode /bin

FROM alpine:3.8
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /go/src/ShyftNetwork/go-empyrean/
COPY --from=build-env /bin/geth /bin/
COPY --from=build-env /bin/bootnode /bin/

COPY ./config.toml ./
COPY ./ShyftNetwork.json ./
COPY ./wait-for.sh ./
COPY ./shyft-cli/initShyftGeth.sh ./shyft-cli/

CMD ["./shyft-cli/initShyftGeth.sh", "geth --config config.toml"]

