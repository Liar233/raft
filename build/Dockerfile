FROM golang:1.12 as build

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o node ./cmd/node/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/node /app/node
COPY --from=build /app/configs/config.yaml ./configs/config.yaml

RUN chmod +x /app/node

EXPOSE 80

ENTRYPOINT ["/app/node",  "/app/configs/config.yaml"]
