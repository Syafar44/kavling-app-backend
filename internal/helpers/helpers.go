package helpers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// Response helpers
func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(201, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func BadRequest(c *gin.Context, message string, errors interface{}) {
	c.JSON(400, gin.H{
		"success": false,
		"message": message,
		"errors":  errors,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(401, gin.H{
		"success": false,
		"message": message,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(403, gin.H{
		"success": false,
		"message": message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(404, gin.H{
		"success": false,
		"message": message,
	})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(500, gin.H{
		"success": false,
		"message": message,
	})
}

func Paginated(c *gin.Context, message string, data interface{}, total int64, page, perPage int) {
	c.JSON(200, gin.H{
		"success":  true,
		"message":  message,
		"data":     data,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// GenerateNoTransaksi generates a transaction number: XXXXX/TRX-KAV/YYYY
func GenerateNoTransaksi(prefix string) string {
	now := time.Now()
	rnd := rand.Intn(90000) + 10000
	return fmt.Sprintf("%d/%s/%d", rnd, prefix, now.Year())
}

// GenerateNoBooking generates booking number: XXXXX/BKG-KAV/YYYY
func GenerateNoBooking() string {
	return GenerateNoTransaksi("BKG-KAV")
}

// GenerateNoPembayaran generates payment number: XXXXX/BYR-KAV/YYYY
func GenerateNoPembayaran() string {
	return GenerateNoTransaksi("BYR-KAV")
}

// GenerateNoArusKas generates cash flow number: XXXXX/ARS-KAV/YYYY
func GenerateNoArusKas() string {
	return GenerateNoTransaksi("ARS-KAV")
}

// GetClientIP returns IP address from context
func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}

// FormatRupiah formats float as Indonesian Rupiah
func FormatRupiah(amount float64) string {
	return fmt.Sprintf("Rp %.0f", amount)
}

// MonthDiff returns the number of months between two dates
func MonthDiff(from, to time.Time) int {
	years := to.Year() - from.Year()
	months := int(to.Month()) - int(from.Month())
	return years*12 + months
}

// Ptr returns pointer to int
func IntPtr(i int) *int { return &i }
