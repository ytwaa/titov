package models

import (
	"fmt"
	"time"
)

type Expense struct {
	ID          int
	Sum         float64
	Description string
	Date        time.Time
}

func (e Expense) SumFormatted() string {
	return fmt.Sprintf("%.2f", e.Sum)
}
