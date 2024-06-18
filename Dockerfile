# Build the frontend
FROM node:22 AS frontend-build
WORKDIR /app/frontend
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install
COPY frontend .
RUN yarn run build

# Build the backend
FROM golang:1.22-alpine AS backend-build
WORKDIR /app/backend
RUN apk update && apk add --no-cache gcc g++
COPY backend .
RUN go mod download
COPY --from=frontend-build /app/frontend/build ./cmd/assets
RUN CGO_ENABLED=1 go build -o kasseapparat ./cmd/main.go
RUN CGO_ENABLED=1 go build -o kasseapparat-tool ./tools/main.go

# Create the final image
FROM alpine:latest
WORKDIR /app
VOLUME [ "/app/data" ]

# Copy frontend build
COPY --from=backend-build /app/backend/kasseapparat ./kasseapparat
COPY --from=backend-build /app/backend/kasseapparat-tool ./kasseapparat-tool
COPY VERSION VERSION

# Copy backend build
RUN chmod +x /app/kasseapparat && chmod +x /app/kasseapparat-tool 

# Expose port (adjust based on your application)
EXPOSE 8080

# Command to run the application (adjust based on your application)
CMD ["/app/kasseapparat", "8080"]