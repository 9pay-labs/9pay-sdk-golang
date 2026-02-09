package models

import "encoding/json"

type APIResponse struct {
	ErrorCode string          `json:"error_code"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"data"`
}

type ResponseInquire struct {
	// common
	ErrorCode     string `json:"error_code"`
	FailureReason string `json:"failure_reason"`
	Message       string `json:"message"`

	// transaction fields
	TransactionInquire
}

type TransactionInquire struct {
	PaymentNo   *int64  `json:"payment_no"`
	InvoiceNo   *string `json:"invoice_no"`
	Currency    *string `json:"currency"`
	Amount      *int64  `json:"amount"`
	Description *string `json:"description"`
	Method      *string `json:"method"`
	CardBrand   *string `json:"card_brand"`
	CreatedAt   *string `json:"created_at"`
	Status      *int    `json:"status"`
}
