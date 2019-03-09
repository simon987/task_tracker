# Build
FROM golang:1.11.5 as go_build
WORKDIR /go/src/github.com/simon987/task_tracker/

COPY . .
RUN go get ./main/ && GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o tt_api ./main/

# Execute in alpine
FROM alpine:3.9.2
WORKDIR /root

COPY --from=go_build ["/go/src/github.com/simon987/task_tracker/tt_api",\
                     "/go/src/github.com/simon987/task_tracker/schema.sql",\
                     "/go/src/github.com/simon987/task_tracker/config.yml",\
                      "./"]
CMD ["./tt_api"]
