FROM golang:1.12-alpine

WORKDIR /application

RUN apk update && apk upgrade && \
    apk add bash git openssh

COPY . ./

RUN go mod download

RUN go build ./cmd/pososyamba_bot/main.go

CMD [ "./main" ]
