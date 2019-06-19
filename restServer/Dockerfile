ARG GO_VERSION=1.11

FROM golang:$GO_VERSION as builder

ENV GO111MODULE=on
WORKDIR $GOPATH/src/planetpositions
ADD . $GOPATH/src/planetpositions
# modules
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/restServer restServer/cmd/*.go

# run options
ENV portNum=5055
EXPOSE 5055
ENTRYPOINT ["restServer"]
