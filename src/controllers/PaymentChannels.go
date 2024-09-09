package controllers

import (
	"encoding/json"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

type PaymentOut = models.PaymentOut
type AllPay = models.AllPay
type ChannelList = models.ChannelList

// PaymentChannels 支付通道
var PaymentChannels = struct {
	List   func(*gin.Context)
	Add    func(*gin.Context)
	Edit   func(*gin.Context)
	SaveDo func(*gin.Context)
	*ActionDelete
	*ActionState
}{
	List: func(c *gin.Context) { //默认首页
		cond := request.GetQueryCond(c, map[string]interface{}{
			"name":      "%",
			"code":      "%",
			"remark":    "%",
			"is_online": "=",
		})
		limit, offset := request.GetOffsets(c)
		Payments := make([]models.Payment, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("payments").Where(cond).OrderBy("id DESC").Limit(limit, offset).FindAndCount(&Payments)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		out := make([]PaymentOut, len(Payments))
		for i, v := range Payments {
			out[i].Id = v.Id
			out[i].Name = v.Name
			out[i].Code = v.Code
			out[i].CbIpList = v.CbIpList
			out[i].IsOnline = v.IsOnline
			out[i].Remark = v.Remark
			out[i].Weight = v.Weight
			out[i].Created = v.Created
			out[i].Updated = v.Updated
			str := []byte(v.ChannelList)
			stu := AllPay{}
			_ = json.Unmarshal(str, &stu)
			out[i].ChannelList = stu
		}
		viewData := pongo2.Context{
			"rows":  out,
			"total": total,
		}
		viewFile := "payment_channels/list.html"
		if request.IsAjax(c) {
			viewFile = "payment_channels/_list.html"
		}
		response.Render(c, viewFile, viewData)
	},
	Add: func(c *gin.Context) {
		response.Render(c, "payment_channels/edit.html", pongo2.Context{
			"method": "add",
			"vipIds": []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // vip等级相关数字
		})
	},
	Edit: func(c *gin.Context) {
		idStr := c.Query("id")
		id, _ := strconv.Atoi(idStr)
		payments := &models.Payment{}
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		if _, err := engine.Table("payments").Where("id=?", id).Get(payments); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		out := PaymentOut{}
		out.Id = payments.Id
		out.Name = payments.Name
		out.Code = payments.Code
		out.CbIpList = payments.CbIpList
		out.IsOnline = payments.IsOnline
		out.Remark = payments.Remark
		out.Weight = payments.Weight
		out.Created = payments.Created
		out.Updated = payments.Updated
		str := []byte(payments.ChannelList)
		stu := AllPay{}
		_ = json.Unmarshal(str, &stu)
		out.ChannelList = stu
		viewData := pongo2.Context{
			"r":      out,
			"method": "edit",
			//"vipIds": []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // vip等级相关数字
		}
		response.Render(c, "payment_channels/edit.html", viewData)
	},
	SaveDo: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		methods := postedData["methods"].(string)
		channelList := make(map[string]ChannelList)
		for _, s := range consts.PayChannelList {
			st := ChannelList{}
			st.Name = s["channel_name"].(string)
			// -- ClientPC - PC 相关配置
			if v, exists := postedData[s["channel_code"].(string)+"_client_pc"]; exists {
				st.ClientPc.Status = v.(string)
			} else if !exists {
				continue
			}
			st.ClientPc.FixedBankLimit = postedData[s["channel_code"].(string)+"_fixed_bank_pc"].(string)
			st.ClientPc.FixedMoneyLimit = postedData[s["channel_code"].(string)+"_fixed_money_pc"].(string)
			st.ClientPc.MinMoneyLimit = postedData[s["channel_code"].(string)+"_min_money_pc"].(string)
			st.ClientPc.MaxMoneyLimit = postedData[s["channel_code"].(string)+"_max_money_pc"].(string)
			st.ClientPc.Recommended = postedData[s["channel_code"].(string)+"_recommended_pc"].(string)
			st.ClientPc.Poll = postedData[s["channel_code"].(string)+"_poll_pc"].(string)
			st.ClientPc.Enable = postedData[s["channel_code"].(string)+"_enable_pc"].(string)
			st.ClientPc.Fare = postedData[s["channel_code"].(string)+"_fare_pc"].(string)
			st.ClientPc.DayMax = postedData[s["channel_code"].(string)+"_day_max_pc"].(string)
			// -- ClientAPP - APP 相关配置
			if v, exists := postedData[s["channel_code"].(string)+"_client_app"]; exists {
				st.ClientApp.Status = v.(string)
			} else if !exists {
				continue
			}
			st.ClientApp.FixedBankLimit = postedData[s["channel_code"].(string)+"_fixed_bank_app"].(string)
			st.ClientApp.FixedMoneyLimit = postedData[s["channel_code"].(string)+"_fixed_money_app"].(string)
			st.ClientApp.MinMoneyLimit = postedData[s["channel_code"].(string)+"_min_money_app"].(string)
			st.ClientApp.MaxMoneyLimit = postedData[s["channel_code"].(string)+"_max_money_app"].(string)
			st.ClientApp.Recommended = postedData[s["channel_code"].(string)+"_recommended_app"].(string)
			st.ClientApp.Poll = postedData[s["channel_code"].(string)+"_poll_app"].(string)
			st.ClientApp.Enable = postedData[s["channel_code"].(string)+"_enable_app"].(string)
			st.ClientApp.Fare = postedData[s["channel_code"].(string)+"_fare_app"].(string)
			st.ClientApp.DayMax = postedData[s["channel_code"].(string)+"_day_max_app"].(string)
			//vipListPc := make([]interface{}, 0)
			//vipListApp := make([]interface{}, 0)
			//for i := 0; i <= 10; i++ {
			//	if v, ok := postedData[s["channel_code"].(string)+"_vip_list_pc["+strconv.Itoa(i)+"]"]; ok {
			//		vipListPc = append(vipListPc, v.(string))
			//	}
			//	if v, ok := postedData[s["channel_code"].(string)+"_vip_list_app["+strconv.Itoa(i)+"]"]; ok {
			//		vipListApp = append(vipListApp, v.(string))
			//	}
			//}
			//st.ClientPc.VipList = vipListPc
			//st.ClientApp.VipList = vipListApp
			//if len(vipListPc) == 11 {
			//	st.ClientPc.Status = "开启"
			//} else {
			//	st.ClientPc.Status = "关闭"
			//}
			//if len(vipListApp) == 11 {
			//	st.ClientApp.Status = "开启"
			//} else {
			//	st.ClientApp.Status = "关闭"
			//}
			//st.ChannelFee = postedData[s["channel_code"].(string)+"_channel_fee"].(string)
			channelList[s["channel_code"].(string)] = st
		}
		s, _ := json.Marshal(channelList)
		postedData["channel_list"] = string(s)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		if methods == "add" { //新增
			imap := map[string]interface{}{
				"name":         postedData["name"],
				"code":         postedData["code"],
				"channel_list": postedData["channel_list"],
				"is_online":    postedData["is_online"],
				"cb_ip_list":   postedData["cb_ip_list"],
				"remark":       postedData["remark"],
				"weight":       postedData["weight"],
				"created":      tools.NowMicro(),
			}
			if _, err := engine.Table("payments").Insert(imap); err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, "新增失败")
				return
			}
		} else { //修改
			idStr, exists := postedData["id"].(string)
			if !exists || exists && idStr == "0" {
				response.Err(c, "id为空")
				return
			}
			id, _ := strconv.Atoi(idStr)
			umap := map[string]interface{}{
				"name":         postedData["name"],
				"code":         postedData["code"],
				"channel_list": postedData["channel_list"],
				"is_online":    postedData["is_online"],
				"cb_ip_list":   postedData["cb_ip_list"],
				"remark":       postedData["remark"],
				"weight":       postedData["weight"],
				"updated":      tools.NowMicro(),
			}
			if _, err := engine.Table("payments").Where("id=?", id).Update(umap); err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, "编辑失败")
				return
			}
		}
		caches.PaymentThirds.Load(platform)
		response.Ok(c)
	},
	ActionDelete: &ActionDelete{
		Model: models.Payments,
	},
	ActionState: &ActionState{
		Model: models.Payments,
		Field: "is_online",
	},
}
