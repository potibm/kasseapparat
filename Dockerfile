ARG VERSION
ARG BUILD_DATE

# Build the frontend
FROM --platform=$BUILDPLATFORM node:24 AS frontend-build
WORKDIR /app/frontend
COPY frontend/package.json frontend/yarn.lock frontend/.yarnrc.yml ./
RUN corepack enable && \
    corepack yarn install --immutable
COPY frontend .
RUN corepack yarn vite build --outDir ./build

# Build the backend
FROM --platform=$BUILDPLATFORM golang:1.24-bookworm AS backend-build
WORKDIR /app/backend
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates g++ gcc make && \
    apt-get clean && rm -rf /var/lib/apt/lists/*
COPY backend .
RUN go mod download
COPY --from=frontend-build /app/frontend/build ./cmd/assets
RUN CGO_ENABLED=1 go build -o kasseapparat ./cmd/main.go && \
    CGO_ENABLED=1 go build -o kasseapparat-tool ./tools/main.go

# Create the final image
FROM --platform=$BUILDPLATFORM debian:bookworm-slim AS runtime

ARG VERSION
ARG BUILD_DATE

WORKDIR /app
VOLUME [ "/app/data" ]
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    useradd -m -s /bin/bash appuser && \
    rm -rf /var/lib/apt/lists/* 

# Copy frontend build
COPY --from=backend-build /app/backend/kasseapparat ./kasseapparat
COPY --from=backend-build /app/backend/kasseapparat-tool ./kasseapparat-tool

# Copy backend build
RUN echo "${VERSION}" > /app/VERSION && \
    mkdir -p /app/data && \
    chown -R appuser:appuser /app/data && \
    chown -R appuser:appuser /app && \
    chmod +x /app/kasseapparat && \
    chmod +x /app/kasseapparat-tool 

USER appuser

# Expose port (adjust based on your application)
EXPOSE 8080

# Command to run the application (adjust based on your application)
CMD ["/app/kasseapparat", "8080"]