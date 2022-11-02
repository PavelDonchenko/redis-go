FROM golang:1.19

RUN mkdir /app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ../.. .
RUN go build -o /main

CMD ["/main"]