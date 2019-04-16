FROM golang:1.12-alpine

WORKDIR /application

RUN apk update && apk upgrade && \
    apk add bash git openssh

# RUN go get -u github.com/influxdata/influxdb1-client/v2
# RUN go get github.com/influxdata/influxdb1-client/v2
# RUN go get -u github.com/influxdata/influxdb1-client
# RUN go get github.com/influxdata/influxdb1-client

COPY . ./
# COPY ./pososyamba_bot.go ./pososyamba_bot.go
# COPY ./go.sum ./go.sum
# COPY ./go.mod ./go.mod

RUN go build pososyamba_bot.go

CMD [ "./pososyamba_bot" ]
