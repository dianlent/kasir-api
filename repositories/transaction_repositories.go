package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetTodayReport() (*models.TodayReport, error) {
	rows, err := repo.db.Query("SELECT id, total_amount, created_at FROM transactions WHERE created_at::date = CURRENT_DATE ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	totalAmount := 0
	totalTransactions := 0

	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.TotalAmount, &tx.CreatedAt); err != nil {
			return nil, err
		}

		totalAmount += tx.TotalAmount
		totalTransactions++
		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.TodayReport{
		Date:              time.Now().Format("2006-01-02"),
		TotalAmount:       totalAmount,
		TotalTransactions: totalTransactions,
		Transactions:      transactions,
	}, nil
}

func (repo *TransactionRepository) GetReport(startDate time.Time, endDate time.Time) (*models.DateRangeReport, error) {
	endExclusive := endDate.AddDate(0, 0, 1)
	rows, err := repo.db.Query(
		"SELECT id, total_amount, created_at FROM transactions WHERE created_at >= $1 AND created_at < $2 ORDER BY created_at DESC",
		startDate, endExclusive,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	totalAmount := 0
	totalTransactions := 0

	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.TotalAmount, &tx.CreatedAt); err != nil {
			return nil, err
		}

		totalAmount += tx.TotalAmount
		totalTransactions++
		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.DateRangeReport{
		StartDate:         startDate.Format("2006-01-02"),
		EndDate:           endDate.Format("2006-01-02"),
		TotalAmount:       totalAmount,
		TotalTransactions: totalTransactions,
		Transactions:      transactions,
	}, nil
}
