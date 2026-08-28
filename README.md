# Keuskupan Laboan Bajo - Accounting API

REST API untuk sistem akuntansi keuskupan/paroki menggunakan Go + Echo + GORM + MySQL.

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Echo v4
- **ORM**: GORM v1
- **Database**: MySQL 8.0+ (siap migrasi ke SQLite)
- **Validation**: go-playground/validator v10
- **Logging**: Echo built-in logger

## Project Structure

```
KeuskupanLaboanBajo/
├── main.go                 # Entry point
├── config/
│   └── database.go         # DB connection + AutoMigrate
├── models/
│   └── schema.go           # 8 models (User, Keuskupan, Paroki, Jenis, Akun, Jurnal, DetilJurnal, Pembatasan)
├── dto/
│   ├── request.go          # Request DTOs dengan validation tags
│   └── response.go         # Response DTOs (CreateResponse, ListResponse, ErrorResponse)
├── controllers/
│   ├── user.go
│   ├── keuskupan.go
│   ├── paroki.go
│   ├── jenis.go
│   ├── akun.go
│   ├── jurnal.go           # Complex: UUID, nested detil_jurnal, transactions
│   ├── pembatasan.go
│   └── detil_jurnal.go     # Stub
├── middleware/
│   └── error_handler.go    # Custom error handler + validator
├── routes/
│   └── routes.go           # Route registration
└── pkg/logger/             # Logger wrapper (optional)
```

## Database Schema (8 Tables)

| Table | Primary Key | Relations |
|-------|-------------|-----------|
| `user` | `id` (uint) | - |
| `keuskupan` | `id` (uint) | has many `paroki`, `jurnal` |
| `paroki` | `id` (uint) | belongs to `keuskupan`, has many `detil_jurnal` |
| `jenis` | `id` (uint) | has many `akun` |
| `akun` | `id` (uint) | belongs to `jenis`, has many `detil_jurnal`, `pembatasan` |
| `jurnal` | `id` (UUID string) | belongs to `keuskupan`, has many `detil_jurnal` |
| `detil_jurnal` | `id` (uint) | belongs to `jurnal`, `akun`, `paroki` |
| `pembatasan` | `id` (uint) | belongs to `akun` (nullable) |

**Catatan**: Semua tabel menggunakan `SingularTable: true` (tidak plural).

## Setup

### Prerequisites
- Go 1.21+
- MySQL 8.0+

### Environment Variables (`.env`)
```env
DB_USER=root
DB_PASS=your_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=db_keuskupan
```

### Run
```bash
# Install deps
go mod tidy

# Run (auto-migrate tables on start)
go run main.go

# Server: http://localhost:8080
```

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

### Auth
Belum ada (TODO: JWT/API Key)

### Resources

#### User
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users` | Create user |
| GET | `/users` | List users (pagination, search) |
| GET | `/users/:id` | Get user |
| PUT | `/users/:id` | Update user |
| DELETE | `/users/:id` | Delete user |

**Create User:**
```json
{
  "nama": "John Doe",
  "email": "john@example.com"
}
```

#### Keuskupan
| Method | Endpoint |
|--------|----------|
| POST | `/keuskupan` |
| GET | `/keuskupan` |
| GET | `/keuskupan/:id` |
| PUT | `/keuskupan/:id` |
| DELETE | `/keuskupan/:id` |

**Create:**
```json
{
  "nama": "Keuskupan Agung Jakarta",
  "alamat": "Jl. Katedral No. 7B, Jakarta"
}
```

#### Paroki
| Method | Endpoint |
|--------|----------|
| POST | `/paroki` |
| GET | `/paroki` |
| GET | `/paroki/:id` |
| PUT | `/paroki/:id` |
| DELETE | `/paroki/:id` |

**Create:**
```json
{
  "nama": "Paroki Katedral",
  "alamat": "Jakarta Pusat",
  "keuskupan_id": 1
}
```

#### Jenis (Master data akun)
| Method | Endpoint |
|--------|----------|
| POST | `/jenis` |
| GET | `/jenis` |
| GET | `/jenis/:id` |
| PUT | `/jenis/:id` |
| DELETE | `/jenis/:id` |

**Create:**
```json
{ "nama": "Aset" }
```

#### Akun (Chart of Accounts)
| Method | Endpoint |
|--------|----------|
| POST | `/akun` |
| GET | `/akun` |
| GET | `/akun/:id` |
| PUT | `/akun/:id` |
| DELETE | `/akun/:id` |

**Create:**
```json
{
  "kode": "1101",
  "nama": "Kas Paroki",
  "jenis_id": 1
}
```

#### Jurnal (Transaksi - Double Entry)
| Method | Endpoint |
|--------|----------|
| POST | `/jurnal` |
| GET | `/jurnal` |
| GET | `/jurnal/:id` |
| PUT | `/jurnal/:id` |
| DELETE | `/jurnal/:id` |
| GET | `/jurnal/:id/detil` | List detil jurnal |

**Create (wajib minimal 2 detil: 1 debit + 1 credit):**
```json
{
  "keuskupan_id": 1,
  "tanggal": "28-08-2026",
  "deskripsi": "Persembahan Minggu",
  "no_bukti": "BKT-001",
  "detil_jurnal": [
    {
      "akun_id": 1,
      "paroki_id": 1,
      "debit": 5000000,
      "kredit": null,
      "keterangan": "Kas masuk"
    },
    {
      "akun_id": 3,
      "paroki_id": 1,
      "debit": null,
      "kredit": 5000000,
      "keterangan": "Pendapatan persembahan"
    }
  ]
}
```

**Format tanggal**: `DD-MM-YYYY` atau `YYYY-MM-DD` atau RFC3339

#### Pembatasan
| Method | Endpoint |
|--------|----------|
| POST | `/pembatasan` |
| GET | `/pembatasan` |
| GET | `/pembatasan/:id` |
| PUT | `/pembatasan/:id` |
| DELETE | `/pembatasan/:id` |

**Create:**
```json
{
  "tipe": "MAX_TRANSAKSI_HARIAN",
  "nilai": 10000000,
  "akun_id": 1
}
```

---

## Query Parameters (List Endpoints)

| Param | Default | Max | Description |
|-------|---------|-----|-------------|
| `page` | 1 | - | Halaman |
| `limit` | 10 | 100 | Per halaman |
| `sort` | `id DESC` | - | Urutan (contoh: `nama ASC`) |
| `search` | - | - | Pencarian teks |

**Contoh**: `GET /api/v1/akun?page=2&limit=20&sort=nama ASC&search=Kas`

---

## Response Format

### Success Create (201)
```json
{
  "message": "Data Berhasil di Buat",
  "data": { ... }
}
```

### Success List (200)
```json
{
  "data": [...],
  "page": 1,
  "limit": 10,
  "total": 25,
  "total_pages": 3
}
```

### Success Get/Update/Delete (200)
```json
{ "message": "Data Berhasil di Perbarui" }
{ "message": "Data Berhasil di Hapus" }
```

### Error (400/404/500)
```json
{
  "code": "VALIDATION_ERROR|HTTP_ERROR|INTERNAL_ERROR",
  "message": "Deskripsi error",
  "details": [
    { "field": "email", "message": "Format email tidak valid" }
  ]
}
```

---

## Known Issues / Limitations

1. **UUID & MySQL**: GORM preload nested relations dengan UUID menyebabkan error "Illegal double". Workaround: raw SQL di `jurnal.go` (Create, Get, Update, Delete, List).
2. **FK RESTRICT**: Delete `akun`/`jenis` gagal jika masih direferensikan `detil_jurnal`. Return 400 dengan pesan jelas.
3. **No Auth**: API terbuka, perlu JWT/API Key untuk production.
4. **No Tests**: Belum ada unit/integration test.

---

## Development Notes

### Menambah Controller Baru
1. Buat file di `controllers/`
2. Follow pattern: `NewXxxController(db)`, methods `List`, `Get`, `Create`, `Update`, `Delete`
3. Pakai `ctx.Bind()`, `ctx.Validate()`, return `dto.CreateResponse` / `dto.ListResponse`
4. Daftarkan di `routes/routes.go`

### Migrasi ke SQLite
Ganti di `config/database.go`:
```go
import "gorm.io/driver/sqlite"

db, err := gorm.Open(sqlite.Open("keuskupan.db"), &gorm.Config{
    NamingStrategy: schema.NamingStrategy{SingularTable: true},
})
```
Hapus `.env` dan godotenv.

---

## License

Internal use - Keuskupan Laboan Bajo