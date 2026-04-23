# Machine Katalog - Website Katalog Mesin Manufaktur (MVP)

Project monorepo untuk katalog mesin manufaktur.

## Struktur Project

- `backend/`: Go (Fiber) + GORM + PostgreSQL
- `frontend-public/`: Next.js 14 + TypeScript + Tailwind
- `frontend-admin/`: React (Vite) + TypeScript + Tailwind

## Prasyarat

- Docker & Docker Compose
- Go 1.22+ (opsional untuk dev lokal)
- Node.js 20+ (opsional untuk dev lokal)

## Cara Menjalankan

### Menggunakan Docker

1. Pastikan Docker sudah berjalan.
2. Jalankan perintah:
   ```bash
   docker-compose up --build
   ```
3. Akses aplikasi:
   - Backend: http://localhost:8080
   - Frontend Public: http://localhost:3000
   - Frontend Admin: http://localhost:3001
   - MinIO Console: http://localhost:9001

## Fitur Milestone 1

- [x] Inisialisasi Monorepo
- [x] Setup Backend Go + Fiber + GORM
- [x] Setup Frontend Public Next.js
- [x] Setup Frontend Admin React Vite
- [x] Docker Compose Configuration
- [x] Database Migration (Admins, Categories, Products, etc.)
