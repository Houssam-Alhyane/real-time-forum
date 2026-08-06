# ---------- Build stage ----------
# Pinned to the bookworm variant so the builder and the runtime stage share
# the same glibc version (the CGO sqlite3 binary links against glibc).
FROM golang:1.25-bookworm AS builder

WORKDIR /app

# Cache dependencies first (only invalidated when go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY main.go ./
COPY backend ./backend
COPY frontend ./frontend

# mattn/go-sqlite3 requires CGO, so we build with it enabled
RUN CGO_ENABLED=1 go build -o /forum .

# ---------- Runtime stage ----------
# glibc-based image: the CGO binary links against glibc at runtime
FROM debian:bookworm-slim

WORKDIR /app

# The app reads ./frontend and ./backend/database/schema.sql at startup
# (relative paths), so both must be present in the image.
COPY --from=builder /forum ./forum
COPY --from=builder /app/frontend ./frontend
COPY --from=builder /app/backend ./backend

EXPOSE 8082

CMD ["./forum"]
