FROM golang:1.11-alpine AS build-env

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
COPY ./unlockPasswords.txt ./
COPY ./shyft-cli/forBridgeDevelopment.sh ./shyft-cli/

EXPOSE 8545
CMD ["./shyft-cli/forBridgeDevelopment.sh", "geth --config config.toml","--disablepg" ," --unlock 0x43EC6d0942f7fAeF069F7F63D0384a27f529B062,0x9e602164C5826ebb5A6B68E4AFD9Cd466043dc4A,0x5Bd738164C61FB50eb12E227846CbaeF2dE965Aa,0xC04eE4131895F1d0C294D508AF65D94060AA42BB,0x07D899C4aC0c1725C35C5f816e60273B33a964F7 --password ./unlockPasswords.txt"]