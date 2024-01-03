FROM golang:1.21-alpine as builder

WORKDIR /go/src/app

COPY . .
RUN go mod download && go mod verify
RUN go build -v -o /go/bin/dockge-gitops ./cmd/main.go

FROM alpine:3.19

RUN addgroup -g 1337 dockge-gitops && \
    adduser -G dockge-gitops -D -H -u 1337 app

COPY --from=builder /go/bin/dockge-gitops /app/cmd/

RUN chown -R app:dockge-gitops /app

USER 1337

WORKDIR /app/cmd

CMD ["/app/cmd/dockge-gitops"]