# Build API
FROM golang:1.14 as go_build
WORKDIR /build/

COPY api api
COPY client client
COPY config config
COPY main main
COPY storage storage
COPY go.mod .
RUN GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o tt_api ./main/

FROM scratch

WORKDIR /root/


COPY --from=go_build ["/build/tt_api", "/root/"]
COPY ["config.yml", "schema.sql",  "/root/"]

CMD ["/root/tt_api"]
