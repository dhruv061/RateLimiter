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
RUN mkdir -p /app/data /app/logs /app/backups
EXPOSE 8080
CMD ["/app/server"]
