FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o exchange ./cmd/main.go
EXPOSE 8088
CMD ["./exchange"]