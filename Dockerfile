FROM golang:1.12-alpine AS pososyamba

WORKDIR /application

RUN apk update && apk upgrade && \
    apk add bash git openssh

COPY . ./

RUN go mod download

RUN go build ./cmd/pososyamba_bot/main.go

CMD [ "./main" ]

FROM alpine:latest
WORKDIR /pososyamba/
COPY --from=pososyamba /application/main pososyamba
CMD ["./pososyamba"]
