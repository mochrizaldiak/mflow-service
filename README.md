# 💸 M-Flow Backend

**M-Flow** adalah REST API backend untuk aplikasi manajemen keuangan pelajar.  
Dibangun menggunakan **Golang**, **Gin**, **GORM**, dan **MySQL**.

---

## ✅ Fitur

- 🔐 Autentikasi menggunakan JWT (register, login)
- 👤 Manajemen user dan saldo
- 📚 Artikel keuangan (dengan pencarian, filter, sorting)
- 💰 Manajemen anggaran (budget)
- 💳 Pencatatan transaksi pemasukan dan pengeluaran
- 📊 Perhitungan **Skor Keuangan** (0–100) per user

---

## 🧰 Teknologi

- Go 1.20+
- Gin Gonic
- GORM (ORM untuk MySQL)
- JWT (JSON Web Token)
- MySQL 5.7+
- Air (opsional - hot reload)

---

## 📦 Struktur Project

```
.
├── controllers      # Handler semua endpoint
├── middleware       # JWT & CORS middleware
├── models           # Struktur database (GORM models)
├── services         # Logika bisnis & DB access
├── main.go          # Entry point server
├── go.mod           # Dependensi Go
└── README.md
```

---

## ⚙️ Setup & Menjalankan

### 1. Clone Repo

```bash
git clone https://github.com/yourusername/mflow-backend.git
cd mflow-backend
```

### 2. Setup Database

Pastikan MySQL berjalan, lalu buat database:

```sql
CREATE DATABASE mflow;
```

### 3. Konfigurasi Database

Edit di `main.go` bagian:

```go
dsn := "root:@tcp(127.0.0.1:3306)/mflow?charset=utf8mb4&parseTime=True&loc=Local"
```

> Ganti `root` dan password sesuai konfigurasi lokalmu.

### 4. Install Dependensi

```bash
go mod tidy
```

### 5. Jalankan Server

```bash
go run main.go
```

> Untuk auto-reload, install [Air](https://github.com/air-verse/air):

```bash
go install github.com/air-verse/air@latest
air
```

---

## 🔐 Autentikasi

Gunakan JWT Bearer token di header:

```http
Authorization: Bearer <token>
```

---

## 🔗 Endpoint Penting

### 🧑‍💼 Users

| Method | Endpoint       | Deskripsi              |
|--------|----------------|------------------------|
| POST   | `/users`       | Register user baru     |
| POST   | `/users/login` | Login, dapatkan token  |
| GET    | `/users/me`    | Profil + skor keuangan |

### 📚 Articles

| Method | Endpoint        | Deskripsi                            |
|--------|-----------------|--------------------------------------|
| GET    | `/articles`     | List artikel (support search/sort)   |

### 💳 Transactions

| Method | Endpoint             | Deskripsi                   |
|--------|----------------------|-----------------------------|
| GET    | `/transactions`      | List transaksi user         |
| POST   | `/transactions`      | Tambah transaksi            |
| GET    | `/transactions/:id`  | Detail transaksi            |
| PUT    | `/transactions/:id`  | Edit transaksi              |
| DELETE | `/transactions/:id`  | Hapus transaksi             |

### 💰 Budgets

| Method | Endpoint       | Deskripsi              |
|--------|----------------|------------------------|
| GET    | `/budgets`     | List anggaran user     |
| POST   | `/budgets`     | Tambah anggaran baru   |

---

## 📊 Skor Keuangan

Endpoint `/users/me` akan otomatis menghitung **skor keuangan (0–100)** berdasarkan:

- Saving rate (pengeluaran vs pemasukan)
- Aktivitas transaksi
- Penggunaan anggaran
- Ragam kategori transaksi

Contoh response:

```json
{
  "id": 1,
  "nama": "Budi",
  "email": "budi@email.com",
  "no_hp": "08123456789",
  "saldo": 120000,
  "skor_keuangan": 82
}
```

---

## 🧪 Tes API

Gunakan Postman / Insomnia / curl / frontend Nuxt kamu.

Contoh:

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/users/me
```

---

## 📄 Lisensi

MIT License © 2025 M-Flow Team

---

## ✨ Kontribusi

Pull request, ide, dan masukan sangat diterima 🙌
