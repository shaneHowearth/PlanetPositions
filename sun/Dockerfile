ARG GO_VERSION=1.11

FROM golang:$GO_VERSION as builder

ENV GO111MODULE=on
ADD . $GOPATH/src/planetpositions
WORKDIR $GOPATH/src/planetpositions
# modules
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/sun sun/cmd/main.go

# run options
ENV PORT_NUM=5055
EXPOSE 5055
ENTRYPOINT ["sun"]
