package filters

import (
	"sports-admin/caches"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
)

// ChannelList 渠道列表
type ChannelList = models.AllPay

// 判断rows的数量 - app
func colspanCountApp(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(ChannelList); isType {
		count := 0
		if value.Alipay.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Ebank.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Weixin.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Jdpay.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Qqpay.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Quickpay.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Unionpay.ClientApp.Status == "开启" {
			count += 1
		}
		return pongo2.AsValue(strconv.Itoa(count)), nil // 默认是 16 行
	}
	return pongo2.AsValue("1"), nil
}

// 判断rows的数量 - pc
func colspanCountPC(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(ChannelList); isType {
		count := 0
		if value.Alipay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Ebank.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Weixin.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Jdpay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Qqpay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Quickpay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Unionpay.ClientPc.Status == "开启" {
			count += 1
		}
		return pongo2.AsValue(strconv.Itoa(count)), nil // 默认是 16 行
	}
	return pongo2.AsValue("1"), nil
}

// 判断rows的数量 - pc
func rowspanCountChannel(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(ChannelList); isType {
		count := 2
		if value.Alipay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Ebank.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Weixin.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Jdpay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Qqpay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Quickpay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Unionpay.ClientPc.Status == "开启" {
			count += 1
		}
		if value.Alipay.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Ebank.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Weixin.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Jdpay.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Qqpay.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Quickpay.ClientApp.Status == "开启" {
			count += 1
		}
		if value.Unionpay.ClientApp.Status == "开启" {
			count += 1
		}
		return pongo2.AsValue(strconv.Itoa(count)), nil // 默认是 16 行
	}
	return pongo2.AsValue("16"), nil
}

// 会员标签
func onlinePaymentName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value := in.Interface().(string)
	arr := strings.Split(value, ":")
	if len(arr) != 2 {
		return pongo2.AsValue("-"), nil
	}

	platform := arr[0]
	name := caches.PaymentThirds.Get(platform, arr[1])
	return pongo2.AsValue(name), nil
}
