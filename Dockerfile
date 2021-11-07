FROM golang:1.17-alpine

WORKDIR /go/src/app
# copy src
COPY . .
# install deps
RUN go get -d -v ./...
# compile binary
RUN go build -o disbot
# run binary
CMD ["./disbot"]
