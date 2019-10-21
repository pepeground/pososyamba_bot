FROM golang:1.13-alpine AS pososyamba

WORKDIR /application

RUN apk update && apk upgrade && \
    apk add bash git openssh

COPY . ./

RUN go mod download

RUN go build ./cmd/pososyamba_bot/main.go

FROM alpine:latest
WORKDIR /application

COPY --from=pososyamba /application /application

CMD [ "./main" ]
