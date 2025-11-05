FROM golang:1.25

WORKDIR /go/src/app

COPY . .

EXPOSE 8000

RUN go mod tidy
RUN go build -o main cmd/main.go

CMD [ "./main" ]