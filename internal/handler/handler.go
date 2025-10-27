package handler

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

	"REST-service-sub/internal/model"
	"REST-service-sub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	svc      service.SubscriptionServiceInterface
	validate *validator.Validate
}

func NewSubscriptionHandler(svc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		svc:      svc,
		validate: validator.New(),
	}
}

func (h *SubscriptionHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/subscriptions", h.Create)
	r.GET("/subscriptions/:id", h.Get)
	r.PUT("/subscriptions/:id", h.Update)
	r.DELETE("/subscriptions/:id", h.Delete)
	r.GET("/subscriptions", h.List)
	r.GET("/subscriptions/aggregate", h.Aggregate)
}

// Create Subscription godoc
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param payload body CreateSubscriptionDTO true "*a field end_date is optional*"
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var dto CreateSubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := h.validate.Struct(dto); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	startDate, err := ParseMonthYear(dto.StartDate)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid start date")
		return
	}

	var endDate *time.Time
	if dto.EndDate != nil {
		ed, err := ParseMonthYear(*dto.EndDate)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "invalid end_date")
			return
		}
		endDate = &ed
	}

	uid, err := uuid.Parse(dto.UserID)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	sub := &model.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      uid,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.svc.Create(sub); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, sub)
}

// Get Subscription godoc
// @Summary Get subscription by ID
// @Description Retrieve a single subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid id")
		return
	}
	sub, err := h.svc.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			respondWithError(c, http.StatusNotFound, "id not found")
			return
		}
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Update Subscription godoc
// @Summary Update subscription
// @Description Update existing subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param payload body CreateSubscriptionDTO true "Updated subscription payload"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var dto CreateSubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	// validate, parse dates, build model same as Create
	startDate, err := ParseMonthYear(dto.StartDate)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid start date")
		return
	}

	var endDate *time.Time
	if dto.EndDate != nil {
		ed, err := ParseMonthYear(*dto.EndDate)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "invalid end date")
			return
		}
		endDate = &ed
	}

	uid, err := uuid.Parse(dto.UserID)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	updated := &model.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      uid,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.svc.Update(id, updated); err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			respondWithError(c, http.StatusNotFound, "subscription not found")
			return
		}
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete Subscription godoc
// @Summary Delete subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			respondWithError(c, http.StatusInternalServerError, "subscription not found")
			return
		}
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// List Subscriptions godoc
// @Summary List subscriptions
// @Description Retrieve all subscriptions with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "Filter by user UUID"
// @Param service_name query string false "Filter by service name"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {array} SubscriptionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	filter := make(map[string]interface{})

	if userID := c.Query("user_id"); userID != "" {
		uid, err := uuid.Parse(userID)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "invalid user_id format")
			return
		}
		filter["user_id"] = uid

	}

	if serviceName := c.Query("service_name"); serviceName != "" {
		filter["service_name"] = serviceName
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	subs, err := h.svc.List(filter, limit, offset)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, subs)
}

// Aggregate Subscriptions godoc
// @Summary Aggregate subscription costs
// @Description Calculate total subscription cost for a given period
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param from query string true "Start of period (MM-YYYY or YYYY-MM)"
// @Param to query string true "End of period (MM-YYYY or YYYY-MM)"
// @Param user_id query string false "Filter by user UUID"
// @Param service_name query string false "Filter by service name"
// @Success 200 {object} models.AggregateResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/aggregate [get]
func (h *SubscriptionHandler) Aggregate(c *gin.Context) {
	from := c.Query("from") // expecting YYYY-MM or MM-YYYY
	to := c.Query("to")
	if from == "" || to == "" {
		respondWithError(c, http.StatusBadRequest, "invalid from or to")
		return
	}
	pFrom, err := ParseMonthYear(from)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid from date")
		return
	}
	pTo, err := ParseMonthYear(to)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid to date")
		return
	}
	var uid *uuid.UUID
	if userID := c.Query("user_id"); userID != "" {
		u, err := uuid.Parse(userID)
		if err == nil {
			uid = &u
		}
	}
	var svcName *string
	if sn := c.Query("service_name"); sn != "" {
		svcName = &sn
	}
	total, err := h.svc.AggregateTotalCost(pFrom, pTo, uid, svcName)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total_cost": total,
		"currency":   "RUB",
		"from":       from,
		"to":         to,
	})
}
