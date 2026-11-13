
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

RUN git clone https://codeberg.org/emersion/hydroxide.git .

RUN go build -trimpath -ldflags="-s -w" \
    -o /hydroxide .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S hydroxide \
    && adduser -S -G hydroxide \
       -h /home/hydroxide hydroxide \
    && mkdir -p /home/hydroxide/.config/hydroxide \
    && chown -R hydroxide:hydroxide /home/hydroxide

COPY --from=builder /hydroxide /usr/local/bin/hydroxide

USER hydroxide

ENV HOME=/home/hydroxide

ENTRYPOINT ["hydroxide"]

CMD ["-imap-host", "0.0.0.0", "-smtp-host", "0.0.0.0", "-disable-carddav", "serve"]
