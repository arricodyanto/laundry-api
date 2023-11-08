package utils

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Validate(err error, message string, c *gin.Context, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}
}

func FormatStringToTime(dateString string) time.Time {
	format := "2006-01-02"
	result, _ := time.Parse(format, dateString)
	return result
}

func FormatStringToInt(char string) int {
	result, _ := strconv.Atoi(char)
	return result
}

func FormatIntToString(number int) string {
	return strconv.Itoa(number)
}
