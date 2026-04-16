# Issue: Implementasi Denah Kavling berbasis SVG

## Ringkasan

Mengubah sistem kavling dari flat list menjadi hierarki **Denah Kavling → Kavling (tanah)**.
User menginput kode SVG melalui textarea, sistem mem-parsing `<path>` yang memiliki atribut `id`,
lalu otomatis membuat entri kavling untuk setiap path tersebut.

---

## Kondisi Saat Ini

- Kavling disimpan sebagai flat list tanpa parent (denah).
- SVG diupload sebagai file (`POST /kavling/upload-svg` dengan `multipart/form-data`).
- Tipe `Kavling` memiliki field: `kode_kavling`, `ukuran`, `luas_tanah`, `hrg_meter`, `hrg_jual`, `status`, `polygon_points`, `path_d`, `transform`.
- Peta interaktif (`PetaInteraktif.tsx`) me-render semua kavling dalam satu SVG viewBox hardcoded `0 0 800 600`.

## Kondisi Yang Diinginkan

- Kavling dikelompokkan di bawah **Denah Kavling** (satu denah = satu SVG lengkap).
- Satu project bisa memiliki **multiple denah kavling**.
- User menambah denah kavling dengan **paste kode SVG ke textarea** (bukan upload file).
- Sistem **mem-parse SVG**, mengekstrak semua `<path id="...">`, lalu **otomatis membuat kavling** per path.
- Setiap kavling memiliki field tambahan: `panjang_kanan`, `panjang_kiri`, `lebar_depan`, `lebar_belakang`, `harga_per_meter`, `harga_jual_cash`.
- Perubahan `status` kavling otomatis mengubah warna `fill` pada path di peta SVG.

---

## Struktur Data Baru

### Tabel: `denah_kavling`

| Field        | Tipe     | Keterangan                                    |
|-------------|----------|-----------------------------------------------|
| id          | int (PK) | Auto increment                                |
| nama        | string   | Nama denah, misal "Blok A"                   |
| svg_content | text     | Kode SVG lengkap yang diinput user            |
| viewbox     | string   | Diambil dari atribut `viewBox` SVG            |
| created_at  | datetime |                                               |
| updated_at  | datetime |                                               |

### Tabel: `kavling` (diperbarui)

| Field            | Tipe     | Keterangan                                                             |
|-----------------|----------|------------------------------------------------------------------------|
| id              | int (PK) | Auto increment                                                         |
| denah_kavling_id| int (FK) | Relasi ke `denah_kavling`                                              |
| kode_kavling    | string   | Diambil dari atribut `id` pada `<path>`, misal "A-1", "A-2"           |
| kode_map        | string   | Diambil dari atribut `d` pada `<path>`, misal "M133.008 158.064L..."  |
| panjang_kanan   | float    | Diisi manual oleh user setelah kavling dibuat                          |
| panjang_kiri    | float    | Diisi manual oleh user setelah kavling dibuat                          |
| lebar_depan     | float    | Diisi manual oleh user setelah kavling dibuat                          |
| lebar_belakang  | float    | Diisi manual oleh user setelah kavling dibuat                          |
| luas_tanah      | float    | Diisi manual atau dihitung dari dimensi                                |
| harga_per_meter | float    | Harga per meter persegi                                                |
| harga_jual_cash | float    | Harga jual cash                                                        |
| status          | int      | 0=Kosong, 1=Booking, 2=Cash, 3=Kredit (mengubah warna fill di SVG)    |
| created_at      | datetime |                                                                        |
| updated_at      | datetime |                                                                        |

> **Catatan:** Field lama `ukuran`, `hrg_meter`, `hrg_jual`, `polygon_points`, `path_d`, `transform` dihapus/diganti.

---

## Tahapan Implementasi

### Tahap 1: Backend - Migrasi Database & Model

**File yang diubah/dibuat:** model, migration, schema

1. Buat migration untuk tabel `denah_kavling` dengan kolom: `id`, `nama`, `svg_content` (TEXT), `viewbox`, `created_at`, `updated_at`.
2. Buat migration untuk mengubah tabel `kavling`:
   - Tambah kolom: `denah_kavling_id` (FK), `kode_map`, `panjang_kanan`, `panjang_kiri`, `lebar_depan`, `lebar_belakang`, `harga_per_meter`, `harga_jual_cash`.
   - Hapus kolom yang tidak dipakai: `ukuran`, `polygon_points`, `path_d`, `transform`.
   - Rename `hrg_meter` → `harga_per_meter`, `hrg_jual` → `harga_jual_cash` (atau buat kolom baru dan migrasi data).
3. Buat model `DenahKavling` dengan relasi `hasMany` ke `Kavling`.
4. Update model `Kavling` dengan relasi `belongsTo` ke `DenahKavling`.

### Tahap 2: Backend - SVG Parser & API Denah Kavling

**File yang diubah/dibuat:** controller, service/helper, routes

1. Buat helper/service **SVG Parser** yang:
   - Menerima string SVG.
   - Validasi format SVG (cek tag `<svg>`, ambil `viewBox`).
   - Extract semua `<path>` yang memiliki atribut `id`.
   - Return array: `[{ kode_kavling: "A-1", kode_map: "M133.008..." }, ...]`.
   - Abaikan `<path>` tanpa `id` (itu dekorasi/teks), `<line>`, `<rect>`, dll.

2. Buat **CRUD API untuk Denah Kavling**:

   | Method | Endpoint                      | Keterangan                                                    |
   |--------|-------------------------------|---------------------------------------------------------------|
   | GET    | `/denah-kavling`              | List semua denah kavling (dengan relasi kavling-kavlingnya)    |
   | GET    | `/denah-kavling/:id`          | Detail satu denah beserta kavling-kavlingnya                   |
   | POST   | `/denah-kavling`              | Buat denah baru: terima `{ nama, svg_content }`               |
   |        |                               | → Parse SVG, simpan denah, otomatis buat kavling per path ber-id |
   | PUT    | `/denah-kavling/:id`          | Update nama denah                                              |
   | DELETE | `/denah-kavling/:id`          | Hapus denah beserta semua kavling di dalamnya                  |

3. **Logika `POST /denah-kavling`**:
   ```
   1. Terima { nama, svg_content }
   2. Panggil SVG Parser → validasi + extract paths
   3. Jika tidak ada path ber-id → return error "SVG tidak memiliki path kavling"
   4. Simpan ke tabel denah_kavling
   5. Untuk setiap path ber-id:
      - Buat record kavling dengan:
        - denah_kavling_id = denah.id
        - kode_kavling = path.id (misal "A-1")
        - kode_map = path.d
        - status = 0 (Kosong)
        - field dimensi & harga = 0 (default, diisi manual nanti)
   6. Return denah beserta kavling-kavlingnya
   ```

4. Update **API Kavling** yang sudah ada:
   - `PUT /kavling/:id` → update field dimensi, harga, status (kode_kavling & kode_map tidak boleh diubah).
   - `GET /kavling?denah_kavling_id=X` → filter kavling berdasarkan denah.
   - Hapus endpoint `POST /kavling/upload-svg` (digantikan `POST /denah-kavling`).
   - Hapus `POST /kavling` individual (kavling dibuat otomatis dari SVG).


## Contoh Alur User

```
1. User buka halaman /kavling
2. Terlihat list denah kavling (kosong jika baru pertama kali)
3. User klik "Tambah Denah Kavling"
4. Muncul modal dengan input nama dan textarea SVG
5. User isi nama: "Blok A"
6. User paste kode SVG ke textarea
7. User klik Submit
8. Backend parse SVG → ditemukan 11 path ber-id (A-1 sampai A-11)
9. Backend buat 1 denah + 11 kavling otomatis
10. Redirect ke /kavling/1 (detail denah)
11. Tampil peta SVG dengan 11 kavling berwarna hijau
12. User klik kavling A-1 di peta → popup muncul
13. User klik Edit → isi dimensi & harga
14. User klik Booking pada kavling A-1 → warna berubah kuning
```

---

## Catatan Penting

- **Hanya `<path>` dengan atribut `id`** yang dijadikan kavling. Path tanpa `id` (biasanya teks/label) dan elemen lain (`<line>`, `<rect>`) tetap di-render sebagai dekorasi statis.
- **SVG disimpan utuh** di database (`svg_content`). Frontend yang bertanggung jawab me-render dan mengganti warna fill sesuai status.
- **`kode_kavling` dan `kode_map` bersifat immutable** setelah dibuat. Jika user ingin mengubah layout SVG, hapus denah lalu buat ulang.
- Validasi SVG minimal: cek ada tag `<svg>`, cek ada minimal 1 `<path>` ber-id.
