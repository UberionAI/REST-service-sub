package handler

import (
	"fmt"
	"time"
)

type CreateSubscriptionDTO struct {
	ServiceName string  `json:"serviceName" validate:"required"`
	Price       int     `json:"price" validate:"required,min=0"`
	UserID      string  `json:"userId" validate:"required,uuid4"`
	StartDate   string  `json:"startDate" validate:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

func ParseMonthYear(s string) (time.Time, error) {
	var t time.Time
	var err error
	layout := []string{"01-2006", "2006-01"}
	for _, l := range layout {
		t, err = time.Parse(l, s)
		if err == nil {
			return time.Date(t.Year(), t.Month(), t.Day(), 1, 0, 0, 0, time.UTC), nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}
