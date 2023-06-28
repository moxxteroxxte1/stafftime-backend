FROM golang:1.19.0

WORKDIR /usr/src/stafftime/backend/

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy

EXPOSE 3000