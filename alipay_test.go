package tests

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
)

const (
	ALIPAY_PRIVATE_KEY_PATH = "./assets/alipay_private_key.pem"
	ALIPAY_OPENAPI_GATEWAY  = "https://openapi.alipay.com/gateway.do"
)

func TestAlipayTradeAppPay(t *testing.T) {
	const (
		TRADE_APP_PAY_SUBJECT      = "iPhone16 Pro Max"
		TRADE_APP_PAY_OUT_TRADE_NO = "9876543211"
		TRADE_APP_PAY_TOTAL_AMOUNT = "0.01"
		TRADE_APP_PAY_PRODUCT_CODE = "QUICK_MEDICAL_PAY"
		TRADE_APP_PAY_NOTIFY_URL   = "https://cloud.haiwell.com/api/v1/alipay/notify"

		TRADE_APP_PAY_GOODS_ID       = "iotcard"
		TRADE_APP_PAY_GOODS_NAME     = "IoT Card"
		TRADE_APP_PAY_GOODS_QUANTITY = 1
		TRADE_APP_PAY_GOODS_PRICE    = 0.01
	)
	appId := os.Getenv("ALIPAY_APP_ID")

	ap := alipay.TradeAppPay{}
	ap.Subject = TRADE_APP_PAY_SUBJECT
	ap.OutTradeNo = TRADE_APP_PAY_OUT_TRADE_NO
	ap.TotalAmount = TRADE_APP_PAY_TOTAL_AMOUNT
	ap.ProductCode = TRADE_APP_PAY_PRODUCT_CODE
	ap.NotifyURL = TRADE_APP_PAY_NOTIFY_URL

	goodsDetail := &alipay.GoodsDetail{
		GoodsId:   TRADE_APP_PAY_GOODS_ID,
		GoodsName: TRADE_APP_PAY_GOODS_NAME,
		Quantity:  TRADE_APP_PAY_GOODS_QUANTITY,
		Price:     TRADE_APP_PAY_GOODS_PRICE,
	}
	ap.GoodsDetail = []*alipay.GoodsDetail{goodsDetail}

	privateKeyFile, err := os.OpenFile(ALIPAY_PRIVATE_KEY_PATH, os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open alipay private key file: %s\n", err.Error())
	}

	pk, err := io.ReadAll(privateKeyFile)
	if err != nil {
		t.Fatalf("Failed to read file content: %s\n", err.Error())
	}

	client, err := alipay.New(appId, string(pk), true)
	result, err := client.TradeAppPay(ap)
	if err != nil {
		t.Fatalf("Failed to generate trade app pay: %s\n", err.Error())
	}

	t.Logf("Trade app pay result: %s\n", result)
}

func TestAlipayTradePagePay(t *testing.T) {
	const (
		TRADE_PAGE_PAY_SUBJECT       = "iPhone16 Pro Max"
		TRADE_PAGE_PAY_OUT_TRADE_NO  = "9876543212"
		TRADE_PAGE_PAY_TOTAL_AMOUNT  = "0.01"
		TRADE_PAGE_PAY_PRODUCT_CODE  = "FAST_INSTANT_TRADE_PAY"
		TRADE_PAGE_PAY_NOTIFY_URL    = "https://cloud.haiwell.com/api/v1/alipay/notify"
		TRADE_PAGE_PAY_MODE_QR_CODE  = "4"
		TRADE_PAGE_PAY_QR_CODE_WIDTH = "100"

		TRADE_PAGE_PAY_GOODS_ID       = "iotcard"
		TRADE_PAGE_PAY_GOODS_NAME     = "IoT Card"
		TRADE_PAGE_PAY_GOODS_QUANTITY = 1
		TRADE_PAGE_PAY_GOODS_PRICE    = 0.01
	)

	appId := os.Getenv("ALIPAY_APP_ID")

	ap := alipay.TradePagePay{}
	ap.Subject = TRADE_PAGE_PAY_SUBJECT
	ap.OutTradeNo = TRADE_PAGE_PAY_OUT_TRADE_NO
	ap.TotalAmount = TRADE_PAGE_PAY_TOTAL_AMOUNT
	ap.ProductCode = TRADE_PAGE_PAY_PRODUCT_CODE
	ap.NotifyURL = TRADE_PAGE_PAY_NOTIFY_URL
	ap.QRPayMode = TRADE_PAGE_PAY_MODE_QR_CODE
	ap.QRCodeWidth = TRADE_PAGE_PAY_QR_CODE_WIDTH

	goodsDetail := &alipay.GoodsDetail{
		GoodsId:   TRADE_PAGE_PAY_GOODS_ID,
		GoodsName: TRADE_PAGE_PAY_GOODS_NAME,
		Quantity:  TRADE_PAGE_PAY_GOODS_QUANTITY,
		Price:     TRADE_PAGE_PAY_GOODS_PRICE,
	}
	ap.GoodsDetail = []*alipay.GoodsDetail{goodsDetail}

	privateKeyFile, err := os.OpenFile(ALIPAY_PRIVATE_KEY_PATH, os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open alipay private key file: %s\n", err.Error())
	}

	pk, err := io.ReadAll(privateKeyFile)
	if err != nil {
		t.Fatalf("Failed to read file content: %s\n", err.Error())
	}

	client, err := alipay.New(appId, string(pk), true)
	result, err := client.TradePagePay(ap)
	if err != nil {
		t.Fatalf("Failed to generate trade app pay: %s\n", err.Error())
	}

	t.Logf("Trade app pay result: %s\n", result)
}

func TestAlipayTradePrecreate(t *testing.T) {
	const (
		TRADE_PRECREATE_METHOD    = "alipay.trade.precreate"
		TRADE_PRECREATE_CHARSET   = "UTF-8"
		TRADE_PRECREATE_SIGN_TYPE = "RSA2"
		TRADE_PRECREATE_TIMESTAMP = "2025-01-07 16:05:10"
		TRADE_PRECREATE_VERSION   = "1.0"
		TRADE_PRECREATE_FORMAT    = "json"

		TRADE_PRECREATE_OUT_TRADE_NO = "987654321"
		TRADE_PRECREATE_TOTAL_AMOUNT = "0.01"
		TRADE_PRECREATE_SUBJECT      = "iPhone16 Pro Max"
		TRADE_PRECREATE_PRODUCT_CODE = "QR_CODE_OFFLINE"
	)
	appId := os.Getenv("ALIPAY_APP_ID")

	bizContent := map[string]any{
		"out_trade_no": TRADE_PRECREATE_OUT_TRADE_NO,
		"total_amount": TRADE_PRECREATE_TOTAL_AMOUNT,
		"subject":      TRADE_PRECREATE_SUBJECT,
		"product_code": TRADE_PRECREATE_PRODUCT_CODE,
	}
	bizContentJSON, _ := json.Marshal(bizContent)

	params := map[string]string{
		"app_id":      appId,
		"method":      TRADE_PRECREATE_METHOD,
		"charset":     TRADE_PRECREATE_CHARSET,
		"sign_type":   TRADE_PRECREATE_SIGN_TYPE,
		"timestamp":   TRADE_PRECREATE_TIMESTAMP,
		"version":     TRADE_PRECREATE_VERSION,
		"format":      TRADE_PRECREATE_FORMAT,
		"biz_content": string(bizContentJSON),
	}

	privateKeyFile, err := os.Open(ALIPAY_PRIVATE_KEY_PATH)
	if err != nil {
		t.Fatalf("Failed to open alipay private key file: %s\n", err.Error())
	}

	pk, err := io.ReadAll(privateKeyFile)
	if err != nil {
		t.Fatalf("Failed to read file content: %s\n", err.Error())
	}

	// Generate signature
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sb := strings.Builder{}
	values := url.Values{}

	for i, k := range keys {
		if params[k] != "" {
			if i > 0 {
				sb.WriteString("&")
			}
			sb.WriteString(k + "=" + params[k])
		}

		values.Set(k, params[k])
	}

	signature, err := generateRSA2048Signature(sb.String(), string(pk))
	if err != nil {
		t.Fatalf("Failed to generate signature: %s\n", err.Error())
	}

	values.Set("sign", signature)

	url := fmt.Sprintf("%s?%s", ALIPAY_OPENAPI_GATEWAY, values.Encode())

	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatalf("Failed to send http request: %s\n", err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read http response: %s\n", err.Error())
	}

	t.Logf("HTTP response: %s\n", bodyBytes)
}

func parseRSA2048PrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, fmt.Errorf("Failed to decode private key")
	}

	// Resolve PKCS1
	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return parsedKey, nil
	}

	// Resolve PKCS8
	parsedKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		if key, ok := parsedKeyInterface.(*rsa.PrivateKey); ok {
			return key, nil
		}
	}
	return nil, fmt.Errorf("Failed to parse private key: unsupported format")
}

// Generate RSA2048 sign
func generateRSA2048Signature(s string, privateKey string) (string, error) {

	privateKeyParsed, err := parseRSA2048PrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	hashed := crypto.SHA256.New()
	hashed.Write([]byte(s))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKeyParsed, crypto.SHA256, hashed.Sum(nil))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
