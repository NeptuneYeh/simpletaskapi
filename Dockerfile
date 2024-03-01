FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o simpletask ./cmd/main.go

FROM scratch
COPY --from=builder /app/simpletask .
EXPOSE 8080

CMD ["./simpletask"]