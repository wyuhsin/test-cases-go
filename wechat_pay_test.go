package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

const (
	WECHATPAY_PRIVATE_KEY_PATH = "./assets/wechatpay_private_key.pem"
)

func TestWechatPayBack(t *testing.T) {
	const (
		HEADER_WECHATPAY_TIMESTAMP = "Wechatpay-Timestamp"
		HEADER_WECHATPAY_NONCE     = "Wechatpay-Nonce"
		HEADER_WECHATPAY_SERIAL    = "Wechatpay-Serial"
		HEADER_WECHATPAY_SIGNATURE = "Wechatpay-Signature"

		EVENT_TYPE_TRANSACTION_SUCCESSED = "TRANSACTION.SUCCESS"
		RESOURCE_TYPE_ENCRYPT            = "encrypt-resource"
		RESOURCE_ALGORITHM               = "AEAD_AES_256_GCM"
	)

	type WechatPayNotification struct {
		ID           string    `json:"id"`
		CreateTime   time.Time `json:"create_time"`
		ResourceType string    `json:"resource_type"`
		EventType    string    `json:"event_type"`
		Summary      string    `json:"summary"`
		Resource     struct {
			OriginalType   string `json:"original_type"`
			Algorithm      string `json:"algorithm"`       // AEAD_AES_256_GCM加密算法
			Ciphertext     string `json:"ciphertext"`      // 加密后的业务数据
			AssociatedData string `json:"associated_data"` // 附加数据(可能为空)
			Nonce          string `json:"nonce"`           // 加密使用的随机串
		} `json:"resource"`
	}

	type PaymentItem struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}

	type DiscountItem struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}

	type RiskFund struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}

	type TimeRange struct {
		StartTime       string `json:"start_time"`
		StartTimeRemark string `json:"start_time_remark"`
		EndTime         string `json:"end_time"`
		EndTimeRemark   string `json:"end_time_remark"`
	}

	type OrderLocation struct {
		StartLocation string `json:"start_location"`
		EndLocation   string `json:"end_location"`
	}

	type CollectionDetail struct {
		Seq           int    `json:"seq"`
		Amount        int    `json:"amount"`
		PaidType      string `json:"paid_type"`
		PaidTime      string `json:"paid_time"`
		TransactionID string `json:"transaction_id"`
	}

	type CollectionInfo struct {
		State        string             `json:"state"`
		TotalAmount  int                `json:"total_amount"`
		PayingAmount int                `json:"paying_amount"`
		PaidAmount   int                `json:"paid_amount"`
		Details      []CollectionDetail `json:"details"`
	}

	type OrderDetail struct {
		ServiceID           string         `json:"service_id"`
		AppID               string         `json:"appid"`
		MchID               string         `json:"mchid"`
		SubAppID            string         `json:"sub_appid"`
		SubMchID            string         `json:"sub_mchid"`
		ChannelID           string         `json:"channel_id"`
		OutOrderNo          string         `json:"out_order_no"`
		SubOpenID           string         `json:"sub_openid"`
		State               string         `json:"state"`
		ServiceIntroduction string         `json:"service_introduction"`
		TotalAmount         int            `json:"total_amount"`
		PostPayments        []PaymentItem  `json:"post_payments"`
		PostDiscounts       []DiscountItem `json:"post_discounts"`
		RiskFund            RiskFund       `json:"risk_fund"`
		TimeRange           TimeRange      `json:"time_range"`
		Location            OrderLocation  `json:"location"`
		Attach              string         `json:"attach"`
		OrderID             string         `json:"order_id"`
		NeedCollection      bool           `json:"need_collection"`
		Collection          CollectionInfo `json:"collection"`
	}

}

func TestWechatPayH5Prepay(t *testing.T) {

	const (
		H5_PREPAY_DESCRIPTION  = "HAIWELL TEST"
		H5_PREPAY_OUT_TRADE_NO = "haiwell_002"
		H5_PREPAY_NOTIFY_URL   = "https://cloud.haiwell.com/api/v1/internal/wechatNotify"
		H5_PREPAY_AMOUNT_TOTAL = 1

		H5_PREPAY_PAYER_CLIENT_IP = "112.48.26.71"
		// Type: Wap / iOS / Android
		H5_PREPAY_H5_INFO_TYPE = "Wap"
	)

	appId := os.Getenv("WECHAT_PAY_APP_ID")
	mchId := os.Getenv("WECHAT_PAY_MCH_ID")
	serialNumber := os.Getenv("WECHAT_PAY_MCH_CERTIFICATE_SERIAL_NUMBER")
	apiKey := os.Getenv("WECHAT_PAY_MCH_API_V3_KEY")

	var (
		ctx context.Context = context.Background()
	)

	pk, err := utils.LoadPrivateKeyWithPath(WECHATPAY_PRIVATE_KEY_PATH)
	if err != nil {
		t.Fatalf("Failed to load merchant private key: %s\n", err.Error())
	}

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchId, serialNumber, pk, apiKey),
	}

	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		t.Fatalf("Failed to new wechat pay client: %s\n", err.Error())
	}

	svc := h5.H5ApiService{Client: client}
	resp, result, err := svc.Prepay(ctx, h5.PrepayRequest{
		Appid:       core.String(appId),
		Mchid:       core.String(mchId),
		Description: core.String(H5_PREPAY_DESCRIPTION),
		OutTradeNo:  core.String(H5_PREPAY_OUT_TRADE_NO),
		NotifyUrl:   core.String(H5_PREPAY_NOTIFY_URL),
		Amount: &h5.Amount{
			Total: core.Int64(H5_PREPAY_AMOUNT_TOTAL),
		},
		SceneInfo: &h5.SceneInfo{
			PayerClientIp: core.String(H5_PREPAY_PAYER_CLIENT_IP),
			H5Info: &h5.H5Info{
				Type: core.String(H5_PREPAY_H5_INFO_TYPE),
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create prepay order: %s\n", err.Error())
	}

	t.Logf("status=%d resp=%s", result.Response.StatusCode, resp)
}

func TestWechatPayAppPrepay(t *testing.T) {

	const (
		APP_PREPAY_DESCRIPTION  = "HAIWELL TEST"
		APP_PREPAY_OUT_TRADE_NO = "haiwell_002"
		APP_PREPAY_NOTIFY_URL   = "https://cloud.haiwell.com/api/v1/internal/wechatNotify"
		APP_PREPAY_AMOUNT_TOTAL = 1
	)

	var (
		ctx context.Context = context.Background()
	)

	appId := os.Getenv("WECHAT_PAY_APP_ID")
	mchId := os.Getenv("WECHAT_PAY_MCH_ID")
	serialNumber := os.Getenv("WECHAT_PAY_MCH_CERTIFICATE_SERIAL_NUMBER")
	apiKey := os.Getenv("WECHAT_PAY_MCH_API_V3_KEY")

	pk, err := utils.LoadPrivateKeyWithPath(WECHATPAY_PRIVATE_KEY_PATH)
	if err != nil {
		t.Fatalf("Failed to load merchant private key: %s\n", err.Error())
	}

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchId, serialNumber, pk, apiKey),
	}

	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		t.Fatalf("Failed to new wechat pay client: %s\n", err.Error())
	}

	svc := app.AppApiService{Client: client}
	resp, result, err := svc.Prepay(ctx, app.PrepayRequest{
		Appid:       core.String(appId),
		Mchid:       core.String(mchId),
		Description: core.String(APP_PREPAY_DESCRIPTION),
		OutTradeNo:  core.String(APP_PREPAY_OUT_TRADE_NO),
		NotifyUrl:   core.String(APP_PREPAY_NOTIFY_URL),
		Amount: &app.Amount{
			Total: core.Int64(APP_PREPAY_AMOUNT_TOTAL),
		},
	})
	if err != nil {
		t.Fatalf("Failed to create prepay order: %s\n", err.Error())
	}

	t.Logf("status=%d resp=%s", result.Response.StatusCode, resp)
}

func TestWechatPayNativePrepay(t *testing.T) {

	const (
		NATIVE_PREPAY_DESCRIPTION  = "HAIWELL TEST"
		NATIVE_PREPAY_OUT_TRADE_NO = "haiwell_001"
		NATIVE_PREPAY_NOTIFY_URL   = "https://cloud.haiwell.com/api/v1/internal/wechatNotify"
		NATIVE_PREPAY_AMOUNT_TOTAL = 1
	)

	var (
		ctx context.Context = context.Background()
	)

	appId := os.Getenv("WECHAT_PAY_APP_ID")
	mchId := os.Getenv("WECHAT_PAY_MCH_ID")
	serialNumber := os.Getenv("WECHAT_PAY_MCH_CERTIFICATE_SERIAL_NUMBER")
	apiKey := os.Getenv("WECHAT_PAY_MCH_API_V3_KEY")

	pk, err := utils.LoadPrivateKeyWithPath(WECHATPAY_PRIVATE_KEY_PATH)
	if err != nil {
		t.Fatalf("Failed to load merchant private key: %s\n", err.Error())
	}

	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchId, serialNumber, pk, apiKey),
	}

	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		t.Fatalf("Failed to new wechat pay client: %s\n", err.Error())
	}

	svc := native.NativeApiService{Client: client}

	resp, result, err := svc.Prepay(ctx, native.PrepayRequest{
		Appid:       core.String(appId),
		Mchid:       core.String(mchId),
		Description: core.String(NATIVE_PREPAY_DESCRIPTION),
		OutTradeNo:  core.String(NATIVE_PREPAY_OUT_TRADE_NO),
		NotifyUrl:   core.String(NATIVE_PREPAY_NOTIFY_URL),
		Amount: &native.Amount{
			Total: core.Int64(NATIVE_PREPAY_AMOUNT_TOTAL),
		},
	})
	if err != nil {
		t.Fatalf("Failed to create NATIVE_PREPAY order: %s\n", err.Error())
	}

	t.Logf("status=%d resp=%s", result.Response.StatusCode, resp)
}
