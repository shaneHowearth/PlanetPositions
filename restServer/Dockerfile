ARG GO_VERSION=1.11

FROM golang:$GO_VERSION as builder

ENV GO111MODULE=on
COPY restServer restServer$GOPATH/src/planetpositions/restServer
WORKDIR $GOPATH/src/planetpositions/restServer
ADD . $GOPATH/src/planetpositions
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

WORKDIR $GOPATH/src/planetpositions
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/restServer restServer/cmd/*.go

ENV portNum=5055
EXPOSE 5055
ENTRYPOINT ["restServer"]
