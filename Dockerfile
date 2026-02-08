FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum  ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend-go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/backend-go .
EXPOSE 8080
ENTRYPOINT [ "./backend-go" ]
