package users

import (
	"sports-admin/validations"
	common "sports-common"
	"sports-common/config"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func (ths *Users) Add(c *gin.Context) {
	postedData := request.GetPostedData(c)
	platform := request.GetPlatform(c)
	//需要做密码校验
	if err := validations.UserCreateSave(postedData); err != nil {
		response.Err(c, err.Error())
		return
	}
	req.SetTimeout(50 * time.Second)

	header := req.Header{
		"Accept": "application/json",
	}

	topId, _ := strconv.Atoi(postedData["top_id"].(string))
	dbSession := common.Mysql(platform)
	defer dbSession.Close()

	param := req.Param{
		"user_name":         postedData["username"].(string),
		"password":          tools.MD5(postedData["password"].(string)),
		"confirm_pwd":       tools.MD5(postedData["re_password"].(string)),
		"a_code":            1001,
		"device_id":         "pc",
		"register_id":       1,
		"register_url":      postedData["register_url"].(string),
		"platform":          platform,
		"withdraw_password": postedData["withdraw_password"].(string),
		"phone":             postedData["phone"].(string),
		"realname":          postedData["realname"].(string),
	}

	baseRegisterUrl := config.Get("admin_register.service") + config.Get("internal_api.register_url")
	registerUrl := baseRegisterUrl + "?platform=" + platform
	re, rerr := req.Post(registerUrl, header, param)
	if rerr != nil {
		response.Err(c, "系统异常")
		log.Err(rerr.Error())
		return
	}
	rres := Register{}
	_ = re.ToJSON(&rres)
	if rres.Errcode == 0 {
		_, ferr := dbSession.Table("users").Where("username=?", postedData["username"]).Update(map[string]interface{}{"top_id": topId, "top_name": postedData["top_name"]})
		if ferr != nil {
			response.Err(c, "程序错误")
			log.Err(ferr.Error())
			return
		}
		response.Ok(c)
		return
	} else {
		response.Err(c, rres.Message)
		return
	}
}
