package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"subscription-service/internal/services"

	"github.com/google/uuid"
)

type SubsHandler struct {
	service ISubsServiceHandler
}

func RegisterSubsRoutes(mux *http.ServeMux, service ISubsServiceHandler) {
	h := &SubsHandler{service: service}

	mux.HandleFunc("/subs/create", h.Create())
	mux.HandleFunc("/subs/get", h.GetByID())
	mux.HandleFunc("/subs/list", h.List())
	mux.HandleFunc("/subs/update", h.Update())
	mux.HandleFunc("/subs/delete", h.Delete())
	mux.HandleFunc("/subs/total", h.CalculateTotal())
}

type CreateSubsReq struct {
	ServiceName string `json:"service_name"`
	Price       uint   `json:"price"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubsReq true "subscription data"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subs/create [post]
func (h *SubsHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var body CreateSubsReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Println(body)
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		userID, err := uuid.Parse(body.UserId)
		if err != nil {
			log.Printf("Create: ivalid user_id type, id=%s", userID)
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("01-2006-02", body.StartDate)
		if err != nil {
			log.Printf("Create: invalid data type, id=%s", startDate)
			http.Error(w, "use MM-YYYY", http.StatusBadRequest)
			return
		}

		var endDate *time.Time
		if body.EndDate != "" {
			t, err := time.Parse("01-2006-02", body.EndDate)
			if err != nil {
				log.Printf("Create: invalid data type, id=%s", endDate)
				http.Error(w, "use MM-YYYY", http.StatusBadRequest)
				return
			}
			endDate = &t
		}

		subs, err := h.service.Create(
			r.Context(),
			userID,
			body.ServiceName,
			body.Price,
			startDate,
			endDate,
		)
		if err != nil {
			log.Printf("Create: internal error, id=%s: %v", userID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(subs)
	}
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Get a subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id query string true "subscription ID (UUID)"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subs/get [get]
func (h *SubsHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userID, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("GetByID: subscription not found, id=%s", userID)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		subs, err := h.service.GetByID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				log.Printf("GetByID: subscription not found, id=%s", userID)
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			log.Printf("GetByID: internal error, id=%s: %v", userID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subs)
	}
}

type UpdateSubsReq struct {
	ServiceName string `json:"service_name"`
	Price       uint   `json:"price"`
	EndDate     string `json:"end_date"`
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update an exsisting subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query string true "subscription ID (UUID)"
// @Param subscription body CreateSubsReq true "subscription update data"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subs/update [put]
func (h *SubsHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userID, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("Update: subscription not found, id=%s", userID)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var req UpdateSubsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		var endDate *time.Time
		if req.EndDate != "" {
			t, err := time.Parse("2006-01-02", req.EndDate)
			if err != nil {
				log.Printf("Update: Incalid date type, date=%s", endDate)
				http.Error(w, "invalid date", http.StatusBadRequest)
				return
			}
			endDate = &t
		}

		subs, err := h.service.Update(
			r.Context(),
			userID,
			req.ServiceName,
			req.Price,
			endDate,
		)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				http.Error(w, "invalid id", http.StatusNotFound)
				return
			}
			log.Printf("Update: internal error not found, id=%s, %v", userID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subs)
	}
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Delete a subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id query string true "subscription ID (UUID)"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subs/delete [delete]
func (h *SubsHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userID, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("Delete: subscription not found, id=%s", userID)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := h.service.Delete(r.Context(), userID); err != nil {
			if errors.Is(err, services.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			log.Printf("Delete: internal error not found, id=%s: %v", userID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ListSubscription godoc
// @Summary List subscription by ID
// @Description List a subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID (UUID)"
// @Param service_name query string false "Filter by service name"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subs/list [get]
func (h *SubsHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userIDParam := r.URL.Query().Get("user_id")
		userID, err := uuid.Parse(userIDParam)
		if err != nil {
			log.Printf("List: subscription not found, id=%s", userID)
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		serviceName := r.URL.Query().Get("service_name")

		subs, err := h.service.List(
			r.Context(),
			userID,
			serviceName,
		)

		if err != nil {
			log.Printf("List: internal error, id=%s: %v", userID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subs)
	}
}

// CalculateCost godoc
// @Summary Calculate subscription cost
// @Description Calculate total cost of subscriptions for a period with optional filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID (UUID)"
// @Param start query string true "Start date (YYYY-MM-DD)"
// @Param end query string true "End date (YYYY-MM-DD)"
// @Param service_name query string false "Filter by service name"
// @Success 200 {object} map[string]uint
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subs/total [get]
func (h *SubsHandler) CalculateTotal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
		if err != nil {
			log.Printf("CalculateTotal: subscription not found, id=%s", userID)
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("2006-01-02", r.URL.Query().Get("start"))
		if err != nil {
			log.Printf("CalculateTotal: invalid start date, id=%v", startDate)
			http.Error(w, "use YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", r.URL.Query().Get("end"))
		if err != nil {
			log.Printf("CalculateTotal: invalid end date, id=%v", endDate)
			http.Error(w, "use YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		total, err := h.service.CalculateTotal(
			r.Context(),
			userID,
			r.URL.Query().Get("service_name"),
			startDate,
			endDate,
		)
		if err != nil {
			log.Printf("CalculateTotal: internal error, id=%s: %v", userID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]uint{"total": total})
	}
}
