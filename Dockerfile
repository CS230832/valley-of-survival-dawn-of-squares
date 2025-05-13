FROM node:24-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend ./
RUN npm run build

FROM golang:1.24 AS backend-builder
WORKDIR /app
COPY backend/go.* ./
RUN go mod download
COPY backend/ ./
RUN go build -o server .

FROM postgres:17
WORKDIR /app
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=backend-builder /app/server ./
COPY --from=frontend-builder /app/frontend/dist /app/static
COPY scripts/init-db.sql /docker-entrypoint-initdb.d/
COPY scripts/start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 8080 8080
CMD ["bash", "/start.sh"]