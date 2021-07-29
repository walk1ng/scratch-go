FROM golang:latest

EXPOSE 2345

RUN GOPROXY=https://goproxy.cn,direct go get github.com/go-delve/delve/cmd/dlv

WORKDIR /go/src/app

COPY . /go/src/app

ENTRYPOINT ["dlv", "debug", "--headless", "--listen=:2345", "--log", "--api-version=2"]