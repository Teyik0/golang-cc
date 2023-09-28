FROM golang:latest as builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o main .

FROM golang:latest

WORKDIR /app

COPY --from=builder /app/main /app/.env ./

ENTRYPOINT ["/app/main"]