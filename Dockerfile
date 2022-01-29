# START: probes
# START: beginning
FROM golang:1.14-alpine AS build
WORKDIR /go/src/proglog
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/prolog ./cmd/prolog
# END: beginning
# START_HIGHLIGHT
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/go/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/tag/${GRPC_HEALTH_PROBE_VERSION} && \
    chmod +x /go/bin/grpc_health_probe
# END_HIGHLIGHT
# START: beginning

FROM scratch
COPY --from=build /go/bin/prolog /bin/prolog
# END: beginning
# START_HIGHLIGHT
COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe
# END_HIGHLIGHT
# START: beginning
ENTRYPOINT ["/bin/prolog"]
# END: beginning
# END: probes
