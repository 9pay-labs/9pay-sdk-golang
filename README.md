# 9Pay Go SDK

Official Go SDK for integrating with **9Pay Payment Gateway**.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Features

- **Payment URL Generation** - Create secure payment links for customers
- **Transaction Inquiry** - Check payment status by invoice number
- **Webhook Verification** - Verify incoming webhook data using checksums
- **Zero Dependencies** - Pure Go implementation using only standard library
- **Secure by Design** - HMAC-SHA256 signing and constant-time comparison

## Installation

```bash
go get github.com/9pay-labs/9pay-sdk-golang
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"

    ninepay "github.com/9pay-labs/9pay-sdk-golang"
    "github.com/9pay-labs/9pay-sdk-golang/pkg/models"
)

func main() {
    // Initialize client
    client := ninepay.New(
        "your-merchant-key",
        "your-secret-key",
        "your-checksum-key",
        "https://sandbox.9pay.mobi",
    )

    // Create payment URL
    req := models.New[models.BuildUrlRequest]()
    req.InvoiceNo = "ORDER_001"
    req.Amount = 50000
    req.ReturnUrl = "https://myshop.com/result"
    req.Description = "Order description"
    req.Time = time.Now().Unix()
    req.MerchantKey = "your-merchant-key"

    url, err := client.BuildPaymentURL(req)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Payment URL:", url)
}
```

## Configuration

Initialize the client with your 9Pay credentials:

```go
client := ninepay.New(key, secret, checksum, endpoint)
```

| Parameter | Description |
|-----------|-------------|
| `key` | Merchant key for authentication |
| `secret` | Secret key for HMAC signing |
| `checksum` | Checksum key for webhook verification |
| `endpoint` | 9Pay API endpoint URL |

**Recommended:** Load credentials from environment variables:

```go
client := ninepay.New(
    os.Getenv("NINEPAY_MERCHANT_KEY"),
    os.Getenv("NINEPAY_SECRET_KEY"),
    os.Getenv("NINEPAY_CHECKSUM_KEY"),
    os.Getenv("NINEPAY_ENDPOINT"),
)
```

## API Reference

### BuildPaymentURL

Generate a secure payment URL for customers to complete transactions.

```go
func (c *Client) BuildPaymentURL(req *models.BuildUrlRequest) (string, error)
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `InvoiceNo` | string | Yes | Unique order/invoice identifier |
| `Amount` | int | Yes | Payment amount (smallest currency unit) |
| `Description` | string | Yes | Order description |
| `ReturnUrl` | string | Yes | URL to redirect after payment |
| `Time` | int64 | No | Unix timestamp |
| `MerchantKey` | string | No | Merchant identifier |

**Example:**

```go
req := models.New[models.BuildUrlRequest]()
req.InvoiceNo = "ORDER_12345"
req.Amount = 100000  // 100,000 VND
req.ReturnUrl = "https://myshop.com/payment/result"
req.Description = "Payment for Order #12345"
req.Time = time.Now().Unix()
req.MerchantKey = "your-merchant-key"

url, err := client.BuildPaymentURL(req)
if err != nil {
    log.Fatal(err)
}
// Redirect customer to this URL
```

### Inquire

Query transaction status by invoice number.

```go
func (c *Client) Inquire(invoiceNo string) (*models.TransactionInquire, error)
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `PaymentNo` | *int | Payment reference number |
| `InvoiceNo` | *string | Original invoice number |
| `Amount` | *int | Paid amount |
| `Status` | *int | Payment status code |
| `Currency` | *string | Currency code |
| `Method` | *string | Payment method used |
| `CardBrand` | *string | Card brand (if applicable) |
| `Description` | *string | Transaction description |
| `CreatedAt` | *int64 | Creation timestamp |

**Example:**

```go
transaction, err := client.Inquire("ORDER_12345")
if err != nil {
    log.Printf("Inquiry failed: %v", err)
    return
}

fmt.Printf("Payment No: %d\n", *transaction.PaymentNo)
fmt.Printf("Status: %d\n", *transaction.Status)
fmt.Printf("Amount: %d\n", *transaction.Amount)
```

### VerifyChecksum

Verify webhook data integrity using checksum.

```go
func (c *Client) VerifyChecksum(data, checksum string) bool
```

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `data` | string | Base64-encoded webhook payload |
| `checksum` | string | Checksum from webhook header |

**Example:**

```go
// In your webhook handler
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    data := r.FormValue("data")
    checksum := r.FormValue("checksum")

    if !client.VerifyChecksum(data, checksum) {
        http.Error(w, "Invalid checksum", http.StatusBadRequest)
        return
    }

    // Process the webhook data
    decoded, _ := base64.StdEncoding.DecodeString(data)
    // Parse and handle the payment notification
}
```

## Security

This SDK implements several security measures:

- **HMAC-SHA256 Signing**: All API requests are signed with your secret key
- **Base64 Encoding**: Payment data is encoded before transmission
- **Constant-Time Comparison**: Webhook verification uses constant-time comparison to prevent timing attacks
- **Request Validation**: All required fields are validated before API calls

## Examples

See the [examples](./examples) directory for complete usage examples:

```bash
cd examples
go run main.go
```

## Requirements

- Go 1.22 or higher

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

For support and questions, please contact 9Pay support or open an issue on GitHub.
