package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ServiceName string     `gorm:"type:text;not null" json:"service_name" validate:"required"`
	Price       int        `gorm:"not null" json:"price" validate:"required,min=0"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id" validate:"required"`
	StartDate   time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate     *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
