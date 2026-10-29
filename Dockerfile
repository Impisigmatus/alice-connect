FROM --platform=$BUILDPLATFORM golang:1.26-trixie AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -mod=vendor -o /out/alice-connect ./cmd/service

FROM debian:trixie-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY autogen/oapi3.yaml /app/autogen/oapi3.yaml
COPY --from=build /out/alice-connect /usr/local/bin/alice-connect
USER 65532:65532
EXPOSE 8000
ENTRYPOINT ["/usr/local/bin/alice-connect"]
