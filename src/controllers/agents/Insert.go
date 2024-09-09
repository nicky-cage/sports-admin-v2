package agents

import (
	"fmt"
	"net/url"
	"sports-admin/controllers/users"
	"sports-admin/validations"
	common "sports-common"
	"sports-common/config"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func (ths *Agents) Insert(c *gin.Context) {
	var plan string
	var agent_type string
	data := request.GetPostedData(c)
	username := data["username"].(string)
	if data["agent_commission"] != nil {
		plan, _ = url.QueryUnescape(data["agent_commission"].(string))
	} else {
		plan = "默认"
	}
	if data["agent_type"] != nil {
		agent_type = data["agent_type"].(string)
	}
	times := time.Now().Unix()
	month := time.Now().Format("2006-01")

	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform, true)
	defer dbSession.Close()

	isBool := "select id from users where username='" + data["username"].(string) + "'"
	boolRes, _ := dbSession.QueryString(isBool)
	if len(boolRes) == 0 {
		//需要做密码校验
		if err := validations.UserCreateSave(data); err != nil {
			response.Err(c, err.Error())
			return
		}
		req.SetTimeout(50 * time.Second)
		header := req.Header{
			"Accept": "application/json",
		}
		param := req.Param{
			"user_name":    data["username"].(string),
			"password":     tools.MD5(data["password"].(string)),
			"confirm_pwd":  tools.MD5(data["re_password"].(string)),
			"a_code":       "1000",
			"platform":     platform,
			"device_id":    "system",
			"register_id":  1,
			"register_url": "https://" + c.Request.Host,
		}
		baseRegisterUrl := config.Get("admin_register.service") + config.Get("internal_api.register_url")
		registerUrl := baseRegisterUrl + "?platform=" + platform
		re, rerr := req.Post(registerUrl, header, param)
		if rerr != nil {
			response.Err(c, "系统异常")
			log.Err(rerr.Error())
			return
		}
		rres := users.Register{}
		_ = re.ToJSON(&rres)
		if rres.Errcode != 0 {
			response.Err(c, rres.Message)
			return
		}
	}

	updateSql := "update users set is_agent=%d,agent_type=%s ,transform_agent=%d,top_name='sys_test_agent',top_id='1000',agent_commission='%s',user_link='%s' where username='%s'"
	comSql := "insert into agent_commission_logs(user_id,type,month,username) values(?,?,?,?)"
	arr := strings.Split(username, ",")

	// checkSql := "select * from agent_commission_plans where agent_commission='" + data["agent_commission"].(string) + "'"
	checkSql := "select * from agent_commission_plans where agent_commission='" + plan + "'"
	checkRes, _ := dbSession.QueryString(checkSql)
	var temps models.ParameterGroup
	dbSession.Table("parameter_groups").Where("title='代理分享链接'").Get(&temps)
	for _, v := range arr {
		if v != "" {
			idSql := "select id,is_agent,top_id,top_name from users where username='%s'"
			idSqll := fmt.Sprintf(idSql, v)
			resId, err := dbSession.QueryString(idSqll)
			if err != nil {
				log.Err(err.Error())
				return
			}
			if len(resId) < 1 {
				response.Err(c, v+"无效会员")
				return
			}
			if resId[0]["is_agent"] == "1" {
				response.Err(c, v+"会员已是代理")
				return
			}

			// var usql string
			usql := fmt.Sprintf(updateSql, 1, agent_type, times, plan, temps.Name+"?a_code="+resId[0]["id"], v)
			_, uperr := dbSession.QueryString(usql)
			if uperr != nil {
				log.Err(uperr.Error())
				return
			}
			//创建代理佣金记录
			_, ierr := dbSession.Exec(comSql, resId[0]["id"], agent_type, month, v)
			if ierr != nil {
				log.Err(ierr.Error())
				return
			}
			name := "普通模式"
			if checkRes[0]["type"] == "2" {
				sql := "insert into agent_commission_plans(user_id,agent_commission,type,rate) values(?,?,?,?)"
				dbSession.Exec(sql, resId[0]["id"], data["agent_commission"], "2", checkRes[0]["rate"])
				name = "占城模式"
			}
			//添加分享链接。
			linkSql := "insert into promotions(user_id,username,agent_commission_type,name,link,plan_name,created) values(?,?,?,?,?,?,?)"
			dbSession.Exec(linkSql, resId[0]["id"], v, checkRes[0]["type"], "会员推广链接", temps.Name+"?a_code="+resId[0]["id"], name, time.Now().Unix())
			//报表也要加记录
			reportSql := "insert into user_daily_reports(user_id,username,day,created,game_code,top_id,top_name,is_agent,game_type) values(?,?,?,?,?,?,?,?,?)"
			_, err = dbSession.Exec(reportSql, resId[0]["id"], v, time.Now().Format("2006-01-02"), time.Now().Unix(), "0", resId[0]["top_id"], resId[0]["top_name"], 1, 0)
			if err != nil {
				log.Err(err.Error())
			}
		}
	}

	response.Result(c, "ok")
}
