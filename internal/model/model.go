package model

import "time"

type Subscription struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"id"`
	ServiceName string    `gorm:"not null" json:"service_name"`
	Price       int       `gorm:"not null" json:"price"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	EndDate     time.Time `gorm:"not null" json:"end_date"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}
