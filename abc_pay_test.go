package tests

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	// "golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"software.sslmate.com/src/go-pkcs12"
)

type Response struct {
	Msg *Msg `json:"MSG"`
}

type Msg struct {
	Message            *ResponseMessage `json:"Message"`
	SignatureAlgorithm string           `json:"Signature-Algorithm"`
	Signature          string           `json:"Signature"`
}

type ResponseMessage struct {
	Version      string    `json:"Version,omitempty"`
	Format       string    `json:"Format,omitempty"`
	Merchant     *Merchant `json:"Merchant,omitempty"`
	ReturnCode   string    `json:"ReturnCode,omitempty"`
	ErrorMessage string    `json:"ErrorMessage,omitempty"`
	TrxType      string    `json:"TrxType,omitempty"`
	OrderNo      string    `json:"OrderNo,omitempty"`
	PaymentURL   string    `json:"PaymentURL,omitempty"`
	OrderAmount  string    `json:"OrderAmount,omitempty"`
}

type Request struct {
	Message            *Message `json:"Message"`
	SignatureAlgorithm string   `json:"Signature-Algorithm"`
	Signature          string   `json:"Signature"`
}
type Message struct {
	Version    string      `json:"Version,omitempty"`
	Format     string      `json:"Format,omitempty"`
	Merchant   *Merchant   `json:"Merchant,omitempty"`
	TrxRequest *TrxRequest `json:"TrxRequest,omitempty"`
}

type Merchant struct {
	ECMerchantType string `json:"ECMerchantType,omitempty"`
	MerchantID     string `json:"MerchantID,omitempty"`
}

type TrxRequest struct {
	TrxType          string `json:"TrxType,omitempty"`
	Order            *Order `json:"Order,omitempty"`
	PaymentType      string `json:"PaymentType,omitempty"`
	PaymentLinkType  string `json:"PaymentLinkType,omitempty"`
	ReceiveAccount   string `json:"ReceiveAccount,omitempty"`
	ReceiveAccName   string `json:"ReceiveAccName,omitempty"`
	NotifyType       string `json:"NotifyType,omitempty"`
	ResultNotifyURL  string `json:"ResultNotifyURL,omitempty"`
	MerchantRemarks  string `json:"MerchantRemarks,omitempty"`
	IsBreakAccount   string `json:"IsBreakAccount,omitempty"`
	SplitAccTemplate string `json:"SplitAccTemplate,omitempty"`
}

type Order struct {
	PayTypeID         string       `json:"PayTypeID,omitempty"`
	OrderDate         string       `json:"OrderDate,omitempty"`
	OrderTime         string       `json:"OrderTime,omitempty"`
	OrderTimeoutDate  string       `json:"orderTimeoutDate,omitempty"`
	OrderNo           string       `json:"OrderNo,omitempty"`
	CurrencyCode      string       `json:"CurrencyCode,omitempty"`
	OrderAmount       string       `json:"OrderAmount,omitempty"`
	SubsidyAmount     string       `json:"SubsidyAmount,omitempty"`
	Fee               string       `json:"Fee,omitempty"`
	AccountNo         string       `json:"AccountNo,omitempty"`
	OrderDesc         string       `json:"OrderDesc,omitempty"`
	OrderURL          string       `json:"OrderURL,omitempty"`
	ReceiverAddress   string       `json:"ReceiverAddress,omitempty"`
	InstallmentMark   string       `json:"InstallmentMark,omitempty"`
	CommodityType     string       `json:"CommodityType,omitempty"`
	BuyIP             string       `json:"BuyIP,omitempty"`
	ExpiredDate       string       `json:"ExpiredDate,omitempty"`
	SplitAccInfoItems string       `json:"SplitAccInfoItems,omitempty"`
	OrderItems        []*OrderItem `json:"OrderItems,omitempty"`
}

type OrderItem struct {
	SubMerName         string `json:"SubMerName,omitempty"`
	SubMerId           string `json:"SubMerId,omitempty"`
	SubMerMCC          string `json:"SubMerMCC,omitempty"`
	SubMerchantRemarks string `json:"SubMerchantRemarks,omitempty"`
	ProductID          string `json:"ProductID,omitempty"`
	ProductName        string `json:"ProductName,omitempty"`
	UnitPrice          string `json:"UnitPrice,omitempty"`
	Qty                string `json:"Qty,omitempty"`
	ProductRemarks     string `json:"ProductRemarks,omitempty"`
	ProductType        string `json:"ProductType,omitempty"`
	ProductDiscount    string `json:"ProductDiscount,omitempty"`
	ProductExpiredDate string `json:"ProductExpiredDate,omitempty"`
}

const (
	SIGNATURE_ALGORITHM = "SHA1withRSA"

	ABC_PAY_MECHANT_CERTIFICATE_PATH = "./assets/abc_pay_merchant_cert.pfx"
	// ABC_PAY_MECHANT_CERTIFICATE_PATH   = "./assets/abc_pay_merchant_cert_haiwell.pfx"
	ABC_PAY_TRUST_PAY_CERTIFICATE_PATH = "./assets/abc_pay_trust_pay.cer"
	ABC_PAY_GATEWAY_URL                = "https://pay.test.abchina.com/ebusold/trustpay/ReceiveMerchantTrxReqServlet"
	ABC_PAY_NOTIFY_URL                 = "https://google.com"

	ABC_PAY_REQUEST_VERSION = "V3.0.0"
	ABC_PAY_REQUEST_FORMAT  = "JSON"
	ABC_PAY_MERCHANT_TYPE   = "EBUS"

	ABC_PAY_TRX_TYPE_PAY_RER         = "PayReq"
	ABC_PAY_ORDER_PAY_TYPE_IMMEDIATE = "ImmediatePay"
)

var (
	privateKeyPassword = os.Getenv("ABC_PAY_PRIVATE_KAY_PASSWORD")
	merchantId         = os.Getenv("ABC_PAY_MERCHANT_ID")
	ip                 = os.Getenv("PUBLIC_IP")
)

func TestAbcPayPayReq(t *testing.T) {

	now := time.Now()
	date := now.Format("2005/01/02")
	tim := now.Format("15:04:05")
	internalOrderNo := "TEST-20250417141600"
	price := "0.01"
	productName := "TEST-PRODUCT"

	merchant := &Merchant{
		MerchantID:     merchantId,
		ECMerchantType: ABC_PAY_MERCHANT_TYPE,
	}

	orderItems := []*OrderItem{
		{
			ProductName: productName,
		},
	}

	order := &Order{
		PayTypeID:   ABC_PAY_ORDER_PAY_TYPE_IMMEDIATE,
		OrderDate:   date,
		OrderTime:   tim,
		OrderNo:     internalOrderNo,
		OrderAmount: price,
		BuyIP:       ip,
		OrderItems:  orderItems,
	}

	trxRequest := &TrxRequest{
		TrxType:         ABC_PAY_TRX_TYPE_PAY_RER,
		PaymentType:     "1",
		PaymentLinkType: "1",
		NotifyType:      "0",
		ResultNotifyURL: ABC_PAY_NOTIFY_URL,
		Order:           order,
	}

	message := &Message{
		Version:    ABC_PAY_REQUEST_VERSION,
		Format:     ABC_PAY_REQUEST_FORMAT,
		Merchant:   merchant,
		TrxRequest: trxRequest,
	}

	merchantCertificate, err := os.ReadFile(ABC_PAY_MECHANT_CERTIFICATE_PATH)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	signature, err := calculateSignature(merchantCertificate, privateKeyPassword, messageBytes)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	request := Request{
		SignatureAlgorithm: SIGNATURE_ALGORITHM,
		Message:            message,
		Signature:          signature,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	resp, err := http.Post(ABC_PAY_GATEWAY_URL, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	defer resp.Body.Close()

	response := Response{}
	if err := json.Unmarshal(respBytes, &response); err != nil {
		t.Fatalf("%s\n", err)
	}

	responseMessageBytes, err := json.Marshal(response.Msg.Message)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	t.Log(string(responseMessageBytes))

	// publicCertificate, err := os.ReadFile(ABC_PAY_TRUST_PAY_CERTIFICATE_PATH)
	// if err != nil {
	// 	t.Fatalf("%s\n", err)
	// }

	// corr, err := verifyResponse(response.Msg.Signature, responseMessageBytes, publicCertificate)
	// if err != nil {
	// 	t.Fatalf("%s\n", err)
	// }

	// assert.True(t, corr)
}

func encodeToGBK(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewEncoder())
	return io.ReadAll(reader)
}

func TestAbcPaySignRequest(t *testing.T) {

	data := []byte("{\"Version\":\"V3.0.0\",\"Format\":\"JSON\",\"Merchant\":{\"ECMerchantType\":\"EBUS\",\"MerchantID\":\"103882200000958\"},\"TrxRequest\":{\"TrxType\":\"PayReq\",\"Order\":{\"PayTypeID\":\"ImmediatePay\",\"OrderDate\":\"2021/02/04\",\"OrderTime\":\"16:36:18\",\"orderTimeoutDate\":\"20211231000000\",\"OrderNo\":\"ON2021456440301001\",\"CurrencyCode\":\"156\",\"OrderAmount\":\"1.00\",\"SubsidyAmount\":\"1.00\",\"Fee\":\"\",\"AccountNo\":\"\",\"OrderDesc\":\"\",\"OrderURL\":\"\",\"ReceiverAddress\":\"北京\",\"InstallmentMark\":\"0\",\"CommodityType\":\"0101\",\"BuyIP\":\"127.0.0.1\",\"ExpiredDate\":\"30\",\"SplitAccInfoItems\":\"\",\"OrderItems\":{\"SubMerName\":\"\",\"SubMerId\":\"\",\"SubMerMCC\":\"\",\"SubMerchantRemarks\":\"\",\"ProductID\":\"\",\"ProductName\":\"中国移动IP卡\",\"UnitPrice\":\"\",\"Qty\":\"\",\"ProductRemarks\":\"\",\"ProductType\":\"\",\"ProductDiscount\":\"\",\"ProductExpiredDate\":\"\"}},\"PaymentType\":\"A\",\"PaymentLinkType\":\"1\",\"ReceiveAccount\":\"\",\"ReceiveAccName\":\"\",\"NotifyType\":\"0\",\"ResultNotifyURL\":\"http://yourwebsite/appname/MerchantResult.jsp\",\"MerchantRemarks\":\"\",\"IsBreakAccount\":\"0\",\"SplitAccTemplate\":\"\"}}")

	// expected := "{\"Message\":{\"Version\":\"V3.0.0\",\"Format\":\"JSON\",\"Merchant\":{\"ECMerchantType\":\"EBUS\",\"MerchantID\":\"103882200000958\"},\"TrxRequest\":{\"TrxType\":\"PayReq\",\"Order\":{\"PayTypeID\":\"ImmediatePay\",\"OrderDate\":\"2021/02/04\",\"OrderTime\":\"16:36:18\",\"orderTimeoutDate\":\"20211231000000\",\"OrderNo\":\"ON2021456440301001\",\"CurrencyCode\":\"156\",\"OrderAmount\":\"1.00\",\"SubsidyAmount\":\"1.00\",\"Fee\":\"\",\"AccountNo\":\"\",\"OrderDesc\":\"\",\"OrderURL\":\"\",\"ReceiverAddress\":\"北京\",\"InstallmentMark\":\"0\",\"CommodityType\":\"0101\",\"BuyIP\":\"127.0.0.1\",\"ExpiredDate\":\"30\",\"SplitAccInfoItems\":\"\",\"OrderItems\":{\"SubMerName\":\"\",\"SubMerId\":\"\",\"SubMerMCC\":\"\",\"SubMerchantRemarks\":\"\",\"ProductID\":\"\",\"ProductName\":\"中国移动IP卡\",\"UnitPrice\":\"\",\"Qty\":\"\",\"ProductRemarks\":\"\",\"ProductType\":\"\",\"ProductDiscount\":\"\",\"ProductExpiredDate\":\"\"}},\"PaymentType\":\"A\",\"PaymentLinkType\":\"1\",\"ReceiveAccount\":\"\",\"ReceiveAccName\":\"\",\"NotifyType\":\"0\",\"ResultNotifyURL\":\"http://yourwebsite/appname/MerchantResult.jsp\",\"MerchantRemarks\":\"\",\"IsBreakAccount\":\"0\",\"SplitAccTemplate\":\"\"}},\"Signature-Algorithm\":\"SHA1withRSA\",\"Signature\":\"FsWwYCwn/vRUt66j5FmvSBlTC2rfagNOlzom60ISEy9TSJus0+lJL/PxsxEiEvQV8jcbM1NDciwk4ffIQl3nnmqVcvHpF2JNXQWev19ELARfukJLsUCmVZuAVW8Na4K0yCvfEZDdc5w/ju+EnZulhgjwYVb/a5JHayzkidINBTM=\"}"

	expected := "FsWwYCwn/vRUt66j5FmvSBlTC2rfagNOlzom60ISEy9TSJus0+lJL/PxsxEiEvQV8jcbM1NDciwk4ffIQl3nnmqVcvHpF2JNXQWev19ELARfukJLsUCmVZuAVW8Na4K0yCvfEZDdc5w/ju+EnZulhgjwYVb/a5JHayzkidINBTM="

	merchantCertificate, err := os.ReadFile(ABC_PAY_MECHANT_CERTIFICATE_PATH)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	signature, err := calculateSignature(merchantCertificate, privateKeyPassword, data)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	assert.Equal(t, expected, signature)
}

func TestAbcPayVerifyResponse(t *testing.T) {
	data := []byte("{\"Version\":\"V3.0.0\",\"Format\":\"JSON\",\"Merchant\":{\"ECMerchantType\":\"EBUS\",\"MerchantID\":\"103882200000958\"},\"ReturnCode\":\"0000\",\"ErrorMessage\":\"交易成功\",\"TrxType\":\"PayReq\",\"OrderNo\":\"ON2021456440301001\",\"PaymentURL\":\"https://pay.test.abchina.com/perbankold/PaymentModeNewAct.ebf?TOKEN=16124277721648016933\",\"OrderAmount\":\"1.00\",\"OneQRForAll\":\"http://mpay.test.abchina.com/mpay/mobileBank/zh_CN/EBusinessModule/BarcodeH5Act.aspx?token=16124277721648016933\"}")

	signatureBase64 := "J1vXDnsTgbgmnpS/yLzu2m94A82mmva+P+2oGX52gqoV7CS2QWdLBqXf7uz/5P6Obq4ow0H7rraT1YA3xNd1FQbuOZCrwEx61yaEJbMludbKjhtm/B8dXcmPqnW+DzOzcuGr2yU8vCMt8DEH0rouei6q3AugOatV6NCf2bTTMyM="

	certificate, err := os.ReadFile(ABC_PAY_TRUST_PAY_CERTIFICATE_PATH)
	if err != nil {
		t.Fatal(err.Error())
	}

	ok, err := verifyResponse(signatureBase64, data, certificate)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.True(t, ok)
}

func calculateSignature(p12Data []byte, password string, message []byte) (string, error) {
	pk, _, err := extractPrivateKey(p12Data, password)
	if err != nil {
		return "", err
	}

	hashed := sha1.Sum([]byte(message))

	signature, err := rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	return signatureBase64, nil

	// signed := SignedRequest{
	// 	Message:            message,
	// 	SignatureAlgorithm: SIGNATURE_ALGORITHM,
	// 	Signature:          signatureBase64,
	// }

	// jsonBytes, err := json.Marshal(signed)
	// if err != nil {
	// 	return "", err
	// }

	// return string(jsonBytes), nil
}

func extractPrivateKey(p12Data []byte, password string) (*rsa.PrivateKey, *x509.Certificate, error) {
	privKey, cert, err := pkcs12.Decode(p12Data, password)
	if err != nil {
		return nil, nil, err
	}

	rsaKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("private key is not RSA")
	}
	return rsaKey, cert, nil
}

func verifyResponse(signatureBase64 string, message, trustPayCert []byte) (bool, error) {
	cert, err := parseX509Cert(trustPayCert)
	if err != nil {
		return false, err
	}

	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("public key is not rsa type")
	}

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, err
	}

	gbkEncoder := simplifiedchinese.GBK.NewEncoder()
	msgBytes, err := gbkEncoder.Bytes(message)
	if err != nil {
		fmt.Println("aaaa")
		return false, err
	}

	hash := sha1.Sum(msgBytes)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hash[:], signature)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func parseX509Cert(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block != nil {
		return x509.ParseCertificate(block.Bytes)
	}

	return x509.ParseCertificate(data)
}
