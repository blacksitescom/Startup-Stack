FROM golang:1.16-alpine

WORKDIR /app
#COPY go.mod ./
#COPY go.sum ./
#RUN go mod download
COPY main.go ./
COPY cmd /app/cmd
RUN apk update && apk add bash git
RUN go mod init github.com/gordonianj/seccloud
RUN go mod tidy
RUN go build