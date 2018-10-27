FROM golang:1.11-alpine3.7 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
# ENV GO111MODULE=on

WORKDIR /go/src/github.com/mvl-at/calendar
RUN apk add --no-cache \
    git \
    musl-dev
COPY . /go/src/github.com/mvl-at/calendar
RUN go get ./...
RUN go install -ldflags '-s -w' ./cmd/calserve

# ---

FROM alpine:edge
COPY --from=build /go/bin/calserve /calserve
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories; apk update; apk add --no-cache wkhtmltopdf ghostscript-fonts
WORKDIR /calendar-data
VOLUME  /calendar-data
EXPOSE  7303
ENTRYPOINT [ "/calserve" ]
