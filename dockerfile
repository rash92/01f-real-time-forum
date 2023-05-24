FROM golang:1.16-alpine as build

LABEL version="1.0"
LABEL description="HMRP Forum"
LABEL authors="HMRP"
LABEL author-usernames="HMRP"

# for bash troubleshooting
RUN apk add --no-cache bash \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download

# build current directory as main
RUN go build -o main .


FROM alpine:3.12

COPY --from=build /app /app
WORKDIR /app

CMD ["/app/main"]