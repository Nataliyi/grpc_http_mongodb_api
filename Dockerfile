# syntax = docker/dockerfile:1.0-experimental

FROM golang:1.17 as build

ARG GRPC_HEALTH_PROBE_VERSION=0.4.5
RUN wget -q -O /bin/grpc_health_probe "https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64" && \
    chmod +x /bin/grpc_health_probe

WORKDIR /src

COPY . /src/
RUN go mod download
RUN go generate ./...
RUN CGO_ENABLED=0 go build -o /bin/grpc-api -ldflags="-s -w" -trimpath .

FROM gcr.io/distroless/base:latest

ENV GOTRACEBACK=single
ENV PORT=8080
EXPOSE 8080

COPY --from=build /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=build /bin/grpc-api /bin/grpc-api

ENTRYPOINT ["/bin/grpc-api", "-log", "debug"]