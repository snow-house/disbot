FROM golang:1.17-alpine

WORKDIR /go/src/app
# copy src
COPY . .

# Download pictures 
RUN wget -O vvimg.jpg https://i.imgur.com/yiXXYMy.jpg

# TODO: download font and save it 
# install zipo
RUN apk update
RUN apk add zip
RUN wget -O comic-sans.zip https://www.wfonts.com/download/data/2014/06/05/comic-sans-ms/comic-sans-ms.zip
RUN unzip comic-sans.zip

# install deps
RUN go get -d -v ./...
# compile binary
RUN go build -o disbot
# run binary
CMD ["./disbot"]
