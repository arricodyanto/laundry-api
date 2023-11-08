package utils

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Validate(err error, message string, c *gin.Context, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		log.Println(err, "Transaction has been rolled back.")
	} else {
		log.Println("Successfully " + message)
	}
}

func ErrorRecover(c *gin.Context) error {
	err := recover()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	return nil
}

func FormatStringToTime(dateString string) time.Time {
	format := "2006-01-02"
	result, err := time.Parse(format, dateString)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func FormatStringToInt(char string) int {
	result, _ := strconv.Atoi(char)
	return result
}

func FormatIntToString(number int) string {
	return strconv.Itoa(number)
}
