FROM golang AS builder
WORKDIR /go/netc
COPY . .
RUN GOOS=linux GOARCH=amd64 go get && go build

FROM ubuntu
WORKDIR /usr/local/mboard
RUN apt-get update -y
COPY --from=builder /go/netc/mboard-go .
EXPOSE 8000
CMD ["/usr/local/mboard/mboard-go"]
