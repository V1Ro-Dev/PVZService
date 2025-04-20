FROM golang:1.24 AS builder

WORKDIR /pvz_app

COPY backend/go.mod backend/go.sum backend/
RUN cd backend && go mod download

COPY backend /pvz_app/backend
COPY deploy/config /pvz_app/deploy/config

WORKDIR /pvz_app/backend

# ✅ Кросс-компиляция под Linux/amd64
RUN go build -o main main.go
RUN chmod +x main

FROM debian:bookworm-slim

WORKDIR /pvz_app/backend

COPY --from=builder /pvz_app/backend/main .
COPY --from=builder /pvz_app/deploy /pvz_app/deploy
COPY deploy/.env /etc/environment

CMD ["./main"]
