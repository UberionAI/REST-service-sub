package service

import (
	"REST-service-sub/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SubscriptionServiceInterface interface {
	Create(*model.Subscription) error
	GetByID(uuid.UUID) (*model.Subscription, error)
	Update(uuid.UUID, *model.Subscription) error
	Delete(uuid.UUID) error
	List(map[string]interface{}, int, int) ([]model.Subscription, error)
	AggregateTotalCost(time.Time, time.Time, *uuid.UUID, *string) (int64, error)
}

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
	sql := `
SELECT COALESCE(SUM(
    ((DATE_PART('year', LEAST(COALESCE(end_date, ?), ?)) * 12 + DATE_PART('month', LEAST(COALESCE(end_date, ?), ?)))
     -
    (DATE_PART('year', GREATEST(start_date, ?)) * 12 + DATE_PART('month', GREATEST(start_date, ?)))
     + 1) * price
),0)::bigint as total
FROM subscriptions
WHERE start_date <= ? AND (end_date IS NULL OR end_date >= ?)
`

	args := []interface{}{periodEnd, periodEnd, periodEnd, periodEnd, periodStart, periodStart, periodEnd, periodStart}

	// динамические фильтры
	if userID != nil {
		sql += " AND user_id = ?"
		args = append(args, *userID)
	}
	if serviceName != nil {
		sql += " AND service_name = ?"
		args = append(args, *serviceName)
	}

	var total int64
	if err := s.db.Raw(sql, args...).Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
