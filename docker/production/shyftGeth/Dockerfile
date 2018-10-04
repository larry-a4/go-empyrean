FROM golang:1.10.3-alpine AS build-env

# To copy complete directory
COPY ./ /go/src/github.com/ShyftEthereum/go-empyrean
WORKDIR /go/src/github.com/ShyftEthereum/go-empyrean

RUN \
  apk add --update git make gcc musl-dev linux-headers ca-certificates && \
  (cd /go/src/github.com/ShyftEthereum/go-empyrean && make geth && make bootnode) && \
  cp -v /go/src/github.com/ShyftEthereum/go-empyrean/build/bin/geth /bin && \
  cp -v /go/src/github.com/ShyftEthereum/go-empyrean/build/bin/bootnode /bin 

FROM alpine:3.8
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /go/src/github.com/ShyftEthereum/go-empyrean
COPY --from=build-env /bin/geth /bin
COPY --from=build-env /bin/bootnode /bin

COPY ./shyft-cli/ ./shyft-cli/
COPY ./unlockPasswords.txt ./
COPY ./wait-for.sh ./
COPY ./ShyftNetwork.json ./
COPY ./config.toml ./

CMD ["./shyft-cli/initShyftGeth.sh", "./shyft-cli/startShyftGeth.sh"]
