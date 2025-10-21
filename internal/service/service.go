package service

import (
	"REST-service-sub/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SubscriptionService struct {
	db *gorm.DB
}

func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

func (s *SubscriptionService) Create(sub *model.Subscription) error {
	return s.db.Create(sub).Error
}

func (s *SubscriptionService) GetByID(id uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription
	if err := s.db.First(&sub, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *SubscriptionService) Update(id uuid.UUID, updated *model.Subscription) error {
	updated.ID = id
	return s.db.Model(&model.Subscription{}).Where("id = ?", id).Updates(updated).Error
}

func (s *SubscriptionService) Delete(id uuid.UUID) error {
	return s.db.Delete(&model.Subscription{}, "id = ?", id).Error
}

func (s *SubscriptionService) List(filter map[string]interface{}, limit, offset int) ([]model.Subscription, error) {
	var subs []model.Subscription
	tx := s.db.Model(&model.Subscription{})
	for k, v := range filter {
		tx = tx.Where(k+" = ?", v)
	}
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	if offset > 0 {
		tx = tx.Offset(offset)
	}
	if err := tx.Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *SubscriptionService) AggregateTotalCost(periodStart, periodEnd time.Time, userID *uuid.UUID, serviceName *string) (int64, error) {
	pFrom := periodStart.Format("2006-01-02")
	pTo := periodEnd.Format("2006-01-02")

	sql := `
SELECT COALESCE(SUM(
  ((DATE_PART('year', LEAST(COALESCE(end_date, :p_to), :p_to)) * 12 + DATE_PART('month', LEAST(COALESCE(end_date, :p_to), :p_to)))
   -
  (DATE_PART('year', GREATEST(start_date, :p_from)) * 12 + DATE_PART('month', GREATEST(start_date, :p_from)))
  + 1) * price
),0)::bigint as total
FROM subscriptions
WHERE start_date <= :p_to AND (end_date IS NULL OR end_date >= :p_from)
`
	// filters
	if userID != nil {
		sql += " AND user_id = :user_id"
	}
	if serviceName != nil {
		sql += " AND service_name = :service_name"
	}

	rows, err := s.db.Raw(sql, map[string]interface{}{
		"p_from":       pFrom,
		"p_to":         pTo,
		"user_id":      userID,
		"service_name": serviceName,
	}).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var total int64
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, nil
}
