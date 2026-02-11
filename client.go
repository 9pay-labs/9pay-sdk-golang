package ninepay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/9pay-labs/9pay-sdk-golang/pkg/config"
	"github.com/9pay-labs/9pay-sdk-golang/pkg/consts"
	"github.com/9pay-labs/9pay-sdk-golang/pkg/models"
	"github.com/9pay-labs/9pay-sdk-golang/pkg/security"
	"github.com/9pay-labs/9pay-sdk-golang/pkg/utils"
)

type Client struct {
	Config     *config.Config
	Signer     *security.Signer
	HttpClient *http.Client
}

func New(key, secret string, checksum string, endpoint string) *Client {
	cfg := config.New(key, secret, checksum, endpoint)
	return &Client{
		Config:     cfg,
		Signer:     security.NewSigner(cfg.SecretKey, cfg.CheckSumKey),
		HttpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) BuildPaymentURL(req *models.BuildUrlRequest) (string, error) {

	if err := req.Validate(); err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}
	ts := time.Now().Unix()
	req.Set("merchantKey", c.Config.MerchantKey)
	req.Set("time", ts)
	dataMap := req.ToMap()
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(dataMap); err != nil {
		return "", err
	}

	jsonBytes := bytes.TrimSpace(buf.Bytes())

	canonicalPayload := utils.BuildCanonicalParams(dataMap)

	fullUrl := c.Config.Endpoint + consts.PathPaymentCreate
	stringToSign := fmt.Sprintf("POST\n%s\n%d\n%s", fullUrl, ts, canonicalPayload)
	signature := c.Signer.Sign(stringToSign)
	base64Data := base64.StdEncoding.EncodeToString(jsonBytes)
	v := url.Values{}
	v.Set("baseEncode", base64Data)
	v.Set("signature", signature)

	return fmt.Sprintf("%s/portal?%s", c.Config.Endpoint, v.Encode()), nil
}

func (c *Client) Inquire(invoiceNo string) (*models.TransactionInquire, error) {
	path := fmt.Sprintf(consts.PathInquire, invoiceNo)

	req := models.New[models.InquireRequest]()

	var res models.ResponseInquire
	// Check docs 9Pay để lấy đúng path check status
	err := c.callAPI("GET", path, req, &res)
	if err != nil {
		return nil, err
	}

	if res.InvoiceNo == nil {
		return nil, fmt.Errorf(
			"9pay inquire failed (%s): %s",
			res.ErrorCode,
			res.FailureReason,
		)
	}

	return &res.TransactionInquire, nil
}
func (c *Client) VerifyChecksum(data, checksum string) bool {
	if c.Signer == nil {
		return false
	}
	return c.Signer.VerifyChecksum(data, checksum)
}

func (c *Client) PayerAuth(req *models.PayerAuthRequest) (*models.PayerAuthResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	path := consts.MchPayerAuth
	var res models.PayerAuthResponse
	if err := c.callAPIMch("POST", path, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) Authorize(req *models.AuthorizeRequest) (*models.AuthorizeResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	path := consts.MchAuthorize

	var res models.AuthorizeResponse
	if err := c.callAPIMch("POST", path, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) Capture(req *models.CaptureRequest) (*models.CaptureResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	path := consts.MchCapture

	var res models.CaptureResponse
	if err := c.callAPIMch("POST", path, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) ReverseAuth(
	req *models.ReverseAuthRequest,
) (*models.ReverseAuthResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	path := consts.MchReverseAuth

	var res models.ReverseAuthResponse
	if err := c.callAPIMch("POST", path, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) RefundCreate(
	req *models.RefundCreateRequest,
) (*models.RefundCreateResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	path := consts.MchRefundCreate

	var res models.RefundCreateResponse
	if err := c.callAPI("POST", path, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) callAPI(method, path string, req models.IBaseRequest, res interface{}) error {
	ts := time.Now().Unix()
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	dataMap := req.ToMap()
	if err := enc.Encode(dataMap); err != nil {
		return err
	}
	jsonBytes := bytes.TrimSpace(buf.Bytes())
	canonicalPayload := utils.BuildCanonicalParams(dataMap)
	fullUrl := c.Config.Endpoint + path
	stringToSign := fmt.Sprintf("%s\n%s\n%d", method, fullUrl, ts)
	if canonicalPayload != "" {
		stringToSign += "\n" + canonicalPayload
	}
	signature := c.Signer.Sign(stringToSign)
	httpReq, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Date", fmt.Sprintf("%d", ts))
	authHeader := fmt.Sprintf(
		"Signature Algorithm=HS256,Credential=%s,SignedHeaders=,Signature=%s",
		c.Config.MerchantKey,
		signature,
	)
	httpReq.Header.Set("Authorization", authHeader)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("http error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return json.Unmarshal(bodyBytes, res)
}

func (c *Client) callAPIMch(method, path string, req models.IBaseRequest, res interface{}) error {
	ts := time.Now().Unix()

	dataMap := req.ToMap()

	jsonBytes, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	encodedJSON := url.QueryEscape(string(jsonBytes))
	canonicalPayload := "json=" + encodedJSON

	fullUrl := c.Config.Endpoint + path

	stringToSign := fmt.Sprintf(
		"%s\n%s\n%d\n%s",
		method,
		fullUrl,
		ts,
		canonicalPayload,
	)

	signature := c.Signer.Sign(stringToSign)

	httpReq, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Date", fmt.Sprintf("%d", ts))
	httpReq.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"Signature Algorithm=HS256,Credential=%s,SignedHeaders=,Signature=%s",
			c.Config.MerchantKey,
			signature,
		),
	)

	// 6. Send request
	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("http error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return json.Unmarshal(bodyBytes, res)
}
