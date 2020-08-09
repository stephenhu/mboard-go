FROM golang AS builder
WORKDIR /go/mboard-go
COPY . .
RUN GOOS=linux GOARCH=amd64 go get && go build

FROM ubuntu
WORKDIR /usr/local/mboard
RUN apt-get update -y
COPY --from=builder /go/mboard-go/mboard-go .
EXPOSE 8000
CMD ["/usr/local/mboard/mboard-go"]
