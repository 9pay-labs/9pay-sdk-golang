package models

import (
	"encoding/json"
	"fmt"
)

type APIResponse struct {
	ErrorCode string          `json:"error_code"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"data"`
}
type StringOrNumber string

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
type PayerAuthResponse struct {
	RequestID     string `json:"request_id"`
	OrderCode     int64  `json:"order_code"`
	ThreeDSURL    string `json:"3ds_url"`
	Status        int    `json:"status"`
	ErrorCode     string `json:"error_code"`
	FailureReason string `json:"failure_reason"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	SubmittedAt   string `json:"submitted_at"` // ISO8601 UTC+0
}

// -------- Response --------

type AuthorizeResponse struct {
	RequestID     string `json:"request_id"`
	OrderCode     int64  `json:"order_code"`
	AuthCode      string `json:"auth_code"`
	Status        int    `json:"status"`
	ErrorCode     string `json:"error_code"`
	FailureReason string `json:"failure_reason"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	SubmittedAt   string `json:"submitted_at"`
}

type CaptureResponse struct {
	RequestID   string `json:"request_id"`
	OrderCode   int64  `json:"order_code"`
	Status      int    `json:"status"`
	ErrorCode   string `json:"error_code"`
	FailureCode string `json:"failure_code"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	SubmittedAt string `json:"submitted_at"`
}

type ReverseAuthResponse struct {
	RequestID   string `json:"request_id"`
	OrderCode   int64  `json:"order_code"`
	Status      int    `json:"status"`
	ErrorCode   string `json:"error_code"`
	FailureCode string `json:"failure_code"`
	SubmittedAt string `json:"submitted_at"`
}

type RefundCreateResponse struct {
	Status        int            `json:"status"`
	ErrorCode     StringOrNumber `json:"error_code"`
	FailureReason string         `json:"failure_reason,omitempty"`
	Amount        float64        `json:"amount"`
	OrderCode     int64          `json:"order_code"`
	RefundNo      int64          `json:"refund_no,omitempty"`
	RequestID     string         `json:"request_id"`
}

func (s *StringOrNumber) UnmarshalJSON(b []byte) error {
	// string
	if b[0] == '"' {
		var str string
		if err := json.Unmarshal(b, &str); err != nil {
			return err
		}
		*s = StringOrNumber(str)
		return nil
	}

	// number
	var num interface{}
	if err := json.Unmarshal(b, &num); err != nil {
		return err
	}

	*s = StringOrNumber(fmt.Sprintf("%v", num))
	return nil
}
