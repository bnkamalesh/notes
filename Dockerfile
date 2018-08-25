FROM golang:1.10.3
COPY ${PWD}/ /go/src/github.com/bnkamalesh/notes/
WORKDIR /go/src/github.com/bnkamalesh/notes/
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o notes cmd/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/bnkamalesh/notes/notes .
CMD ["./notes"]
