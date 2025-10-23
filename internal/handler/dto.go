package handler

import (
	"fmt"
	"time"
)

// CreateSubscriptionDTO represents the data transfer object for creating or updating a subscription.
//
//swagger:model CreateSubscriptionDTO
type CreateSubscriptionDTO struct {
	//Name of the service (for example, "Spotify Premium")
	//required: true
	ServiceName string `json:"service_name" validate:"required"`
	//Price monthly in RUB
	//required: true
	Price int `json:"price" validate:"required,min=0"`
	//User UUID
	//required: true
	UserID string `json:"user_id" validate:"required,uuid4"`
	//Start date (formated as MM-YYYY)
	//required: true
	StartDate string `json:"start_date" validate:"required"`
	//End date (formated as MM-YYYY) *Optional*
	EndDate *string `json:"end_date,omitempty"`
}

// SubscriptionResponse represents the response structure for a subscription.
//
//swagger:model SubscriptionResponse
type SubscriptionResponse struct {
	ID          string     `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AggregatedResponse represents the response structure for aggregated data.
// swagger:model AggregatedResponse
type AggregatedResponse struct {
	//Total cost of the subscriptions with or without filters (service name or user id)
	TotalCost int `json:"total_cost" example:"10000"`
	//Currency of the total cost (RUB for example)
	Currency string `json:"currency" example:"RUB"`
	//Start point of the aggregated data
	FromDate string `json:"from_date" example:"01-2023"`
	//End point of the aggregated data (inclusive)
	ToDate string `json:"to_date" example:"02-2023"`
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
