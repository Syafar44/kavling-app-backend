package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"backend-kavling/internal/config"
	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
)

type CustomerV2Handler struct {
	repo         *repositories.CustomerV2Repository
	kavlingRepo  *repositories.KavlingV2Repository
	lokasiRepo   *repositories.LokasiKavlingRepository
	tagihanRepo  *repositories.TagihanRepository
	trxKavRepo   *repositories.TransaksiKavlingRepository
	keuanganRepo *repositories.KeuanganRepository
}

func NewCustomerV2Handler(
	r   *repositories.CustomerV2Repository,
	k   *repositories.KavlingV2Repository,
	l   *repositories.LokasiKavlingRepository,
	t   *repositories.TagihanRepository,
	trx *repositories.TransaksiKavlingRepository,
	keu *repositories.KeuanganRepository,
) *CustomerV2Handler {
	return &CustomerV2Handler{r, k, l, t, trx, keu}
}

// GET /api/customer
func (h *CustomerV2Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	list, total, err := h.repo.List(
		c.Query("q"), c.Query("id_lokasi"),
		c.Query("status"), c.Query("jenis_pembelian"),
		page, perPage,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.PaginatedResponse{Success: true, Data: list, Total: total, Page: page, PerPage: perPage})
}

// GET /api/customer/:id
func (h *CustomerV2Handler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "customer tidak ditemukan"})
		return
	}
	kavlings, _ := h.repo.GetKavlings(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: gin.H{
		"customer": m,
		"kavlings": kavlings,
	}})
}

// POST /api/customer
func (h *CustomerV2Handler) Create(c *gin.Context) {
	var body struct {
		models.Customer
		IDKavling *int `json:"id_kavling"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	m := body.Customer

	// Auto-gen kode kontrak
	namaLokasi := ""
	if m.IDLokasi != nil {
		if lok, err := h.lokasiRepo.FindByID(*m.IDLokasi); err == nil {
			namaLokasi = lok.Nama
		}
	}
	kode, _ := h.repo.GenerateKodeKontrak(namaLokasi, m.JenisPembelian)
	m.KodeKontrak = kode

	if err := h.repo.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	log.Printf("[CustomerCreate] customer=%d id_kavling=%v jenis=%q", m.ID, body.IDKavling, m.JenisPembelian)
	if body.IDKavling != nil && *body.IDKavling > 0 {
		h.assignKavling(m.ID, *body.IDKavling, m.JenisPembelian)
		h.generatePembayaran(&m, *body.IDKavling)
	} else {
		log.Printf("[CustomerCreate] SKIP generatePembayaran: id_kavling nil/0")
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "customer berhasil dibuat", Data: m})
}

// PUT /api/customer/:id
func (h *CustomerV2Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	var body struct {
		models.Customer
		IDKavling *int `json:"id_kavling"`
	}
	body.Customer = *existing
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	body.Customer.ID = id
	if err := h.repo.Update(&body.Customer); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	if body.IDKavling != nil {
		h.reassignKavling(id, *body.IDKavling, body.Customer.JenisPembelian)
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Data: body.Customer})
}

// assignKavling links a kavling to the customer: creates the join row and sets
// kavling.id_customer + flips status from Ready (0) to BF (2).
// Status yang diset tergantung jenis pembelian:
// - BOOKING FEE → BF (2)
// - CASH KERAS  → BF (2), akan otomatis jadi LUNAS (5) saat lunas
// - KREDIT      → BF (2), akan jadi AKAD (3) saat proses
func (h *CustomerV2Handler) assignKavling(idCustomer, idKavling int, jenisPembelian string) {
	_ = h.repo.AddKavling(idCustomer, idKavling)
	k, err := h.kavlingRepo.FindByID(idKavling)
	if err != nil {
		return
	}
	newStatus := k.Status
	if newStatus == 0 {
		switch strings.ToUpper(strings.TrimSpace(jenisPembelian)) {
		case "BOOKING FEE", "CASH KERAS", "KREDIT":
			newStatus = 2 // BF — status awal semua jenis pembelian
		default:
			newStatus = 2
		}
	}
	_ = h.kavlingRepo.UpdateStatus(idKavling, newStatus, &idCustomer)
}

// reassignKavling replaces the customer's current kavling links with the one
// provided. If idKavling == 0, only clears existing links.
func (h *CustomerV2Handler) reassignKavling(idCustomer, idKavling int, jenisPembelian string) {
	existingLinks, _ := h.repo.GetKavlings(idCustomer)
	for _, ck := range existingLinks {
		if ck.IDKavling == idKavling {
			continue
		}
		_ = h.repo.RemoveKavling(idCustomer, ck.IDKavling)
		k, err := h.kavlingRepo.FindByID(ck.IDKavling)
		if err != nil {
			continue
		}
		k.IDCustomer = nil
		if k.Status == 2 {
			k.Status = 0
		}
		_ = h.kavlingRepo.Update(k)
	}
	if idKavling > 0 {
		alreadyLinked := false
		for _, ck := range existingLinks {
			if ck.IDKavling == idKavling {
				alreadyLinked = true
				break
			}
		}
		if !alreadyLinked {
			h.assignKavling(idCustomer, idKavling, jenisPembelian)
		}
	}
}

// DELETE /api/customer/:id
func (h *CustomerV2Handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Reset status semua kavling milik customer ini sebelum dihapus.
	links, _ := h.repo.GetKavlings(id)
	for _, ck := range links {
		k, err := h.kavlingRepo.FindByID(ck.IDKavling)
		if err != nil {
			continue
		}
		k.IDCustomer = nil
		k.Status = 0
		_ = h.kavlingRepo.Update(k)
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "berhasil dihapus"})
}

// ─── Customer File (Upload) ───────────────────────────────────────────────────

// GET /api/customer/:id/file
func (h *CustomerV2Handler) ListFile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := h.repo.ListFiles(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

// POST /api/customer/:id/file
func (h *CustomerV2Handler) UploadFile(c *gin.Context) {
	idCustomer, _ := strconv.Atoi(c.Param("id"))
	namaFile := c.PostForm("nama_file")
	keterangan := c.PostForm("keterangan")
	tanggalStr := c.PostForm("tanggal")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: "file wajib diupload"})
		return
	}

	// Simpan file
	now := time.Now()
	dir := filepath.Join(config.AppConfig.UploadPath, "customer", now.Format("2006"), now.Format("01"))
	os.MkdirAll(dir, 0755)
	ext := filepath.Ext(file.Filename)
	filename := strconv.Itoa(idCustomer) + "_" + now.Format("20060102150405") + ext
	dstPath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: "gagal menyimpan file"})
		return
	}
	relativePath := "customer/" + now.Format("2006") + "/" + now.Format("01") + "/" + filename

	tgl := now
	if tanggalStr != "" {
		tgl, _ = time.Parse("2006-01-02", tanggalStr)
	}

	m := &models.CustomerFile{
		IDCustomer: idCustomer,
		Tanggal:    tgl,
		NamaFile:   namaFile,
		PathFile:   relativePath,
		Keterangan: keterangan,
	}
	if err := h.repo.CreateFile(m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

// DELETE /api/customer/:id/file/:id_file
func (h *CustomerV2Handler) DeleteFile(c *gin.Context) {
	idFile, _ := strconv.Atoi(c.Param("id_file"))
	if err := h.repo.DeleteFile(idFile); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "file dihapus"})
}

// ─── Arsip Customer ───────────────────────────────────────────────────────────

// GET /api/customer/arsip
func (h *CustomerV2Handler) ListArsip(c *gin.Context) {
	list, err := h.repo.ListArsip(c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

// DELETE /api/customer/arsip/:id
func (h *CustomerV2Handler) DeleteArsip(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.repo.DeleteArsip(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "arsip dihapus"})
}

// ─── Prospek ──────────────────────────────────────────────────────────────────

// GET /api/prospek
func (h *CustomerV2Handler) ListProspek(c *gin.Context) {
	list, err := h.repo.ListProspek(c.Query("q"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

// POST /api/prospek
func (h *CustomerV2Handler) CreateProspek(c *gin.Context) {
	var m models.Prospek
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	if m.Tanggal.IsZero() {
		m.Tanggal = time.Now()
	}
	if err := h.repo.CreateProspek(&m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.Response{Success: true, Data: m})
}

// PUT /api/prospek/:id
func (h *CustomerV2Handler) UpdateProspek(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	existing, err := h.repo.FindProspekByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	c.ShouldBindJSON(existing)
	existing.ID = id
	h.repo.UpdateProspek(existing)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: existing})
}

// DELETE /api/prospek/:id
func (h *CustomerV2Handler) DeleteProspek(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.repo.DeleteProspek(id)
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "dihapus"})
}

// POST /api/prospek/:id/convert — prospek → customer
func (h *CustomerV2Handler) ConvertProspek(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	p, err := h.repo.FindProspekByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "prospek tidak ditemukan"})
		return
	}
	// Return prospek data yang bisa di-copy ke form customer
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "data prospek siap dikonversi",
		Data: gin.H{
			"nama":         p.Nama,
			"no_telp":      p.NoTelp,
			"pekerjaan":    p.Pekerjaan,
			"id_marketing": p.IDMarketing,
			"id_prospek":   p.ID,
		},
	})
}

// generatePembayaran creates a TransaksiKavling + Tagihan records based on jenis_pembelian.
func (h *CustomerV2Handler) generatePembayaran(customer *models.Customer, idKavling int) {
	log.Printf("[generatePembayaran] START customer=%d kavling=%d jenis=%q harga=%v jumlah_pembayaran=%v",
		customer.ID, idKavling, customer.JenisPembelian, customer.HargaJual, customer.JumlahPembayaran)

	now := time.Now()

	jenisPembelianInt := 1
	switch customer.JenisPembelian {
	case "CASH KERAS":
		jenisPembelianInt = 2
	case "KREDIT":
		jenisPembelianInt = 3
	}

	noTransaksi := fmt.Sprintf("TRX-%s-%d", now.Format("20060102"), customer.ID)

	trx := &models.TransaksiKavling{
		NoTransaksi:    noTransaksi,
		IDKavling:      idKavling,
		IDCustomer:     customer.ID,
		IDMarketing:    customer.IDMarketing,
		JenisPembelian: jenisPembelianInt,
		HargaJual:      customer.HargaJual,
		TglTransaksi:   now,
	}
	if customer.CicilanPerBulan != nil {
		trx.CicilanPerBulan = *customer.CicilanPerBulan
	}
	if customer.Tenor != nil {
		trx.LamaCicilan = *customer.Tenor
	}

	// TASK 6: Set tgl_mulai_cicilan otomatis untuk KREDIT
	if strings.ToUpper(strings.TrimSpace(customer.JenisPembelian)) == "KREDIT" &&
		customer.JatuhTempo != nil && *customer.JatuhTempo != "" {
		jtStr := strings.TrimSpace(*customer.JatuhTempo)

		if len(jtStr) >= 10 {
			if parsedDate, err := time.Parse("2006-01-02", jtStr[:10]); err == nil {
				trx.TglMulaiCicilan = &parsedDate
				log.Printf("[generatePembayaran] tgl_mulai_cicilan=%s (dari full date)", parsedDate.Format("2006-01-02"))
			}
		} else {
			// Fallback: jika isinya hanya angka 1-31 (format lama)
			tglInt, err := strconv.Atoi(jtStr)
			if err == nil && tglInt >= 1 && tglInt <= 31 {
				nextMonth := now.AddDate(0, 1, 0)
				mulaiCicilan := time.Date(nextMonth.Year(), nextMonth.Month(), tglInt, 0, 0, 0, 0, now.Location())
				trx.TglMulaiCicilan = &mulaiCicilan
				log.Printf("[generatePembayaran] tgl_mulai_cicilan=%s (dari angka hari)", mulaiCicilan.Format("2006-01-02"))
			}
		}
	}

	if err := h.trxKavRepo.Create(trx); err != nil {
		log.Printf("[generatePembayaran] ERROR create TransaksiKavling: %v", err)
		return
	}
	log.Printf("[generatePembayaran] TransaksiKavling %s dibuat (id=%d)", trx.NoTransaksi, trx.ID)

	repositories.CreateTagihanForCustomer(trx.ID, customer, h.tagihanRepo)
	log.Printf("[generatePembayaran] DONE tagihan created for trx=%d", trx.ID)

	// Catat pendapatan ke arus kas berdasarkan jenis pembelian
	h.catatPendapatan(customer, trx.ID)
}

// catatPendapatan mencatat transaksi pemasukan ke tabel transaksi (arus kas)
// berdasarkan jenis pembelian customer yang baru disimpan.
func (h *CustomerV2Handler) catatPendapatan(customer *models.Customer, idTrx int) {
	if h.keuanganRepo == nil {
		log.Printf("[catatPendapatan] keuanganRepo nil, skip")
		return
	}

	jenis := strings.ToUpper(strings.TrimSpace(customer.JenisPembelian))

	var nominal float64
	var kategori string
	var keterangan string

	switch jenis {
	case "BOOKING FEE":
		nominal = customer.JumlahPembayaran // pembayaran BF yang sudah dibayar
		kategori = "Booking Fee"
		keterangan = fmt.Sprintf("Booking Fee - %s", customer.Nama)
	case "CASH KERAS":
		nominal = customer.HargaJual // harga jual kavling penuh
		kategori = "Pembayaran Kavling"
		keterangan = fmt.Sprintf("Penjualan Cash - %s", customer.Nama)
	case "KREDIT":
		// DP awal: ambil dari field uang_muka jika ada, fallback ke JumlahPembayaran
		nominal = customer.JumlahPembayaran
		if nominal <= 0 && customer.CicilanPerBulan != nil {
			nominal = 0 // belum ada DP yang masuk
		}
		kategori = "Pembayaran Kavling"
		keterangan = fmt.Sprintf("DP Kredit - %s", customer.Nama)
	default:
		log.Printf("[catatPendapatan] jenis pembelian %q tidak dikenal, skip", jenis)
		return
	}

	if nominal <= 0 {
		log.Printf("[catatPendapatan] nominal 0 untuk %s (%s), skip pencatatan arus kas", customer.Nama, jenis)
		return
	}

	noTrx := h.keuanganRepo.GenerateNoTransaksi("IN")
	trx := &models.Transaksi{
		NoTransaksi:   noTrx,
		Jenis:         "Pemasukan",
		Kategori:      kategori,
		Nominal:       nominal,
		Tanggal:       time.Now(),
		ReferensiTipe: "pembayaran",
		ReferensiID:   &idTrx,
		Keterangan:    keterangan,
	}
	if err := h.keuanganRepo.CreateTransaksi(trx); err != nil {
		log.Printf("[catatPendapatan] ERROR CreateTransaksi: %v", err)
		return
	}
	log.Printf("[catatPendapatan] %s dicatat: %s %.0f (trx=%d)", jenis, noTrx, nominal, idTrx)
}
