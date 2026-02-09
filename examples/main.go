package main

import (
	"fmt"
	"log"
	"time"

	ninepay "github.com/9pay-labs/9pay-sdk-golang"
	"github.com/9pay-labs/9pay-sdk-golang/pkg/models"
)

func main() {
	// 1. Khởi tạo
	// const MerchantKey =  os.Getenv("NINEPAY_MERCHANT_KEY") // "pxxxw"
	// const SecretKey = os.Getenv("NINEPAY_SECRET_KEY") //"narlsaxxxxxxxxxxxxAtvJgAKSiQOg"
	// const CheckSum = os.Getenv("NINEPAY_CHECKSUM_KEY") //"s6KiGBywxxxxxxxxxxxxxxsx4QHM2YWzLC"
	// const Endpoint = os.Getenv("NINEPAY_ENDPOINT") // https://xxxxx.9pay.mobi

	const MerchantKey = "pxxxw"
	const SecretKey = "narlsaxxxxxxxxxxxxAtvJgAKSiQOg"
	const CheckSum = "s6KiGBywxxxxxxxxxxxxxxsx4QHM2YWzLC"
	const Endpoint = "https://xxxxx.9pay.mobi"

	client := ninepay.New(MerchantKey, SecretKey, CheckSum, Endpoint)

	// 2. Tạo Link Thanh Toán
	req := models.New[models.BuildUrlRequest]()
	req.InvoiceNo = "ORDER_001"
	req.Amount = 50000
	req.ReturnUrl = "https://myshop.com/result"
	req.Description = "Mô tả thông tin đơn hàng"
	req.Time = time.Now().Unix()
	req.MerchantKey = MerchantKey

	url, err := client.BuildPaymentURL(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("👉 Payment URL:", url)

	// // 3. Kiểm tra trạng thái
	fmt.Println("\nChecking Status for ORDER_001...")
	transaction, err := client.Inquire("ORDER_001")
	if err != nil {
		log.Printf("❌ Check failed: %v", err)
	} else {
		fmt.Printf("✅ Status: %d \n PaymentNo: %d", *transaction.Status, *transaction.PaymentNo)
	}

	// // 4. Test Verify Webhook (Giả lập)
	fmt.Println("\nTesting Verify Webhook...")

	const data = "eyJhbW91bnQiOxxxxxm9yIjpudWxsfQ"
	const checksum = "8FD0C7C97ACE326xxxxxxxxxB798F818B8FCB049B6A"
	IsValid := client.VerifyChecksum(data, checksum)
	log.Printf("❌ Check sum: %v", IsValid)
}
