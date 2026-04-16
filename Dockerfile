# ============================================================================
# KAVLING BACKEND - DOCKERFILE
# ============================================================================
# 
# File ini digunakan untuk membungkus aplikasi Go ke dalam Docker container.
# Saat deploy ke Railway, file ini akan dibaca otomatis untuk build image.
#
# Dokumentasi: https://docs.docker.com/engine/reference/builder/
#
# ============================================================================

# ─────────────────────────────────────────────────────────────────────────
# STAGE 1: BUILD STAGE (Compile Go Application)
# ─────────────────────────────────────────────────────────────────────────
#
# Fungsi: Mengompilasi source code Go menjadi binary executable.
#         Stage ini menggunakan golang image yang besar karena butuh compiler.
#         Tapi hasilnya hanya binary kecil yang di-copy ke stage 2.
#
# Keuntungan: Image final jadi lebih kecil (hanya ~15MB vs ~800MB jika pakai golang image)
#

FROM golang:1.26.1-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# Copy go.mod dan go.sum (dependency files)
COPY go.mod go.sum ./

# Download semua dependencies
# Layer ini di-cache, jadi cepat kalau dependency tidak berubah
RUN go mod download

# Copy seluruh source code
COPY . .

# Compile Go application menjadi binary executable
# Penjelasan flags:
#   CGO_ENABLED=0  → Disable C compiler (Go native, no system libs needed)
#   GOOS=linux     → Target OS adalah Linux (Railway/Railway pakai Linux)
#   -o server      → Output binary name = "server"
#   ./cmd/server/main.go → Lokasi main.go
#
# Hasil: binary file bernama "server" (± 30MB)
#
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go


# ─────────────────────────────────────────────────────────────────────────
# STAGE 2: RUNTIME STAGE (Minimal Image)
# ─────────────────────────────────────────────────────────────────────────
#
# Fungsi: Menjalankan binary yang sudah dikompilasi di stage 1.
#         Menggunakan alpine:latest (image Linux minimal, hanya ~5MB).
#         Ini adalah image final yang akan di-deploy.
#
# Alur:
#   1. Download binary dari stage 1 ("server")
#   2. Tambah CA certificates (untuk HTTPS)
#   3. Set working directory
#   4. Copy binary & migration files
#   5. Expose port 8080
#   6. Run application
#

FROM alpine:latest

# Install CA certificates (untuk koneksi HTTPS ke external services)
# Contoh: Database, Email, WhatsApp API, dll
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary "server" dari builder stage
# Format: COPY --from=<stage_name> <source> <destination>
# Di container ini, akan tersimpan di /root/server
COPY --from=builder /app/server .

# Copy migration files (untuk setup database)
# Structure:
#   /root/migrations/
#   ├── 001_init.sql
#   └── 002_denah_kavling.sql
COPY --from=builder /app/migrations ./migrations

# Copy .env file (optional, bisa juga pakai environment variable)
# Kalau pakai Railway, lebih baik set variable langsung di dashboard
COPY --from=builder /app/.env ./

# Expose port 8080
# Catatan: Ini hanya dokumentasi, tidak benar-benar membuka port
# Port sebenarnya dibuka via environment variable PORT=8080
EXPOSE 8080

# Health check (optional, bagus untuk monitoring)
# Railway akan menjalankan ini untuk cek aplikasi masih hidup
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD sh -c 'wget --spider -q http://localhost:8080/health || exit 1'

# Start application
# Container akan menjalankan binary "server"
# Environment variables akan di-inject oleh Railway
CMD ["./server"]


# ============================================================================
# PENJELASAN MULTI-STAGE BUILD
# ============================================================================
#
# Kenapa 2 stage? → Untuk membuat image final lebih kecil!
#
# Tanpa multi-stage:
#   - Image: golang:1.26.1-alpine (± 800MB)
#   - Dikompres: masih ~500MB+
#   - Lambat upload ke Railway
#
# Dengan multi-stage:
#   - Stage 1 (builder): golang image (800MB) → compile
#   - Stage 2 (final): alpine image (5MB)
#   - Hanya copy binary dari stage 1 ke stage 2
#   - Final image: ~50-100MB total (kecil!)
#
# Benefit:
#   ✓ Deploy lebih cepat (bandwith lebih sedikit)
#   ✓ Pull image lebih cepat
#   ✓ Start container lebih cepat
#   ✓ Server cost lebih murah (storage lebih kecil)
#
# ============================================================================

# ============================================================================
# CARA MENGGUNAKAN DOCKERFILE SECARA LOCAL
# ============================================================================
#
# 1. Build image:
#    docker build -t kavling-backend:latest .
#
# 2. Run container (development):
#    docker run -p 8080:8080 \
#      -e DATABASE_URL="postgresql://user:pass@localhost/kavling" \
#      -e JWT_SECRET="your-secret-here" \
#      kavling-backend:latest
#
# 3. Run container (production):
#    docker run -d \
#      --name kavling-backend \
#      -p 8080:8080 \
#      -e DATABASE_URL="postgresql://user:pass@prod-server/kavling" \
#      -e JWT_SECRET="your-secret-here" \
#      -e APP_ENV="production" \
#      kavling-backend:latest
#
# 4. Check logs:
#    docker logs kavling-backend
#
# 5. Stop container:
#    docker stop kavling-backend
#
# ============================================================================

# ============================================================================
# CARA DEPLOY KE RAILWAY
# ============================================================================
#
# Railway secara otomatis akan:
#   1. Detect file Dockerfile
#   2. Build image menggunakan Dockerfile ini
#   3. Deploy container ke cloud
#   4. Expose endpoint ke public URL
#
# Environment variables di-set di Railway dashboard:
#   Settings → Environment Variables → Add Variable
#
# Railway akan inject ke container sebagai env vars,
# yang kemudian dibaca oleh Go app via os.Getenv()
#
# ============================================================================
