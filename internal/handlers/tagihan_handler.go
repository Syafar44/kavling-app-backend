package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
)

type TagihanHandler struct {
	tagihanRepo  *repositories.TagihanRepository
	pembayaranV2 *repositories.PembayaranV2Repository
	kavlingRepo  *repositories.KavlingV2Repository
	keuanganRepo *repositories.KeuanganRepository
	trxKavRepo   *repositories.TransaksiKavlingRepository
}

func NewTagihanHandler(
	tr *repositories.TagihanRepository,
	pv2 *repositories.PembayaranV2Repository,
	kr *repositories.KavlingV2Repository,
	keu *repositories.KeuanganRepository,
	trx *repositories.TransaksiKavlingRepository,
) *TagihanHandler {
	return &TagihanHandler{tr, pv2, kr, keu, trx}
}

// GET /api/pembayaran — list semua transaksi kavling + total tagihan & terbayar
func (h *TagihanHandler) List(c *gin.Context) {
	// Tidak ada backfill — semua data dibuat saat customer disimpan di customer handler.
	list, err := h.pembayaranV2.List(c.Query("q"), c.Query("progres"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: list})
}

// GET /api/pembayaran/debug — counts of related tables (troubleshooting)
func (h *TagihanHandler) Debug(c *gin.Context) {
	// Read-only — hanya tampilkan counts tanpa modifikasi data.
	data := h.trxKavRepo.DebugCounts()
	c.JSON(http.StatusOK, models.Response{Success: true, Data: data})
}

// GET /api/pembayaran/:id_transaksi — detail (tagihan + pemasukan)
func (h *TagihanHandler) DetailPembayaran(c *gin.Context) {
	idTrx, _ := strconv.Atoi(c.Param("id_transaksi"))
	data, err := h.pembayaranV2.Detail(idTrx)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "transaksi tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: data})
}

// POST /api/pembayaran/:id_transaksi/tagihan
func (h *TagihanHandler) TambahTagihan(c *gin.Context) {
	idTrx, _ := strconv.Atoi(c.Param("id_transaksi"))
	var req struct {
		IDKategori int     `json:"id_kategori" binding:"required"`
		Deskripsi  string  `json:"deskripsi"`
		Nominal    float64 `json:"nominal" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}
	m := &models.Tagihan{
		IDTransaksi: idTrx,
		IDKategori:  req.IDKategori,
		Deskripsi:   req.Deskripsi,
		Nominal:     req.Nominal,
	}
	if err := h.tagihanRepo.Create(m); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "tagihan ditambahkan", Data: m})
}

// PUT /api/pembayaran/:id_transaksi/tagihan/:id
func (h *TagihanHandler) UpdateTagihan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	idTrx, _ := strconv.Atoi(c.Param("id_transaksi"))
	existing, err := h.tagihanRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tagihan tidak ditemukan"})
		return
	}
	var req struct {
		IDKategori int     `json:"id_kategori"`
		Deskripsi  string  `json:"deskripsi"`
		Nominal    float64 `json:"nominal"`
		SyncHarga  bool    `json:"sync_harga"` // jika true, ambil harga terbaru dari kavling
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	if req.SyncHarga {
		// Ambil harga kavling terbaru (sumber kebenaran)
		trxRec, trxErr := h.trxKavRepo.FindByID(idTrx)
		if trxErr == nil && trxRec != nil {
			kav, kavErr := h.kavlingRepo.FindByID(trxRec.IDKavling)
			if kavErr == nil && kav != nil && kav.HargaJualCash > 0 {
				existing.Nominal = kav.HargaJualCash
				existing.Deskripsi = "Harga Unit Kavling"
			}
		}
	} else {
		// Update manual
		if req.IDKategori > 0 {
			existing.IDKategori = req.IDKategori
		}
		if req.Deskripsi != "" {
			existing.Deskripsi = req.Deskripsi
		}
		if req.Nominal > 0 {
			existing.Nominal = req.Nominal
		}
	}

	if err := h.tagihanRepo.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "berhasil diupdate", Data: existing})
}

// DELETE /api/pembayaran/:id_transaksi/tagihan/:id
func (h *TagihanHandler) HapusTagihan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.tagihanRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "tagihan dihapus"})
}

// POST /api/pembayaran/:id_transaksi/pemasukan
func (h *TagihanHandler) TambahPemasukan(c *gin.Context) {
	idTrx, _ := strconv.Atoi(c.Param("id_transaksi"))
	var req struct {
		Tanggal    string  `json:"tanggal" binding:"required"`
		IDKategori int     `json:"id_kategori" binding:"required"`
		CaraBayar  string  `json:"cara_bayar" binding:"required"`
		IDBank     *int    `json:"id_bank"`
		Nominal    float64 `json:"nominal" binding:"required,gt=0"`
		Keterangan string  `json:"keterangan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	tgl, _ := time.Parse("2006-01-02", req.Tanggal)

	// Validate tidak melebihi tagihan
	totalTagihan, _ := h.tagihanRepo.SumByTransaksi(idTrx)
	totalBayar := h.pembayaranV2.SumPembayaran(idTrx)
	if totalBayar+req.Nominal > totalTagihan && totalTagihan > 0 {
		c.JSON(http.StatusUnprocessableEntity, models.Response{
			Message: "nominal melebihi sisa tagihan",
		})
		return
	}

	// Ambil id_customer & id_kavling dari transaksi_kavling (kolom NOT NULL di pembayaran)
	trxRec, _ := h.trxKavRepo.FindByID(idTrx)

	noPembayaran := "PAY-" + time.Now().Format("20060102") + "-" + strconv.Itoa(idTrx) + "-" + strconv.Itoa(h.pembayaranV2.NextPembayaranKe(idTrx))
	pembayaran := &models.Pembayaran{
		NoPembayaran: noPembayaran,
		IDTransaksi:  idTrx,
		IDKategori:   &req.IDKategori,
		CaraBayar:    req.CaraBayar,
		IDBank:       req.IDBank,
		Tanggal:      tgl,
		PembayaranKe: h.pembayaranV2.NextPembayaranKe(idTrx),
		JumlahBayar:  req.Nominal,
		Keterangan:   req.Keterangan,
	}
	if trxRec != nil {
		pembayaran.IDCustomer = trxRec.IDCustomer
		pembayaran.IDKavling = trxRec.IDKavling
	}
	if err := h.pembayaranV2.CreatePembayaran(pembayaran); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	// Catat ke arus kas
	noTrx := h.keuanganRepo.GenerateNoTransaksi("IN")
	trx := &models.Transaksi{
		NoTransaksi:   noTrx,
		Jenis:         "Pemasukan",
		Kategori:      "Pembayaran Kavling",
		Nominal:       req.Nominal,
		IDBank:        req.IDBank,
		Tanggal:       tgl,
		ReferensiTipe: "pembayaran",
		ReferensiID:   &pembayaran.ID,
		Keterangan:    req.Keterangan,
	}
	h.keuanganRepo.CreateTransaksi(trx)

	// Update saldo bank
	if req.IDBank != nil {
		h.keuanganRepo.AddSaldo(*req.IDBank, req.Nominal)
	}

	// TASK 4: Auto-update status kavling ke LUNAS (5) jika sudah lunas
	newTotalBayar := totalBayar + req.Nominal
	if totalTagihan > 0 && newTotalBayar >= totalTagihan && trxRec != nil {
		_ = h.kavlingRepo.UpdateStatus(trxRec.IDKavling, 5, nil) // 5 = LUNAS
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "pemasukan dicatat", Data: pembayaran})
}

// DELETE /api/pembayaran/:id_transaksi/pemasukan/:id
func (h *TagihanHandler) HapusPemasukan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	p, err := h.pembayaranV2.FindPembayaranByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Message: "tidak ditemukan"})
		return
	}
	// Kurangi saldo bank
	if p.IDBank != nil {
		h.keuanganRepo.SubtractSaldo(*p.IDBank, p.JumlahBayar)
	}
	if err := h.pembayaranV2.DeletePembayaran(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "pemasukan dihapus"})
}
