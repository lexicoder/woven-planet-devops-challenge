FROM golang:1.21-bullseye as base

WORKDIR $GOPATH/src/storage-server/

COPY src/ .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /storage-server main.go

FROM gcr.io/distroless/static-debian11

COPY --from=base /storage-server .

ENV SERVER_PORT 5000
ENV STORAGE_PATH /storage/data

CMD ["./storage-server"]
