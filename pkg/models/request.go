package models

import (
	"errors"
	"fmt"
	"strings"
)

type BaseRequest struct {
	Extra map[string]interface{}
}

func (b *BaseRequest) Init() {
	b.Extra = make(map[string]interface{})
}

func (b *BaseRequest) Set(key string, val interface{}) {
	if b.Extra == nil {
		b.Extra = make(map[string]interface{})
	}
	b.Extra[key] = val
}

type IBaseRequest interface {
	Init()
	ToMap() map[string]interface{}
	Set(key string, val interface{})
}
type BuildUrlRequest struct {
	BaseRequest
	InvoiceNo   string `json:"invoice_no,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
	ReturnUrl   string `json:"return_url,omitempty"`
	Time        int64  `json:"time,omitempty"`
	MerchantKey string `json:"merchantKey,omitempty"`

	extra map[string]interface{}
}

type InquireRequest struct {
	BaseRequest
}

type PayerAuthRequest struct {
	BaseRequest
	RequestID string `json:"request_id"`

	Currency string `json:"currency"`

	Amount float64 `json:"amount"`

	Card *CardInfo `json:"card"`

	ReturnURL string `json:"return_url"`

	Off3DS int8 `json:"off_3ds,omitempty"`
}

type AuthorizeRequest struct {
	BaseRequest
	RequestID string `json:"request_id"`

	Currency string `json:"currency"`

	Amount float64 `json:"amount"`

	Card *CardInfo `json:"card"`

	OrderCode int64 `json:"order_code"`
}

type CaptureRequest struct {
	BaseRequest
	RequestID string `json:"request_id"`
	Currency  string `json:"currency"`
	Amount    int64  `json:"amount"`
	OrderCode int64  `json:"order_code"`
}

type ReverseAuthRequest struct {
	BaseRequest
	RequestID string `json:"request_id"`
	OrderCode int64  `json:"order_code"`
}
type RefundCreateRequest struct {
	BaseRequest
	RequestID     string  `json:"request_id"`
	PaymentNo     int64   `json:"payment_no"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	Bank          *string `json:"bank,omitempty"`
	AccountNumber *string `json:"account_number,omitempty"`
	CustomerName  *string `json:"customer_name,omitempty"`
}

type CardInfo struct {
	CardNumber string `json:"card_number"`

	CardName string `json:"hold_name"`

	CardMonth string `json:"exp_month"`

	CardYear string `json:"exp_year"`

	CVV string `json:"cvv"`
}

func New[T any, P interface {
	*T
	IBaseRequest
}]() P {
	ptr := P(new(T))
	ptr.Init()

	return ptr
}

func (r *BuildUrlRequest) Validate() error {
	if strings.TrimSpace(r.InvoiceNo) == "" {
		return errors.New("missing required field: invoice_no")
	}
	if r.Amount <= 0 {
		return errors.New("invalid amount: must be greater than 0")
	}
	if strings.TrimSpace(r.Description) == "" {
		return errors.New("missing required field: description")
	}
	if strings.TrimSpace(r.ReturnUrl) == "" {
		return errors.New("missing required field: return_url")
	}
	if r.Time <= 0 {
		return errors.New("missing required field: time")
	}
	if strings.TrimSpace(r.MerchantKey) == "" {
		return errors.New("missing required field: merchantKey")
	}
	return nil
}

func (r *InquireRequest) Validate() error {
	return nil
}

func (r *PayerAuthRequest) Validate() error {
	if r == nil {
		return errors.New("request is nil")
	}

	if r.RequestID == "" {
		return errors.New("request_id is required")
	}
	if len(r.RequestID) > 30 {
		return errors.New("request_id max length is 30")
	}

	if r.Currency != "VND" {
		return errors.New("currency must be VND")
	}

	if r.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if r.Card == nil {
		return errors.New("card is required")
	}
	if err := r.Card.Validate(); err != nil {
		return errors.New("card: " + err.Error())
	}

	if r.ReturnURL == "" {
		return errors.New("return_url is required")
	}
	if r.Off3DS < 0 {
		if r.Off3DS != 0 && r.Off3DS != 1 {
			return fmt.Errorf("off_3ds must be 0 or 1")
		}
	}

	return nil
}

func (r *AuthorizeRequest) Validate() error {
	if r == nil {
		return errors.New("request is nil")
	}

	if r.RequestID == "" {
		return errors.New("request_id is required")
	}
	if len(r.RequestID) > 30 {
		return errors.New("request_id max length is 30")
	}

	if r.Currency != "VND" {
		return errors.New("currency must be VND")
	}

	if r.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if r.Card == nil {
		return errors.New("card is required")
	}
	if err := r.Card.Validate(); err != nil {
		return errors.New("card: " + err.Error())
	}

	if r.OrderCode <= 0 {
		return errors.New("order_code is required")
	}

	return nil
}

func (r *CaptureRequest) Validate() error {
	if r == nil {
		return errors.New("request is nil")
	}

	if r.RequestID == "" {
		return errors.New("request_id is required")
	}
	if len(r.RequestID) > 30 {
		return errors.New("request_id max length is 30")
	}

	if r.Currency != "VND" {
		return errors.New("currency must be VND")
	}

	if r.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if r.OrderCode <= 0 {
		return errors.New("order_code is required")
	}

	return nil
}

func (c *CardInfo) Validate() error {
	if c.CardNumber == "" {
		return fmt.Errorf("CardNumber is required")
	}
	if len(c.CardNumber) < 13 || len(c.CardNumber) > 19 {
		return fmt.Errorf("CardNumber length invalid")
	}

	if c.CardName == "" {
		return errors.New("CardName is required")
	}

	if len(c.CardMonth) != 2 {
		return errors.New("CardMonth must be 2 digits")
	}

	if len(c.CardYear) != 2 {
		return errors.New("CardYear must be 2 digits")
	}

	if len(c.CVV) < 3 || len(c.CVV) > 4 {
		return errors.New("cvv invalid")
	}

	return nil
}

func (r *ReverseAuthRequest) Validate() error {
	if r.RequestID == "" {
		return errors.New("request_id is required")
	}
	if len(r.RequestID) > 30 {
		return errors.New("request_id max length is 30")
	}
	if r.OrderCode <= 0 {
		return errors.New("order_code must be greater than 0")
	}
	return nil
}

func (r *RefundCreateRequest) Validate() error {
	if strings.TrimSpace(r.RequestID) == "" {
		return errors.New("request_id is required")
	}
	if len(r.RequestID) > 30 {
		return errors.New("request_id max length is 30")
	}
	if r.PaymentNo <= 0 {
		return errors.New("payment_no must be greater than 0")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if strings.TrimSpace(r.Description) == "" {
		return errors.New("description is required")
	}
	if len(r.Description) > 255 {
		return errors.New("description max length is 255")
	}

	if r.Bank != nil {
		if r.AccountNumber == nil || strings.TrimSpace(*r.AccountNumber) == "" {
			return errors.New("account_number is required when bank is provided")
		}
		if r.CustomerName == nil || strings.TrimSpace(*r.CustomerName) == "" {
			return errors.New("customer_name is required when bank is provided")
		}
	}

	return nil
}

func (r *ReverseAuthRequest) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"request_id": r.RequestID,
		"order_code": r.OrderCode,
	}
}

func (r *RefundCreateRequest) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"request_id":  r.RequestID,
		"payment_no":  r.PaymentNo,
		"amount":      r.Amount,
		"description": r.Description,
	}

	if r.Bank != nil {
		m["bank"] = *r.Bank
	}
	if r.AccountNumber != nil {
		m["account_number"] = *r.AccountNumber
	}
	if r.CustomerName != nil {
		m["customer_name"] = *r.CustomerName
	}

	return m
}
func (r *BuildUrlRequest) ToMap() map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range r.extra {
		merged[k] = v
	}
	if r.InvoiceNo != "" {
		merged["invoice_no"] = r.InvoiceNo
	}
	if r.Amount != 0 {
		merged["amount"] = r.Amount
	}
	if r.Description != "" {
		merged["description"] = r.Description
	}
	if r.ReturnUrl != "" {
		merged["return_url"] = r.ReturnUrl
	}
	if r.Time != 0 {
		merged["time"] = r.Time
	}

	if r.MerchantKey != "" {
		merged["merchantKey"] = r.MerchantKey
	}

	return merged
}

func (r *InquireRequest) ToMap() map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range r.Extra {
		merged[k] = v
	}
	return merged
}

func (r *PayerAuthRequest) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"request_id": r.RequestID,
		"currency":   r.Currency,
		"amount":     r.Amount,
		"return_url": r.ReturnURL,
	}

	if r.Off3DS > 0 {
		m["off_3ds"] = r.Off3DS
	}

	if r.Card != nil {
		m["card"] = r.Card.ToMap()
	}

	return m
}

func (r *AuthorizeRequest) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"request_id": r.RequestID,
		"currency":   r.Currency,
		"amount":     r.Amount,
		"order_code": r.OrderCode,
	}

	if r.Card != nil {
		m["card"] = r.Card.ToMap()
	}

	return m
}

func (r *CaptureRequest) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"request_id": r.RequestID,
		"currency":   r.Currency,
		"amount":     r.Amount,
		"order_code": r.OrderCode,
	}
	return m
}

func (c *CardInfo) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"card_number": c.CardNumber,
		"exp_month":   c.CardMonth,
		"cvv":         c.CVV,
		"hold_name":   c.CardName,
		"exp_year":    c.CardYear,
	}
}
