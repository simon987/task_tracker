# Build API
FROM golang:1.11.5 as go_build
WORKDIR /go/src/github.com/simon987/task_tracker/

COPY .git .git
COPY api api
COPY client client
COPY config config
COPY main main
COPY storage storage
RUN go get ./main/ && GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o tt_api ./main/

# Build Web
FROM node:10-alpine as npm_build
COPY ./web/ ./
RUN cd ./angular/ && npm install
RUN cd ./angular/ && ./node_modules/@angular/cli/bin/ng build --prod --optimization --output-path "/webroot"

FROM nginx:alpine
WORKDIR /root

COPY nginx.conf schema.sql config.yml ./
COPY --from=go_build ["/go/src/github.com/simon987/task_tracker/tt_api", "./"]
COPY --from=npm_build ["/webroot", "/webroot"]

CMD ["sh", "-c", "nginx -c /root/nginx.conf && /root/tt_api"]
