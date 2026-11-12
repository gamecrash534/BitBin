FROM golang:1.26-alpine AS build
WORKDIR /
COPY . ./
RUN go build -ldflags "-w" /

FROM busybox
COPY --from=build /bitbin /

ENV BITBIN_PORT=8080 \
    BITBIN_DIRECTORY="/data/bitbin/" \
    BITBIN_LINK="https://git.gamecrash.dev/game.crash/BitBin" \
    BITBIN_LINK_TITLE="source" \
    BITBIN_PASTE_TTL="43200"

EXPOSE 8080

ENTRYPOINT ["/bitbin"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --spider http://localhost:${BITBIN_PORT} || exit 1
