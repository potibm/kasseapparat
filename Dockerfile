# syntax=docker/dockerfile:1.7

ARG VERSION
ARG BUILD_DATE

# Build the frontend
FROM --platform=$BUILDPLATFORM node:23 AS frontend-build
WORKDIR /app/frontend
COPY frontend/package.json frontend/yarn.lock frontend/.yarnrc.yml ./
RUN --mount=type=cache,target=/app/frontend/.yarn/cache \
    corepack enable && \
    corepack yarn install --immutable
COPY frontend .
RUN corepack yarn vite build --outDir ./build

# Build the backend
FROM --platform=$BUILDPLATFORM golang:1.24-bookworm AS backend-build
WORKDIR /app/backend

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates g++ gcc make
COPY backend .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY --from=frontend-build /app/frontend/build ./cmd/assets
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go build -o kasseapparat ./cmd/main.go && \
    CGO_ENABLED=1 go build -o kasseapparat-tool ./tools/main.go
RUN strip kasseapparat-tool && \
    strip kasseapparat

# Create the final image
FROM --platform=$BUILDPLATFORM debian:bookworm-slim AS runtime

ARG VERSION
ARG BUILD_DATE

WORKDIR /app
VOLUME [ "/app/data" ]
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    useradd -m -s /bin/bash appuser

# Copy backend build
RUN echo "${VERSION}" > /app/VERSION && \
    mkdir -p /app/data && \
    chown -R appuser:appuser /app

# Copy frontend build
COPY --chown=appuser:appuser --from=backend-build /app/backend/kasseapparat ./kasseapparat
COPY --chown=appuser:appuser --from=backend-build /app/backend/kasseapparat-tool ./kasseapparat-tool

USER appuser

# Expose port (adjust based on your application)
EXPOSE 8080

# Command to run the application (adjust based on your application)
CMD ["/app/kasseapparat", "8080"]