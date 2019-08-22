FROM golang:1.12 AS stage
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-w -s" -o /mqtt-emulator -a -installsuffix cgo cmd/*.go

FROM alpine:3.9
COPY config.cfg config.cfg
COPY --from=stage /mqtt-emulator /opt/application
ENTRYPOINT [ "/opt/application" ]