FROM golang:1.23-alpine
WORKDIR /nikiax-testing-back
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /nikiax-testing-back cmd/nikiax-testing-back/main.go
RUN chmod +x /nikiax-testing-back
CMD ["./main"]