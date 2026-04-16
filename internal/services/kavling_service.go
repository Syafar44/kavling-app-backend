package services

import (
	"encoding/xml"
	"errors"
	"strings"

	"backend-kavling/internal/repositories"
)

type KavlingService struct {
	repo *repositories.KavlingRepository
}

func NewKavlingService(repo *repositories.KavlingRepository) *KavlingService {
	return &KavlingService{repo: repo}
}

// ─── SVG Parsing ──────────────────────────────────────────────────────────────

type svgRoot struct {
	XMLName xml.Name   `xml:"svg"`
	Viewbox string     `xml:"viewBox,attr"`
	Groups  []svgGroup `xml:"g"`
	Paths   []svgPath  `xml:"path"`
}

type svgGroup struct {
	XMLName xml.Name   `xml:"g"`
	Groups  []svgGroup `xml:"g"`
	Paths   []svgPath  `xml:"path"`
}

type svgPath struct {
	ID string `xml:"id,attr"`
	D  string `xml:"d,attr"`
}

// ParsedPath represents a single kavling extracted from SVG
type ParsedPath struct {
	KodeKavling string
	KodeMap     string
}

// ParseSVGString validates and extracts <path id="..."> elements from an SVG string.
// Returns viewbox and list of parsed paths.
func ParseSVGString(svgContent string) (viewbox string, paths []ParsedPath, err error) {
	if !strings.Contains(svgContent, "<svg") {
		return "", nil, errors.New("konten bukan SVG yang valid: tidak ditemukan tag <svg>")
	}

	var root svgRoot
	if e := xml.Unmarshal([]byte(svgContent), &root); e != nil {
		return "", nil, errors.New("format SVG tidak valid: " + e.Error())
	}

	viewbox = root.Viewbox

	// Collect <path id="..."> recursively
	var collect func(groups []svgGroup, directPaths []svgPath)
	collect = func(groups []svgGroup, directPaths []svgPath) {
		for _, p := range directPaths {
			if strings.TrimSpace(p.ID) != "" {
				paths = append(paths, ParsedPath{
					KodeKavling: strings.TrimSpace(p.ID),
					KodeMap:     strings.TrimSpace(p.D),
				})
			}
		}
		for _, g := range groups {
			collect(g.Groups, g.Paths)
		}
	}
	collect(root.Groups, root.Paths)

	if len(paths) == 0 {
		return "", nil, errors.New("SVG tidak memiliki path kavling: tidak ditemukan <path> dengan atribut id")
	}

	return viewbox, paths, nil
}

// ValidateDelete checks if kavling can be deleted (only if no active transactions)
func (s *KavlingService) ValidateDelete(id int) error {
	k, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("kavling tidak ditemukan")
	}
	if k.Status != 0 {
		return errors.New("kavling sudah ada transaksi, tidak bisa dihapus")
	}
	return nil
}
