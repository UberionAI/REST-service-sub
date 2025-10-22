package handler

import (
	"REST-service-sub/internal/model"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockService struct {
	CreatedSub *model.Subscription
}

func (m *mockService) Create(sub *model.Subscription) error {
	sub.ID = uuid.New()
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()
	m.CreatedSub = sub
	return nil
}

func (m *mockService) GetByID(id uuid.UUID) (*model.Subscription, error) {
	return &model.Subscription{
		ID:          id,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}, nil
}

func (m *mockService) Update(id uuid.UUID, sub *model.Subscription) error {
	return nil
}

func (m *mockService) Delete(id uuid.UUID) error {
	return nil
}

func (m *mockService) List(filter map[string]interface{}, limit, offset int) ([]model.Subscription, error) {
	return []model.Subscription{
		{
			ID:          uuid.New(),
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      uuid.New(),
			StartDate:   time.Now(),
		},
	}, nil
}

func (m *mockService) AggregateTotalCost(start, end time.Time, userID *uuid.UUID, serviceName *string) (int64, error) {
	return 800, nil
}

func newTestHandler() *SubscriptionHandler {
	mockSvc := &mockService{}
	return &SubscriptionHandler{
		svc:      mockSvc,
		validate: validator.New(),
	}
}

//
// 3️⃣ Тесты
//

func TestCreateSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()
	router := gin.New()
	h.RegisterRoutes(router)

	dto := map[string]interface{}{
		"service_name": "Yandex Plus",
		"price":        400,
		"user_id":      uuid.New().String(),
		"start_date":   "07-2025",
	}
	body, _ := json.Marshal(dto)

	req, _ := http.NewRequest("POST", "/subscriptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp model.Subscription
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Yandex Plus", resp.ServiceName)
	assert.NotEmpty(t, resp.ID)
}

func TestGetSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()
	router := gin.New()
	h.RegisterRoutes(router)

	id := uuid.New()
	req, _ := http.NewRequest("GET", "/subscriptions/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var sub model.Subscription
	err := json.Unmarshal(w.Body.Bytes(), &sub)
	assert.NoError(t, err)
	assert.Equal(t, "Yandex Plus", sub.ServiceName)
}

func TestListSubscriptions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()
	router := gin.New()
	h.RegisterRoutes(router)

	req, _ := http.NewRequest("GET", "/subscriptions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var subs []model.Subscription
	err := json.Unmarshal(w.Body.Bytes(), &subs)
	assert.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, "Yandex Plus", subs[0].ServiceName)
}

func TestUpdateSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()
	router := gin.New()
	h.RegisterRoutes(router)

	id := uuid.New()
	dto := map[string]interface{}{
		"service_name": "Updated Name",
		"price":        999,
	}
	body, _ := json.Marshal(dto)

	req, _ := http.NewRequest("PUT", "/subscriptions/"+id.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()
	router := gin.New()
	h.RegisterRoutes(router)

	id := uuid.New()
	req, _ := http.NewRequest("DELETE", "/subscriptions/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestAggregateTotalCost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newTestHandler()
	router := gin.New()
	h.RegisterRoutes(router)

	req, _ := http.NewRequest("GET", "/subscriptions/aggregate?from=07-2025&to=08-2025", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(800), resp["total_cost"])
}
