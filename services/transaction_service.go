package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetTodayReport() (*models.TodayReport, error) {
	return s.repo.GetTodayReport()
}

func (s *TransactionService) GetReport(startDate time.Time, endDate time.Time) (*models.DateRangeReport, error) {
	return s.repo.GetReport(startDate, endDate)
}
