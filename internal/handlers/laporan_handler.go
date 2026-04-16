package handlers

import (
	"fmt"
	"net/http"
	"time"

	"backend-kavling/internal/helpers"
	"backend-kavling/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type LaporanHandler struct {
	arusKasRepo *repositories.ArusKasRepository
}

func NewLaporanHandler(arusKasRepo *repositories.ArusKasRepository) *LaporanHandler {
	return &LaporanHandler{arusKasRepo: arusKasRepo}
}

// Umum godoc
//
//	@Summary		Laporan arus kas umum
//	@Description	Mengambil laporan arus kas dengan filter rentang tanggal beserta total pemasukan, pengeluaran, dan saldo
//	@Tags			Laporan
//	@Produce		json
//	@Security		BearerAuth
//	@Param			dari	query		string	false	"Dari tanggal (YYYY-MM-DD)"
//	@Param			sampai	query		string	false	"Sampai tanggal (YYYY-MM-DD)"
//	@Success		200		{object} object
//	@Router			/laporan/umum [get]
func (h *LaporanHandler) Umum(c *gin.Context) {
	dari := c.Query("dari")
	sampai := c.Query("sampai")

	list, err := h.arusKasRepo.GetLaporanUmum(dari, sampai)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil laporan")
		return
	}

	pemasukan, pengeluaran := h.arusKasRepo.SumByJenis(dari, sampai)

	helpers.OK(c, "OK", gin.H{
		"transaksi":   list,
		"pemasukan":   pemasukan,
		"pengeluaran": pengeluaran,
		"saldo":       pemasukan - pengeluaran,
		"dari":        dari,
		"sampai":      sampai,
	})
}

// PerCustomer godoc
//
//	@Summary		Laporan per customer
//	@Description	Mengambil riwayat pembayaran/transaksi untuk customer tertentu dalam rentang tanggal
//	@Tags			Laporan
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id_customer	query		int		true	"Customer ID"
//	@Param			dari		query		string	false	"Dari tanggal (YYYY-MM-DD)"
//	@Param			sampai		query		string	false	"Sampai tanggal (YYYY-MM-DD)"
//	@Success		200			{object} object
//	@Failure		400			{object} object
//	@Router			/laporan/per-customer [get]
func (h *LaporanHandler) PerCustomer(c *gin.Context) {
	idCustomer := 0
	fmt.Sscanf(c.Query("id_customer"), "%d", &idCustomer)
	if idCustomer == 0 {
		helpers.BadRequest(c, "id_customer wajib diisi", nil)
		return
	}

	dari := c.Query("dari")
	sampai := c.Query("sampai")

	list, err := h.arusKasRepo.FindByCustomer(idCustomer, dari, sampai)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil laporan per customer")
		return
	}

	helpers.OK(c, "OK", gin.H{
		"transaksi":   list,
		"id_customer": idCustomer,
		"dari":        dari,
		"sampai":      sampai,
	})
}

// PerKategori godoc
//
//	@Summary		Laporan per kategori
//	@Description	Mengambil transaksi arus kas berdasarkan kategori tertentu dalam rentang tanggal
//	@Tags			Laporan
//	@Produce		json
//	@Security		BearerAuth
//	@Param			kategori	query		string	true	"Nama kategori (contoh: Bayar Angsuran)"
//	@Param			dari		query		string	false	"Dari tanggal (YYYY-MM-DD)"
//	@Param			sampai		query		string	false	"Sampai tanggal (YYYY-MM-DD)"
//	@Success		200			{object} object
//	@Failure		400			{object} object
//	@Router			/laporan/per-kategori [get]
func (h *LaporanHandler) PerKategori(c *gin.Context) {
	kategori := c.Query("kategori")
	if kategori == "" {
		helpers.BadRequest(c, "kategori wajib diisi", nil)
		return
	}

	dari := c.Query("dari")
	sampai := c.Query("sampai")

	list, err := h.arusKasRepo.FindByKategori(kategori, dari, sampai)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil laporan per kategori")
		return
	}

	helpers.OK(c, "OK", gin.H{
		"transaksi": list,
		"kategori":  kategori,
		"dari":      dari,
		"sampai":    sampai,
	})
}

// ExportExcel godoc
//
//	@Summary		Export laporan arus kas ke Excel
//	@Description	Mengunduh laporan arus kas dalam format .xlsx dengan filter rentang tanggal
//	@Tags			Laporan
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Security		BearerAuth
//	@Param			dari	query		string	false	"Dari tanggal (YYYY-MM-DD)"
//	@Param			sampai	query		string	false	"Sampai tanggal (YYYY-MM-DD)"
//	@Success		200		{file}		binary
//	@Router			/laporan/export-excel [get]
func (h *LaporanHandler) ExportExcel(c *gin.Context) {
	dari := c.Query("dari")
	sampai := c.Query("sampai")

	list, err := h.arusKasRepo.GetLaporanUmum(dari, sampai)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data untuk export")
		return
	}

	f := excelize.NewFile()
	sheet := "Laporan Arus Kas"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"No", "No. Transaksi", "Tanggal", "Jenis", "Kategori", "Keterangan", "Nominal", "Bank"}
	for i, h := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s1", col), h)
	}

	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
	})
	f.SetCellStyle(sheet, "A1", "H1", style)

	for i, t := range list {
		row := i + 2
		bankNama := ""
		if t.Bank != nil {
			bankNama = t.Bank.NamaBank
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t.NoTransaksi)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), t.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.Jenis)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), t.Kategori)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), t.Keterangan)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), t.Nominal)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), bankNama)
	}

	filename := fmt.Sprintf("laporan_arus_kas_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Cache-Control", "no-cache")

	if err := f.Write(c.Writer); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}

// ExportExcelPembayaran godoc
//
//	@Summary		Export laporan pembayaran cicilan ke Excel
//	@Description	Mengunduh laporan riwayat pembayaran cicilan dalam format .xlsx
//	@Tags			Laporan
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Security		BearerAuth
//	@Param			dari	query		string	false	"Dari tanggal (YYYY-MM-DD)"
//	@Param			sampai	query		string	false	"Sampai tanggal (YYYY-MM-DD)"
//	@Success		200		{file}		binary
//	@Router			/laporan/export-excel-pembayaran [get]
func (h *LaporanHandler) ExportExcelPembayaran(c *gin.Context) {
	dari := c.Query("dari")
	sampai := c.Query("sampai")

	list, err := h.arusKasRepo.FindPembayaranForExport(dari, sampai)
	if err != nil {
		helpers.InternalError(c, "Gagal mengambil data pembayaran untuk export")
		return
	}

	f := excelize.NewFile()
	sheet := "Laporan Pembayaran"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"No", "No. Pembayaran", "Tanggal", "Kavling", "Customer", "Cicilan Ke", "Jumlah Bayar", "Bank"}
	for i, h := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s1", col), h)
	}

	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
	})
	f.SetCellStyle(sheet, "A1", "H1", style)

	for i, p := range list {
		row := i + 2
		kavlingKode, customerNama, bankNama := "", "", ""
		if p.Kavling != nil {
			kavlingKode = p.Kavling.KodeKavling
		}
		if p.Customer != nil {
			customerNama = p.Customer.Nama
		}
		if p.Bank != nil {
			bankNama = p.Bank.NamaBank
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.NoPembayaran)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), kavlingKode)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), customerNama)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.PembayaranKe)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), p.JumlahBayar)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), bankNama)
	}

	filename := fmt.Sprintf("laporan_pembayaran_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	f.Write(c.Writer)
}
