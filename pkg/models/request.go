package models

import (
	"errors"
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
