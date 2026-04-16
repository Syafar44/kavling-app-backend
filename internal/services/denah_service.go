package services

import (
	"errors"

	"backend-kavling/internal/models"
	"backend-kavling/internal/repositories"
)

type DenahService struct {
	denahRepo   *repositories.DenahRepository
	kavlingRepo *repositories.KavlingRepository
}

func NewDenahService(denahRepo *repositories.DenahRepository, kavlingRepo *repositories.KavlingRepository) *DenahService {
	return &DenahService{denahRepo: denahRepo, kavlingRepo: kavlingRepo}
}

// CreateDenah parses the SVG string, saves the denah, and batch-inserts all kavlings.
func (s *DenahService) CreateDenah(nama, svgContent string) (*models.DenahKavling, error) {
	viewbox, paths, err := ParseSVGString(svgContent)
	if err != nil {
		return nil, err
	}

	denah := &models.DenahKavling{
		Nama:       nama,
		SvgContent: svgContent,
		Viewbox:    viewbox,
	}

	if err := s.denahRepo.Create(denah); err != nil {
		return nil, errors.New("gagal menyimpan denah kavling: " + err.Error())
	}

	kavlings := make([]models.KavlingPeta, 0, len(paths))
	for _, p := range paths {
		denahID := denah.ID
		kavlings = append(kavlings, models.KavlingPeta{
			DenahKavlingID: &denahID,
			KodeKavling:    p.KodeKavling,
			KodeMap:        p.KodeMap,
			Status:         0,
		})
	}

	if err := s.kavlingRepo.CreateBatch(kavlings); err != nil {
		// Rollback denah on failure
		_ = s.denahRepo.Delete(denah.ID)
		return nil, errors.New("gagal menyimpan kavling: " + err.Error())
	}

	// Return denah with kavlings populated
	result, err := s.denahRepo.FindByID(denah.ID)
	if err != nil {
		return denah, nil
	}
	return result, nil
}

// UpdateSVG replaces the SVG content of an existing denah (without changing kavlings).
func (s *DenahService) UpdateSVG(id int, nama, svgContent string) (*models.DenahKavling, error) {
	denah, err := s.denahRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("denah tidak ditemukan")
	}

	denah.Nama = nama
	if svgContent != "" {
		viewbox, _, err := ParseSVGString(svgContent)
		if err != nil {
			return nil, err
		}
		denah.SvgContent = svgContent
		denah.Viewbox = viewbox
	}

	if err := s.denahRepo.Update(denah); err != nil {
		return nil, errors.New("gagal update denah: " + err.Error())
	}

	return denah, nil
}

// DeleteDenah removes a denah only if it has no transactions.
func (s *DenahService) DeleteDenah(id int) error {
	if s.denahRepo.HasTransaksi(id) {
		return errors.New("denah memiliki kavling dengan transaksi aktif, tidak bisa dihapus")
	}
	return s.denahRepo.Delete(id)
}
