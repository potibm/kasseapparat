# ==========================================
# Build the frontend
# ==========================================
FROM --platform=$BUILDPLATFORM node:25 AS frontend-build
WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend .
RUN npm run build -- --outDir ./build

# ==========================================
# Build the backend
# ==========================================
FROM golang:1.26-bookworm AS backend-build
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .
COPY --from=frontend-build /app/frontend/build ./cmd/assets

ARG VERSION
RUN CGO_ENABLED=1 go build -ldflags "-X main.version=${VERSION}" -o kasseapparat ./cmd/main.go && \
    CGO_ENABLED=1 go build -ldflags "-X main.version=${VERSION}" -o kasseapparat-tool ./tools/main.go

# ==========================================
# Create the final image
# ==========================================
FROM debian:bookworm-slim AS runtime
WORKDIR /app

RUN apt-get update -o Acquire::http::No-Cache=True && \
    apt-get install -y --no-install-recommends ca-certificates && \
    useradd -m -s /bin/bash appuser && \
    rm -rf /var/lib/apt/lists/* 

ARG BUILD_DATE
ARG VERSION

RUN mkdir -p /app/data && chown -R appuser:appuser /app

# Copy backend build
COPY --from=backend-build --chown=appuser:appuser /app/backend/kasseapparat ./kasseapparat
COPY --from=backend-build --chown=appuser:appuser /app/backend/kasseapparat-tool ./kasseapparat-tool

RUN chmod +x /app/kasseapparat && \
    chmod +x /app/kasseapparat-tool 

USER appuser

VOLUME [ "/app/data" ]

EXPOSE 8080

CMD ["/app/kasseapparat", "--port", "8080"]