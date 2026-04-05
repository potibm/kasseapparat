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
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/potibm/kasseapparat/cmd.Version=${VERSION}" -o kasseapparat . 

# ==========================================
# Create the final image
# ==========================================
FROM alpine:3.23 AS runtime
WORKDIR /app

RUN apk update --no-cache && \
    apk add --no-cache ca-certificates bash tzdata && \
    adduser -D -h /app -s /bin/bash appuser

ARG BUILD_DATE
ARG VERSION

RUN mkdir -p /app/data && chown -R appuser:appuser /app

# Copy backend build
COPY --from=backend-build --chown=appuser:appuser /app/backend/kasseapparat ./kasseapparat

RUN chmod +x /app/kasseapparat

USER appuser

VOLUME [ "/app/data" ]

EXPOSE 8080

ENTRYPOINT ["/app/kasseapparat"]
CMD ["serve"]