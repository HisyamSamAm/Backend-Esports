# EMBECK API - Team Endpoints Documentation

## Swagger Documentation

Dokumentasi API ini telah dilengkapi dengan Swagger untuk endpoint Team. 

### Cara Menggunakan

1. Jalankan server:
   ```bash
   go run main.go
   ```

2. Buka browser dan akses:
   ```
   http://localhost:1010/swagger/
   ```

### Team Endpoints

API menyediakan endpoint berikut untuk manajemen team:

- **GET** `/api/team` - Mendapatkan semua team
- **GET** `/api/team/{id}` - Mendapatkan team berdasarkan ID
- **POST** `/api/team` - Membuat team baru
- **PUT** `/api/team/{id}` - Update team berdasarkan ID
- **DELETE** `/api/team/{id}` - Hapus team berdasarkan ID

### Model Team

```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Fnatic ONIC",
  "alias": "FNOC",
  "logo_url": "https://example.com/logo.png"
}
```

### Regenerate Documentation

Jika ada perubahan pada komentar Swagger, jalankan:
```bash
swag init
```

### Response Format

Semua endpoint menggunakan format response yang konsisten:

**Success Response:**
```json
{
  "status": 200,
  "message": "success message",
  "data": "response data"
}
```

**Error Response:**
```json
{
  "status": 400/404/500,
  "message": "error message",
  "data": null
}
```
