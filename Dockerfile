FROM golang:1.17-alpine as builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stresstest .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/stresstest .

CMD ["./stresstest"]