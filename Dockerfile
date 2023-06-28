FROM golang:1.19.0

WORKDIR /usr/src/stafftime/backend/

COPY . .
RUN go mod tidy

EXPOSE 3000