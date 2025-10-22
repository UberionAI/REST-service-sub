package service

import (
	"REST-service-sub/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "postgres://postgres:5qL2eTbSj1@127.0.0.1:5432/subscription_db?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Exec("DELETE FROM subscriptions").Error; err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCRUD(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSubscriptionService(db)

	sub := &model.Subscription{
		ID:          uuid.New(),
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
	}
	err := svc.Create(sub)
	assert.NoError(t, err)

	got, err := svc.GetByID(sub.ID)
	assert.NoError(t, err)
	assert.Equal(t, sub.ServiceName, got.ServiceName)

	sub.Price = 500
	err = svc.Update(sub.ID, sub)
	assert.NoError(t, err)

	got, _ = svc.GetByID(sub.ID)
	assert.Equal(t, 500, got.Price)

	list, err := svc.List(nil, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	err = svc.Delete(sub.ID)
	assert.NoError(t, err)

	_, err = svc.GetByID(sub.ID)
	assert.Error(t, err)
}

func TestAggregateTotalCost(t *testing.T) {
	db := setupTestDB(t)
	svc := NewSubscriptionService(db)

	userID := uuid.New()
	db.Create(&model.Subscription{
		ID:          uuid.New(),
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      userID,
		StartDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	})

	from := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)

	total, err := svc.AggregateTotalCost(from, to, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(400), total)
}
