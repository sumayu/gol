FROM golang:1.21-alpine AS builder 
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download 
COPY ./src /app/src
COPY ./configYML /app/configYML
COPY ./src/cmd/server-starter/.env /app/.env
RUN go build -o server ./src/cmd/server-starter
FROM alpine:latest 
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
COPY --from=builder /app/configYML app/configYML
CMD [ "./server" ]