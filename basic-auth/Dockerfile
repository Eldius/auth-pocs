FROM golang:1.22.3-alpine3.19 AS builder

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot
ENV GOPRIVATE="bitbucket.org"
ARG SSH_PRIVATE_KEY
WORKDIR /app

COPY ./ /app

RUN apk update && \
    apk add --no-cache git ca-certificates tzdata openssh make && \
    CGO_ENABLED=0 go build -a -v -o app -a -ldflags '-extldflags "-static"' ./cmd/main.go

FROM gcr.io/distroless/static
COPY --from=builder /etc/passwd /etc/passwd
USER nonroot
WORKDIR /app
COPY --from=builder /app/app app
EXPOSE 8080
ENTRYPOINT ["/app/app"]
