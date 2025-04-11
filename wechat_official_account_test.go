package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func TestPushOfficialAccountMessage(t *testing.T) {
	const ()
	appId := os.Getenv("WECHAT_OFFICIAL_ACCOUNT_APP_ID")
	appSecret := os.Getenv("WECHAT_OFFICIAL_ACCOUNT_APP_SECRET")
	templateId := os.Getenv("WECHAT_OFFICIAL_ACCOUNT_TEMPLATE_ID")
	openId := os.Getenv("WECHAT_OFFICIAL_ACCOUNT_OPEN_ID")
	miniPagePath := os.Getenv("WECHAT_OFFICIAL_ACCOUNT_MINIPROGRAM_PAGE_PATH")
	miniAppId := os.Getenv("WECHAT_OFFICIAL_ACCOUNT_MINIPROGRAM_APP_ID")

	client := wechat.Wechat{}
	client.SetCache(cache.NewMemory())

	template := client.GetOfficialAccount(&config.Config{
		AppID:     appId,
		AppSecret: appSecret,
	}).GetTemplate()

	data := make(map[string]*message.TemplateDataItem)
	data["thing2"] = &message.TemplateDataItem{Value: "测试测试"}
	data["time4"] = &message.TemplateDataItem{Value: "2025-04-10 15:00:00"}
	data["thing5"] = &message.TemplateDataItem{Value: "内部变量_1.是否报警"}
	data["const8"] = &message.TemplateDataItem{Value: "Alert"}

	// data["thing2"] = &message.TemplateDataItem{Value: req.KeyWord1, Color: "#ff0000"}
	// data["time4"] = &message.TemplateDataItem{Value: req.KeyWord2, Color: "#ff0000"}
	// data["const8"] = &message.TemplateDataItem{Value: req.KeyWord4, Color: "#ff0000"}
	// data["thing5"] = &message.TemplateDataItem{Value: req.KeyWord5, Color: "#ff0000"}

	// data["keyword1"] = &message.TemplateDataItem{Value: "别动我的正式企业设备"}
	// data["keyword2"] = &message.TemplateDataItem{Value: "2025-02-18 14:00:00"}
	// data["keyword3"] = &message.TemplateDataItem{Value: "内部变量_1.是否报警"}
	// data["keyword4"] = &message.TemplateDataItem{Value: "Alert"}
	// data["keyword5"] = &message.TemplateDataItem{Value: "报警内容"}

	templateMessageReq := &message.TemplateMessage{
		ToUser:     openId,
		TemplateID: templateId,
		Data:       data,
	}

	templateMessageReq.MiniProgram.PagePath = miniPagePath
	templateMessageReq.MiniProgram.AppID = miniAppId

	_, err := template.Send(templateMessageReq)
	if err != nil {
		t.Fatalf("Failed to send template message: %s\n", err.Error())
	}

}
