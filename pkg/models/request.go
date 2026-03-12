package models

import (
	"github.com/9pay-labs/9pay-sdk-golang/pkg/validator"
)

var engine = validator.New()

// --- INTERFACES & BASE ---

type IBaseRequest interface {
	Init()
	Validate() error
	ToMap() map[string]interface{}
	Set(key string, val interface{})
}

type BaseRequest struct {
	Extra map[string]interface{}
}

func (b *BaseRequest) Init() {
	if b.Extra == nil {
		b.Extra = make(map[string]interface{})
	}
}

func (b *BaseRequest) Set(key string, val interface{}) {
	b.Init()
	b.Extra[key] = val
}

func New[T any, P interface {
	*T
	IBaseRequest
}]() P {
	ptr := P(new(T))
	ptr.Init()
	return ptr
}

// --- MODELS ---

type CardInfo struct {
	CardNumber string `json:"card_number" validate:"required,numeric,min=13,max=19"`
	CardName   string `json:"hold_name" validate:"required"`
	CardMonth  string `json:"exp_month" validate:"required,len=2,numeric"`
	CardYear   string `json:"exp_year" validate:"required,len=2,numeric"`
	CVV        string `json:"cvv" validate:"required,min=3,max=4,numeric"`
}

type BuildUrlRequest struct {
	BaseRequest
	InvoiceNo   string `json:"invoice_no" validate:"required"`
	Amount      int64  `json:"amount" validate:"required,min=1"`
	Description string `json:"description" validate:"required"`
	ReturnUrl   string `json:"return_url" validate:"required,url"`
	Time        int64  `json:"time" validate:"required,min=1"`
	MerchantKey string `json:"merchantKey" validate:"required"`
	extra       map[string]interface{}
}

type PayerAuthRequest struct {
	BaseRequest
	RequestID string    `json:"request_id" validate:"required,max=30"`
	Currency  string    `json:"currency" validate:"required,is_vnd"`
	Amount    float64   `json:"amount" validate:"required,min=1"`
	Card      *CardInfo `json:"card" validate:"required"`
	ReturnURL string    `json:"return_url" validate:"required,url"`
	Off3DS    int8      `json:"off_3ds" validate:"min=0,max=1"`
}

type AuthorizeRequest struct {
	BaseRequest
	RequestID string    `json:"request_id" validate:"required,max=30"`
	Currency  string    `json:"currency" validate:"required,is_vnd"`
	Amount    float64   `json:"amount" validate:"required,min=1"`
	Card      *CardInfo `json:"card" validate:"required"`
	OrderCode int64     `json:"order_code" validate:"required,min=1"`
}

type CaptureRequest struct {
	BaseRequest
	RequestID string `json:"request_id" validate:"required,max=30"`
	Currency  string `json:"currency" validate:"required,is_vnd"`
	Amount    int64  `json:"amount" validate:"required,min=1"`
	OrderCode int64  `json:"order_code" validate:"required,min=1"`
}

type ReverseAuthRequest struct {
	BaseRequest
	RequestID string `json:"request_id" validate:"required,max=30"`
	OrderCode int64  `json:"order_code" validate:"required,min=1"`
}

type RefundCreateRequest struct {
	BaseRequest
	RequestID     string  `json:"request_id" validate:"required,max=30"`
	PaymentNo     int64   `json:"payment_no" validate:"required,min=1"`
	Amount        float64 `json:"amount" validate:"required,min=1"`
	Description   string  `json:"description" validate:"required,max=255"`
	Bank          *string `json:"bank,omitempty"`
	AccountNumber *string `json:"account_number,omitempty"`
	CustomerName  *string `json:"customer_name,omitempty"`
}
type InquireRequest struct {
	BaseRequest
}

func (r *BuildUrlRequest) Validate() error     { return engine.Validate(r) }
func (r *PayerAuthRequest) Validate() error    { return engine.Validate(r) }
func (r *AuthorizeRequest) Validate() error    { return engine.Validate(r) }
func (r *CaptureRequest) Validate() error      { return engine.Validate(r) }
func (r *ReverseAuthRequest) Validate() error  { return engine.Validate(r) }
func (r *RefundCreateRequest) Validate() error { return engine.Validate(r) }
func (r *InquireRequest) Validate() error      { return nil }
func (c *CardInfo) Validate() error            { return engine.Validate(c) }

func (r *BuildUrlRequest) ToMap() map[string]interface{}     { return validator.ToMap(r) }
func (r *PayerAuthRequest) ToMap() map[string]interface{}    { return validator.ToMap(r) }
func (r *AuthorizeRequest) ToMap() map[string]interface{}    { return validator.ToMap(r) }
func (r *CaptureRequest) ToMap() map[string]interface{}      { return validator.ToMap(r) }
func (r *ReverseAuthRequest) ToMap() map[string]interface{}  { return validator.ToMap(r) }
func (r *RefundCreateRequest) ToMap() map[string]interface{} { return validator.ToMap(r) }
func (r *InquireRequest) ToMap() map[string]interface{}      { return validator.ToMap(r) }
