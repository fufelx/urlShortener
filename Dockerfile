FROM golang:1.23

WORKDIR /urlShortener
COPY . .

RUN go mod tidy
RUN go build -o main /urlShortener/cmd/main/main.go

CMD ["./main"]