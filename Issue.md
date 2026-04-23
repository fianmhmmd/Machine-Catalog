# Project Planning — Website Katalog Mesin Manufaktur (MVP)

Versi yang lebih ramping dan realistis untuk perusahaan yang baru merintis. Fokus pada **MVP (Minimum Viable Product)** — fitur inti dulu, bisa dikembangkan bertahap nanti.

---

## 🎯 Prinsip Perencanaan

```
✅ Sederhana dulu, sempurnakan nanti
✅ Fokus fitur yang langsung berdampak ke customer
✅ Single developer friendly
✅ Budget server minimal (~$5-10/bulan)
✅ Target launch: 3 minggu
```

---

## 🏷️ Labels

```
setup        → hijau
backend      → biru
frontend     → oranye
admin        → kuning
bug          → merah
enhancement  → cyan (untuk pengembangan masa depan)
```

---

## 📅 MILESTONE 1: Fondasi Project

**Durasi:** 2–3 hari
**Deskripsi:** Setup semua kebutuhan dasar agar bisa langsung mulai develop fitur.

---

### Issue #1 — Setup Repository, Database & Docker

**Labels:** `setup`

**Deskripsi:**
Inisialisasi seluruh project sekaligus — repo, backend Go + Fiber, frontend public Next.js, frontend admin React + Vite, database PostgreSQL, dan Docker Compose untuk development.

**Tasks:**

- [ ] Buat repository GitHub dengan struktur monorepo:
  ```
  katalog-mesin/
  ├── backend/
  ├── frontend-public/
  ├── frontend-admin/
  ├── docker-compose.yml
  ├── .gitignore
  └── README.md
  ```
- [ ] Init backend Go: install Fiber, GORM, JWT, gomail, validator
- [ ] Init frontend public: Next.js 14 + TypeScript + Tailwind + ShadCN/UI
- [ ] Init frontend admin: React + Vite + TypeScript + Tailwind + ShadCN/UI
- [ ] Buat `docker-compose.yml` (api, web, admin, postgres, minio)
- [ ] Buat migration database:
  - `admins` (id, email, password_hash, name)
  - `categories` (id, name, slug)
  - `products` (id, category_id, name, slug, description, specifications JSONB, contact_phone, contact_name, is_published, timestamps)
  - `product_images` (id, product_id, image_url, is_primary, sort_order)
  - `product_analytics` (id, product_id, event_type, visitor_ip, created_at)
  - `inquiries` (id, product_id, customer_name, customer_email, customer_phone, message, is_read, created_at)
- [ ] Seed admin default
- [ ] Verify: `docker-compose up` → semua service jalan, DB termigrate

**Selesai ketika:**
- [ ] Semua service berjalan via Docker
- [ ] Backend merespons `GET /api/v1/health`
- [ ] Frontend public & admin menampilkan halaman default

---

## 📅 MILESTONE 2: Backend API

**Durasi:** 5–6 hari
**Deskripsi:** Seluruh API yang dibutuhkan frontend.

---

### Issue #2 — API Auth & Middleware

**Labels:** `backend`
**Depends on:** #1

**Deskripsi:**
Login admin (JWT), middleware auth, dan CORS.

**Endpoints:**
```
POST /api/v1/auth/login       → return access token
POST /api/v1/auth/refresh     → return new access token
```

**Tasks:**

- [ ] Login: validasi email + password (bcrypt), generate JWT
- [ ] Middleware: validasi JWT di semua route `/admin/*`
- [ ] CORS middleware (izinkan frontend origins)
- [ ] Access token expiry 1 jam, refresh token 7 hari

**Selesai ketika:**
- [ ] Login berhasil dengan credentials benar, gagal dengan credentials salah
- [ ] Endpoint admin return 401 tanpa token

---

### Issue #3 — API CRUD Produk & Kategori + Upload Foto

**Labels:** `backend`
**Depends on:** #2

**Deskripsi:**
CRUD lengkap untuk produk dan kategori, termasuk upload foto ke MinIO.

**Endpoints:**
```
Public:
GET  /api/v1/categories
GET  /api/v1/products              → list (filter category, search, pagination)
GET  /api/v1/products/:slug        → detail + images
GET  /api/v1/products/:slug/related

Admin:
POST   /api/v1/admin/categories
PUT    /api/v1/admin/categories/:id
DELETE /api/v1/admin/categories/:id

GET    /api/v1/admin/products
POST   /api/v1/admin/products
PUT    /api/v1/admin/products/:id
DELETE /api/v1/admin/products/:id

POST   /api/v1/admin/products/:id/images
DELETE /api/v1/admin/products/:id/images/:imageId
```

**Tasks:**

- [ ] CRUD kategori (auto-generate slug)
- [ ] CRUD produk (auto-generate slug, relasi ke kategori)
- [ ] Upload foto ke MinIO (validasi: max 5MB, jpg/png/webp)
- [ ] Multiple images per produk, set primary image
- [ ] Related products: ambil produk dari kategori yang sama
- [ ] List produk: filter by category, search by nama, pagination
- [ ] Public endpoint hanya tampilkan `is_published = true`
- [ ] Hapus produk → hapus juga foto di MinIO

**Selesai ketika:**
- [ ] Semua endpoint berfungsi (test via Postman/Thunder Client)
- [ ] Foto tersimpan di MinIO dan URL bisa diakses

---

### Issue #4 — API Analytics & Email Inquiry

**Labels:** `backend`
**Depends on:** #3

**Deskripsi:**
Tracking views/clicks dan sistem contact form yang kirim email ke admin.

**Endpoints:**
```
Public:
POST /api/v1/products/:id/view
POST /api/v1/products/:id/click
POST /api/v1/inquiry

Admin:
GET /api/v1/admin/analytics/overview    → total views, clicks, top products
GET /api/v1/admin/inquiries
PUT /api/v1/admin/inquiries/:id/read
```

**Tasks:**

- [ ] Track view & click (simpan ke DB: product_id, event_type, IP, timestamp)
- [ ] Simple dedup: 1 view per IP per produk per 24 jam (cek di DB, tanpa Redis)
- [ ] Inquiry: validasi input → simpan ke DB → kirim email ke admin via SMTP
- [ ] Email template HTML sederhana (nama, email, phone, pesan, produk)
- [ ] Rate limit inquiry: max 3 per IP per jam (simple in-memory atau DB check)
- [ ] Analytics overview: total views & clicks (7 hari terakhir), top 5 produk
- [ ] List inquiry: paginated, filter read/unread

**Selesai ketika:**
- [ ] View/click tercatat di DB
- [ ] Submit inquiry → email terkirim ke admin
- [ ] Admin bisa lihat overview analytics dan list inquiry

---

## 📅 MILESTONE 3: Frontend Public (Next.js)

**Durasi:** 5–6 hari
**Deskripsi:** Website yang dilihat customer.

---

### Issue #5 — Layout, Home & Halaman Katalog

**Labels:** `frontend`
**Depends on:** #3

**Deskripsi:**
Buat layout utama, halaman home, dan halaman list katalog produk.

**Tasks:**

- [ ] Layout: Header (logo, nav: Home, Katalog) + Footer (info perusahaan, copyright)
- [ ] Responsive: mobile-first, hamburger menu
- [ ] **Halaman Home** (`/`):
  - [ ] Hero section (gambar + tagline + tombol "Lihat Katalog")
  - [ ] Grid kategori produk
  - [ ] Grid 6 produk terbaru
  - [ ] Gunakan ISR (revalidate 5 menit)
- [ ] **Halaman Katalog** (`/katalog`):
  - [ ] Grid produk (responsive: 1/2/3 kolom)
  - [ ] Filter kategori
  - [ ] Search bar (debounced)
  - [ ] Pagination
  - [ ] Skeleton loading
  - [ ] Empty state jika tidak ada hasil
- [ ] SEO: meta title & description per halaman

**Selesai ketika:**
- [ ] Home menampilkan kategori dan produk terbaru
- [ ] Katalog bisa difilter, dicari, dan dipaginasi
- [ ] Responsive di mobile dan desktop
- [ ] Lighthouse Performance > 85

---

### Issue #6 — Halaman Detail Produk

**Labels:** `frontend`
**Depends on:** #5, #4

**Deskripsi:**
Halaman detail produk lengkap dengan semua fitur yang diminta client.

**Tasks:**

- [ ] Halaman `/katalog/[slug]` dengan SSG + ISR
- [ ] **Gallery foto**: main image + thumbnail, klik untuk ganti
- [ ] **Info produk**: nama, kategori, deskripsi, tabel spesifikasi
- [ ] **Tombol WhatsApp**:
  - [ ] Tampilkan nama contact person
  - [ ] Klik → buka `wa.me/{nomor}?text={pesan otomatis}`
  - [ ] Track click ke API
- [ ] **Form Inquiry** (kirim email ke admin):
  - [ ] Fields: Nama, Email, No. HP, Pesan
  - [ ] Validasi (React Hook Form + Zod)
  - [ ] Submit → API → toast success/error
- [ ] **Produk Terkait**: grid 4 produk dari kategori yang sama
- [ ] Track page view saat halaman dibuka
- [ ] SEO: dynamic meta tags + og:image dari foto utama

**Selesai ketika:**
- [ ] Semua section tampil dengan benar
- [ ] WhatsApp button berfungsi (buka chat WA)
- [ ] Form inquiry berhasil kirim email
- [ ] Produk terkait muncul
- [ ] Analytics (view & click) tercatat

---

## 📅 MILESTONE 4: Admin Dashboard

**Durasi:** 3–4 hari
**Deskripsi:** Panel admin untuk kelola produk dan lihat statistik.

---

### Issue #7 — Login & Layout Admin

**Labels:** `admin`
**Depends on:** #2

**Deskripsi:**
Halaman login dan kerangka dasar admin dashboard.

**Tasks:**

- [ ] Halaman Login (email + password)
- [ ] Auth state: simpan token, auto redirect jika belum login
- [ ] Layout dashboard: sidebar (Dashboard, Produk, Kategori, Inquiry) + header (nama admin, logout)
- [ ] Protected routes

**Selesai ketika:**
- [ ] Login/logout berfungsi
- [ ] Sidebar navigasi berfungsi
- [ ] Halaman admin tidak bisa diakses tanpa login

---

### Issue #8 — Admin: Kelola Produk & Kategori

**Labels:** `admin`
**Depends on:** #7, #3

**Deskripsi:**
Interface CRUD produk dan kategori.

**Tasks:**

- [ ] **List Produk**: tabel (nama, kategori, status, tanggal) + search + aksi (edit, hapus, toggle publish)
- [ ] **Form Produk** (create & edit):
  - [ ] Input: nama, kategori (dropdown), deskripsi, spesifikasi (dynamic key-value), contact person, no WA
  - [ ] Upload foto: drag & drop, preview, set primary, hapus
  - [ ] Toggle published
- [ ] **Kategori**: list + tambah/edit/hapus (bisa pakai modal sederhana)
- [ ] Konfirmasi sebelum hapus
- [ ] Toast notification success/error

**Selesai ketika:**
- [ ] Bisa tambah, edit, hapus produk dan kategori
- [ ] Upload foto berfungsi
- [ ] Data yang diubah langsung terefleksi di website public

---

### Issue #9 — Admin: Dashboard Analytics & Inquiry

**Labels:** `admin`
**Depends on:** #7, #4

**Deskripsi:**
Dashboard statistik sederhana dan list inquiry masuk.

**Tasks:**

- [ ] **Dashboard**:
  - [ ] 3 kartu: Total Views (7 hari), Total Clicks (7 hari), Inquiry Bulan Ini
  - [ ] Tabel top 5 produk by views
- [ ] **Halaman Inquiry**:
  - [ ] Tabel inquiry (tanggal, nama, email, produk, status)
  - [ ] Klik row → lihat detail pesan
  - [ ] Mark as read
  - [ ] Badge unread count di sidebar

**Selesai ketika:**
- [ ] Dashboard menampilkan data analytics yang benar
- [ ] List inquiry berfungsi dengan mark as read

---

## 📅 MILESTONE 5: Deploy & Go-Live

**Durasi:** 2 hari
**Deskripsi:** Deploy ke server dan pastikan semua berfungsi.

---

### Issue #10 — Deployment & Final Check

**Labels:** `setup`
**Depends on:** #6, #8, #9

**Deskripsi:**
Deploy ke VPS, setup domain, SSL, dan final testing.

**Tasks:**

- [ ] **Server**: Sewa VPS (2 vCPU, 2-4GB RAM) — contoh: DigitalOcean $12/bln
- [ ] **Setup**: Docker, firewall (UFW), SSH key only
- [ ] **Domain**: Setup DNS, arahkan ke server
- [ ] **Cloudflare**: DNS, SSL (Full Strict), basic caching
- [ ] **Nginx**: reverse proxy, HTTPS redirect, gzip compression
- [ ] **Deploy**: `docker-compose up -d` di production
- [ ] **Backup**: setup cron job `pg_dump` harian
- [ ] **Monitoring**: daftar UptimeRobot (gratis) untuk uptime alert

**Final Checklist:**
- [ ] Login admin berfungsi
- [ ] CRUD produk + upload foto berfungsi
- [ ] Website publik menampilkan produk
- [ ] Detail produk: gallery, WA button, form inquiry semua berfungsi
- [ ] Email inquiry terkirim ke admin
- [ ] Dashboard analytics menampilkan data
- [ ] HTTPS aktif
- [ ] Test di Chrome & Safari (desktop + mobile)
- [ ] Serahkan credentials & panduan singkat ke client

**Selesai ketika:**
- [ ] Website live dan bisa diakses publik via domain
- [ ] Semua fitur berfungsi di production
- [ ] Client sudah bisa login dan kelola produk sendiri

---
