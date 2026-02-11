package utils

import "testing"

func TestBuildCanonicalParams_Basic(t *testing.T) {
	data := map[string]interface{}{
		"amount":     10000,
		"currency":   "VND",
		"request_id": "REQ001",
	}

	got := BuildCanonicalParams(data)
	expected := "amount=10000&currency=VND&request_id=REQ001"

	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestBuildCanonicalParams_SortedKeys(t *testing.T) {
	data := map[string]interface{}{
		"b": 2,
		"a": 1,
		"c": 3,
	}

	got := BuildCanonicalParams(data)
	expected := "a=1&b=2&c=3"

	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestBuildCanonicalParams_URLEncode(t *testing.T) {
	data := map[string]interface{}{
		"return_url": "https://merchant.example.com/return",
	}

	got := BuildCanonicalParams(data)
	expected := "return_url=https%3A%2F%2Fmerchant.example.com%2Freturn"

	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestBuildCanonicalParams_SpecialChars(t *testing.T) {
	data := map[string]interface{}{
		"description": "Thanh toán đơn hàng #123",
	}

	got := BuildCanonicalParams(data)
	expected := "description=Thanh+to%C3%A1n+%C4%91%C6%A1n+h%C3%A0ng+%23123"

	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestBuildCanonicalParams_Empty(t *testing.T) {
	data := map[string]interface{}{}

	got := BuildCanonicalParams(data)
	if got != "" {
		t.Fatalf("expected empty string, got %s", got)
	}
}
