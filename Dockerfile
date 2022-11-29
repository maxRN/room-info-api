FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY .env ./
COPY *.go ./

RUN go build -o /room-info

EXPOSE 8080

CMD [ "/room-info" ]
