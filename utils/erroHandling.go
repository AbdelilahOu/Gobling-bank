package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func ErrInvalidCurrency(accountID uuid.UUID, AccCurrency, currency string) error {
	return fmt.Errorf("account %d currency mismatch %s vs %s", accountID, AccCurrency, currency)
}
