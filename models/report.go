package models

type TodayReport struct {
	Date              string        `json:"date"`
	TotalAmount       int           `json:"total_amount"`
	TotalTransactions int           `json:"total_transactions"`
	Transactions      []Transaction `json:"transactions"`
}

type DateRangeReport struct {
	StartDate         string        `json:"start_date"`
	EndDate           string        `json:"end_date"`
	TotalAmount       int           `json:"total_amount"`
	TotalTransactions int           `json:"total_transactions"`
	Transactions      []Transaction `json:"transactions"`
}
