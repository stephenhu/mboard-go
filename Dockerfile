FROM golang AS builder
WORKDIR /work/src/github.com/stephenhu/mboard-go
COPY . .
RUN go get && CGO_ENABLED=0 GOOS=linux go build

FROM scratch
WORKDIR /usr/local/mboard-go
COPY --from=builder /work/src/github.com/stephenhu/mboard-go/mboard-go .
EXPOSE 8000
CMD ["/usr/local/mboard-go/mboard-go"]
