FROM node:22-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.22-alpine AS backend-build
WORKDIR /app/backend
RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates sqlite-libs tzdata
COPY --from=backend-build /app/server /app/server
COPY --from=frontend-build /app/frontend/dist /app/public
RUN mkdir -p /app/data /app/logs /app/backups
ENV PUBLIC_DIR=/app/public
EXPOSE 8080
CMD ["/app/server"]
