package controllers

import (
	"fmt"
	"sports-common/request"
	models "sports-models"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

/// 获取支付渠道信息
var getPaymentChannels = func(platform string) string {
	channels := models.PaymentThirds.GetRelationChannels(platform)
	result := []string{}
	hasElectricBank := false
	for _, channel := range channels {
		// 因为网银可选银行列表过多, 所以把网银缩成一条支付渠道
		if channel.ChannelType == "ebank" && !hasElectricBank {
			code := channel.PayList[0]["code"]
			result = append(result, fmt.Sprintf("%s-ebank|%s - %s", code, code, channel.ChannelName))
			hasElectricBank = true
			continue
		}
		for _, v := range channel.PayList {
			result = append(result, fmt.Sprintf("%s-%s|%s - %s", v["code"], v["type"], v["code"], v["name"]))
		}
	}
	return strings.Join(result, ",")
}

// PaymentGroups 支付分组
var PaymentGroups = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
}{
	ActionList: &ActionList{
		Model:    models.PaymentGroups,
		ViewFile: "payment_groups/list.html",
		Rows: func() interface{} {
			return &[]models.PaymentGroup{}
		},
		QueryCond: map[string]interface{}{
			"name": "%",
			"type": "=",
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.PaymentGroups,
		ViewFile: "payment_groups/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			return pongo2.Context{
				"channels": getPaymentChannels(platform),
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.PaymentGroups,
		ViewFile: "payment_groups/edit.html",
		Row: func() interface{} {
			return &models.PaymentGroup{}
		},
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			return pongo2.Context{
				"channels": getPaymentChannels(platform),
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.PaymentGroups,
	},
	ActionDelete: &ActionDelete{
		Model: models.PaymentGroups,
	},
}
