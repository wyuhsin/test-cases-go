package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestMobSendSms(t *testing.T) {
	const (
		SMS_VERIFY_URL = "https://webapi.sms.mob.com/sms/verify"
		SMS_SEND_URL   = "https://webapi.sms.mob.com/sms/sendmsg"

		SMS_SEND_PHONE        = ""
		SMS_SEND_COUNTRY_CODE = "86"
	)

	appKey := os.Getenv("MOB_APP_KEY")

	var (
		resultMap map[string]any
	)

	resp, err := http.Post(SMS_SEND_URL, "text/html", strings.NewReader(url.Values{
		"phone":  {SMS_SEND_PHONE},
		"appkey": {appKey},
		"zone":   {SMS_SEND_COUNTRY_CODE},
	}.Encode()))
	if err != nil {
		t.Fatalf("Failed to send http request: %s\n", err.Error())
	}
	defer resp.Body.Close()

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read http response: %s\n", err.Error())
	}

	if err := json.Unmarshal(buff, &resultMap); err != nil {
		t.Fatalf("Failed to unmarshal json: %s\n", err.Error())
	}

	t.Logf("Http response: %#v\n", resultMap)
}
