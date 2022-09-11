FROM golang:1.19-alpine3.16

RUN mkdir -p "$GOPATH/src/"

ADD ./ /go/src

WORKDIR /go/src

RUN go build -v

EXPOSE 3000

CMD ["./YoutubeDataService"]