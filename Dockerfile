FROM golang:1.22

WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# COPY internal/db/db.go ./internal/db
# COPY main.go ./

COPY . ./

RUN go mod tidy
RUN go build -o backend-go

EXPOSE 8080

CMD [ "./backend-go" ]
