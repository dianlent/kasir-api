package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"time"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// multiple item apa aja, quantity nya
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) HandleTodayReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTodayReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	items, err := parseCheckoutItems(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) GetTodayReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetTodayReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *TransactionHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, time.Local)
	if err != nil {
		http.Error(w, "start_date must be in YYYY-MM-DD format", http.StatusBadRequest)
		return
	}

	endDate, err := time.ParseInLocation("2006-01-02", endDateStr, time.Local)
	if err != nil {
		http.Error(w, "end_date must be in YYYY-MM-DD format", http.StatusBadRequest)
		return
	}

	if endDate.Before(startDate) {
		http.Error(w, "end_date must be greater than or equal to start_date", http.StatusBadRequest)
		return
	}

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func parseCheckoutItems(r *http.Request) ([]models.CheckoutItem, error) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("failed to read request body")
	}

	if len(bytes.TrimSpace(rawBody)) == 0 {
		return nil, errors.New("request body is empty")
	}

	var req models.CheckoutRequest
	if err := json.Unmarshal(rawBody, &req); err == nil {
		return validateCheckoutItems(req.Items)
	}

	var items []models.CheckoutItem
	if err := json.Unmarshal(rawBody, &items); err == nil {
		return validateCheckoutItems(items)
	}

	return nil, errors.New("invalid request body")
}

func validateCheckoutItems(items []models.CheckoutItem) ([]models.CheckoutItem, error) {
	if len(items) == 0 {
		return nil, errors.New("items is required")
	}

	for i, item := range items {
		if item.ProductID <= 0 {
			return nil, fmt.Errorf("items[%d].product_id must be greater than 0", i)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("items[%d].quantity must be greater than 0", i)
		}
	}

	return items, nil
}
