FROM golang:1.16-alpine AS build
ENV GO111MODULE on
ENV GOPROXY="https://goproxy.cn,direct"
WORKDIR /src/
ADD . .
RUN CGO_ENABLED=0 go build -o /bin/demo
FROM scratch
COPY --from=build /bin/demo /bin/demo
COPY conf/sample/conf.yaml .
COPY --from=busybox:1.28 /bin/busybox /bin/busybox
ENTRYPOINT ["/bin/demo","conf.yaml"]