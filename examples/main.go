package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	ninepay "github.com/9pay-labs/9pay-sdk-golang"
	"github.com/9pay-labs/9pay-sdk-golang/pkg/models"
)

func main() {
	// 1. Khởi tạo
	// const MerchantKey =  os.Getenv("NINEPAY_MERCHANT_KEY") // "pxxxw"
	// const SecretKey = os.Getenv("NINEPAY_SECRET_KEY") //"narlsaxxxxxxxxxxxxAtvJgAKSiQOg"
	// const CheckSum = os.Getenv("NINEPAY_CHECKSUM_KEY") //"s6KiGBywxxxxxxxxxxxxxxsx4QHM2YWzLC"
	// const Endpoint = os.Getenv("NINEPAY_ENDPOINT")

	const MerchantKey = "pxxxw"
	const SecretKey = "narlsaxxxxxxxxxxxxAtvJgAKSiQOg"
	const CheckSum = "s6KiGBywxxxxxxxxxxxxxxsx4QHM2YWzLC"
	const Endpoint = ""

	client := ninepay.New(MerchantKey, SecretKey, CheckSum, Endpoint)
	// // 2. Tạo Link Thanh Toán
	// req := models.New[models.BuildUrlRequest]()
	// req.InvoiceNo = "ORDER_001"
	// req.Amount = 50000
	// req.ReturnUrl = "https://myshop.com/result"
	// req.Description = "Mô tả thông tin đơn hàng"
	// req.Time = time.Now().Unix()
	// req.MerchantKey = MerchantKey

	// url, err := client.BuildPaymentURL(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("👉 Payment URL:", url)

	// // // 3. Kiểm tra trạng thái
	// fmt.Println("\nChecking Status for ORDER_001...")
	// transaction, err := client.Inquire("ORDER_001")
	// if err != nil {
	// 	log.Printf("❌ Check failed: %v", err)
	// } else {
	// 	fmt.Printf("✅ Status: %d \n PaymentNo: %d", *transaction.Status, *transaction.PaymentNo)
	// }

	// // // 4. Test Verify Webhook (Giả lập)
	// fmt.Println("\nTesting Verify Webhook...")

	// const data = "eyJhbW91bnQiOxxxxxm9yIjpudWxsfQ"
	// const checksum = "8FD0C7C97ACE326xxxxxxxxxB798F818B8FCB049B6A"
	// IsValid := client.VerifyChecksum(data, checksum)
	// log.Printf("❌ Check sum: %v", IsValid)

	// 5. PayerAuth
	// Build request

	RequestIdPayerAuth := fmt.Sprintf(
		"REQ_%s_%03d",
		time.Now().Format("20060102"),
		rand.Intn(1000), // 000 - 999
	)
	req := &models.PayerAuthRequest{
		RequestID: RequestIdPayerAuth,
		Currency:  "VND",
		Amount:    150000,
		ReturnURL: "https://merchant.example.com/return",
		Off3DS:    1,
	}
	req.Card = &models.CardInfo{
		CardNumber: "4456530000001005",
		CardName:   "NGUYEN VAN A",
		CardMonth:  "01",
		CardYear:   "34",
		CVV:        "111",
	}
	transaction, err := client.PayerAuth(req)
	if err != nil {
		log.Fatal("create payment error:", err)
	}
	fmt.Println("OrderCode:", transaction.OrderCode)
	fmt.Printf("transaction: %+v\n", transaction)

	// 6. Authorize

	RequestIdAuth := fmt.Sprintf(
		"REQ_%s_%03d",
		time.Now().Format("20060102"),
		rand.Intn(1000), // 000 - 999
	)
	reqAuth := &models.AuthorizeRequest{
		RequestID: RequestIdAuth,
		Currency:  "VND",
		Amount:    150000,
		OrderCode: transaction.OrderCode,
	}

	reqAuth.Card = &models.CardInfo{
		CardNumber: "4456530000001005",
		CardName:   "NGUYEN VAN A",
		CardMonth:  "01",
		CardYear:   "34",
		CVV:        "111",
	}

	transactionAuth, err := client.Authorize(reqAuth)
	if err != nil {
		log.Fatal("create authorize error:", err)
	}
	fmt.Println("OrderCode:", transactionAuth.OrderCode)
	fmt.Printf("transaction: %+v\n", transactionAuth)

	// 6. Authorize
	RequestIdCapture := fmt.Sprintf(
		"REQ_%s_%03d",
		time.Now().Format("20060102"),
		rand.Intn(1000), // 000 - 999
	)
	reqCap := &models.CaptureRequest{
		RequestID: RequestIdCapture,
		Currency:  "VND",
		Amount:    150000,
		OrderCode: transaction.OrderCode,
	}
	transactionCapture, err := client.Capture(reqCap)
	if err != nil {
		log.Fatal("Capture error:", err)
	}
	fmt.Println("OrderCode:", transactionCapture.OrderCode)
	fmt.Printf("transaction: %+v\n", transactionCapture)

	// 7. ReverseAuth
	// RequestIdReverse := fmt.Sprintf(
	// 	"REQ_%s_%03d",
	// 	time.Now().Format("20060102"),
	// 	rand.Intn(1000), // 000 - 999
	// )
	// reqReverse := &models.ReverseAuthRequest{
	// 	RequestID: RequestIdReverse,
	// 	OrderCode: transaction.OrderCode,
	// }
	// transactionReverse, err := client.ReverseAuth(reqReverse)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Reverse status=%d error=%s\n", transactionReverse.Status, transactionReverse.ErrorCode)
	// fmt.Printf("transaction: %+v\n", transactionReverse)
	//8.Refund
	RequestIdRefund := fmt.Sprintf(
		"REQ_%s_%03d",
		time.Now().Format("20060102"),
		rand.Intn(1000), // 000 - 999
	)

	reqRefund := &models.RefundCreateRequest{
		RequestID:   RequestIdRefund,
		PaymentNo:   transaction.OrderCode,
		Amount:      15000,
		Description: "Hoàn tiền đơn hàng ORDER_001",
	}

	refundRes, err := client.RefundCreate(reqRefund)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Refund status=%d refund_no=%v\n", refundRes.Status, refundRes.RefundNo)
	fmt.Printf("Refund: %+v\n", refundRes)

}
