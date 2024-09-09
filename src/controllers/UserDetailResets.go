package controllers

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var UserDetailResets = struct {
	Index func(c *gin.Context)
	*ActionSave
	SaveMore func(c *gin.Context)
}{
	Index: func(c *gin.Context) {
		//账户调整，建表？s
		var money string
		var username string
		var topName string
		var topId string
		var vip string
		id := c.Query("id")
		platform := request.GetPlatform(c)
		myClient := common.Mysql(platform)
		defer myClient.Close()
		sql := "select available,username from accounts where user_id = " + id
		res, err := myClient.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		sqll := "select vip,top_name,top_id from users where id=" + id
		vRes, err := myClient.QueryString(sqll)
		if err != nil {
			log.Err(err.Error())
			return
		}

		if len(res) > 0 {
			money = res[0]["available"]
			username = res[0]["username"]
		}
		if len(vRes) > 0 {
			vip = vRes[0]["vip"]
			topName = vRes[0]["top_name"]
			topId = vRes[0]["top_id"]
		}
		admin := GetLoginAdmin(c)
		var dataView = pongo2.Context{
			"id":       id,
			"money":    money,
			"admin":    admin.Name,
			"username": username,
			"vip":      vip,
			"top_name": topName,
			"top_id":   topId,
		}
		response.Render(c, "users/detail_resets.html", dataView)
	},
	ActionSave: &ActionSave{
		Model: models.UserResets,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			BillNo := tools.GetBillNo("T", 5)
			(*m)["bill_no"] = BillNo

			if (*m)["adjust_type"] == nil {
				return errors.New("必须选择调整类型")
			}
			if (*m)["adjust_method"] == nil {
				return errors.New("必须选择调整方式")
			}

			if (*m)["flow_limit"] == nil {
				return errors.New("必须选择流水限制")
			}
			if (*m)["remark"] == nil {
				return errors.New("必须填写备注")
			}
			flowMultiple, _ := strconv.ParseFloat((*m)["flow_multiple"].(string), 64)
			if (*m)["flow_limit"] == "2" {
				if flowMultiple <= 0 {
					return errors.New("流水倍数必须大于0")
				}
			} else {
				(*m)["flow_limit"] = "0"
			}

			money, _ := strconv.ParseFloat((*m)["adjust_money"].(string), 64)
			if money <= 0 {
				return errors.New("调整金额必须大于0")
			}
			return nil
		},
	},
	SaveMore: func(c *gin.Context) {
		data := request.GetPostedData(c)
		user_id := data["user_id"].(string)             // 用户编号
		username := data["username"].(string)           // 用户名称
		center_money := data["center_money"].(string)   // 中心钱包
		adjust_type := data["adjust_type"].(string)     // 调整类型
		adjust_method := data["adjust_method"].(string) // 调整方法
		flow_limit := data["flow_limit"].(string)       // 流水限制
		flow_multiple := data["flow_multiple"].(string) // 流水倍数
		adjust_money := data["adjust_money"].(string)   // 调整金额
		remark := data["remark"].(string)               // 备注
		admin_name := data["admin_name"].(string)       // 管理员名称
		bill_no := tools.GetBillNo("hl", 5)             // 订单号
		vip := data["vip"].(string)                     // VIP
		top_id := data["top_id"].(string)               // 上级编号
		//amoney, _ := strconv.ParseFloat(data["adjust_money"].(string), 64)
		if len(adjust_type) == 0 {
			response.Err(c, "必须选择调整类型")
			return
		}
		if len(adjust_method) == 0 {
			response.Err(c, "必须选择调整方式")
			return
		}

		if len(flow_limit) == 0 {
			response.Err(c, "必须选择流水限制")
			return
		}
		if len(remark) == 0 {
			response.Err(c, "必须填写备注")
			return
		}
		flowMultiple, _ := strconv.ParseFloat(flow_multiple, 64)
		if flow_limit == "2" {
			if flowMultiple <= 0 {
				response.Err(c, "流水倍数必须大于0")
				return
			}
		} else {
			flow_limit = "0"
		}

		money, _ := strconv.ParseFloat(adjust_money, 64)
		if money <= 0 {
			response.Err(c, "调整金额必须大于0")
			return
		}

		platform := request.GetPlatform(c)
		myClient := common.Mysql(platform)
		defer myClient.Close()
		if err := myClient.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = myClient.Rollback()
			response.Err(c, "事务启动失败")
			return
		}

		vSQL := " INSERT INTO user_resets " +
			"(user_id, username, center_money, adjust_type, adjust_method, flow_limit, flow_multiple, adjust_money, remark, updated, admin_name, bill_no, vip, top_id, created) " +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) "
		//myClient := common.Mysql(platform)
		currentTime := tools.NowMicro()
		_, err := myClient.Exec(vSQL, user_id, username, center_money, adjust_type, adjust_method, flow_limit, flow_multiple, adjust_money, remark, currentTime, admin_name, bill_no, vip, top_id, currentTime)
		if err != nil {
			log.Err(err.Error())
			return
		}
		// type_id := 57       //增加
		// multiple_type := 22 //父类手动调整
		// balance, _ := strconv.ParseFloat(center_money, 64)
		// if adjust_method == "2" {
		// 	amoney = -amoney
		// 	balance = balance + amoney
		// 	type_id = 58 //减少
		// 	remark = "(调整减少):" + remark
		// } else {
		// 	balance = balance + amoney
		// 	remark = "(调整增加):" + remark
		// }
		// mSQL := " update accounts set balance=balance+ ? ,available=available+?, updated= ?  where  user_id =? "
		// _, merr := myClient.Exec(mSQL, amoney, amoney, time.Now().Unix(), user_id)
		// if merr != nil {
		// 	log.Logger.Error(merr.Error())
		// 	_ = myClient.Rollback()
		// 	log.Err(merr.Error())
		// 	return
		// }
		// //账变数据
		// if adjust_type == "1" {
		// 	remark = "红利补发" + remark
		// } else if adjust_type == "2" {
		// 	remark = "系统调整" + remark
		// } else if adjust_type == "3" {
		// 	remark = "输赢调整" + remark
		// }

		// // 父类  22   multiple_type   multiple_name 调整    子类 57 调整增加  58	调整减少 type_id
		// IP := c.ClientIP()
		// administrator := GetLoginAdmin(c)
		// nSQL := "INSERT INTO transactions (serial_number,user_id,username,type_id,multiple_type,  multiple_name,description,amount,previous_balance,previous_available,balance,available,administrator,ip,safekey,previous_frozen,frozen,created,updated) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?, ?, ?)  "
		// _, nerr := myClient.Exec(nSQL, bill_no, user_id, username, type_id, multiple_type, "调整", remark, amoney, center_money, center_money, balance, balance, administrator.Name, IP, "", 0, 0, time.Now().Unix(), time.Now().Unix())
		// if nerr != nil {
		// 	log.Logger.Error(nerr.Error())
		// 	_ = myClient.Rollback()
		// 	log.Err(nerr.Error())
		// 	return
		// }
		_ = myClient.Commit()
		response.Ok(c)
	},
}
