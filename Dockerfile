FROM golang:1.12-alpine

WORKDIR /application

RUN apk update && apk upgrade && \
    apk add bash git openssh

COPY . ./

RUN go build pososyamba_bot.go

CMD [ "./pososyamba_bot" ]
