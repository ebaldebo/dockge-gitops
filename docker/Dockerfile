FROM golang:1.21-alpine as builder

WORKDIR /go/src/app

COPY . .
RUN go mod download && go mod verify
RUN go build -v -o /go/bin/dockge-gitops ./cmd/main.go

FROM alpine:3.19

COPY --from=builder /go/bin/dockge-gitops /app/cmd/

WORKDIR /app/cmd

CMD ["/app/cmd/dockge-gitops"]