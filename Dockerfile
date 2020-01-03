# Build API
FROM golang:1.13 as go_build
WORKDIR /go/src/github.com/simon987/task_tracker/

COPY .git .git
COPY api api
COPY client client
COPY config config
COPY main main
COPY storage storage
RUN go get ./main/ && GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o tt_api ./main/

FROM scratch

WORKDIR /root/

COPY --from=go_build ["/go/src/github.com/simon987/task_tracker/tt_api", "/root/"]
COPY ["config.yml", "schema.sql",  "/root/"]

CMD ["/root/tt_api"]
