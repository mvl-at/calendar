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

FROM scratch
COPY --from=build /go/bin/calserve /calserve
WORKDIR /calendar-data
VOLUME  /calendar-data
EXPOSE  7303
ENTRYPOINT [ "/calserve" ]
