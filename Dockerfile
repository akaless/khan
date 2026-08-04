# ═══════════════════════════════════════════════════════════
# خان (Khan) — Dockerfile
# چت سازمانی سبک فارسی — مولتی‌استیج بیلد
# ═══════════════════════════════════════════════════════════

# ─── مرحله ۱: بیلد ───
FROM golang:1.22-alpine AS builder

WORKDIR /build

# کپی و دانلود وابستگی‌ها (کش بهینه)
COPY code/go.mod code/go.sum ./
RUN go mod download

# کپی سورس
COPY code/ .

# بیلد استاتیک
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /khan ./cmd/server

# ─── مرحله ۲: اجرا ───
FROM alpine:3.20

# نصب CA certificates برای TLS
RUN apk --no-cache add ca-certificates tzdata

# کاربر غیر root (امنیت)
RUN addgroup -S khan && adduser -S khan -G khan

WORKDIR /app
COPY --from=builder /khan /app/khan

# پوشه داده
RUN mkdir -p /app/data && chown -R khan:khan /app

USER khan

# پورت خان
EXPOSE 1727

# سلامت
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://localhost:1727/api/settings/info || exit 1

VOLUME ["/app/data"]

ENTRYPOINT ["/app/khan"]
CMD ["-config", "/app/config.json", "-data", "/app/data", "-host", "0.0.0.0", "-port", "1727"]
