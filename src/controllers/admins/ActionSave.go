package admins

import (
	"errors"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-admin/validations"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionSave = &base_controller.ActionSave{
	Validator: validations.Admins,
	Model:     models.Admins,
	CreateBefore: func(c *gin.Context, m *map[string]interface{}) error {
		name := (*m)["name"].(string) // 检测此用户名称是否存在
		cond := builder.NewCond().And(builder.Eq{"name": name})
		admin := models.Admin{}
		platform := request.GetPlatform(c)
		if exists, err := models.Admins.Find(platform, &admin, cond); exists {
			return errors.New("相同用户名称已经存在")
		} else if err != nil {
			return err
		}

		if password, exists := (*m)["password"]; !exists {
			return errors.New("必须输入登录密码")
		} else if rePassword, exists := (*m)["re_password"]; !exists {
			return errors.New("必须输入重复密码")
		} else if password.(string) != rePassword.(string) {
			return errors.New("两次输入的密码不相等 ")
		}

		return nil
	},
	SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
		platform := request.GetPlatform(c)
		roleIdStr := (*m)["role_id"].(string)
		roleId, _ := strconv.Atoi(roleIdStr)
		adminRoles := caches.AdminRoles.All(platform)

		// 如果有密码
		if v, exists := (*m)["password"]; exists {
			password := v.(string)
			if password != "" {
				rePasswordIn, exists := (*m)["re_password"]
				if !exists {
					return errors.New("必须输入重复密码")
				}
				rePassword := rePasswordIn.(string)
				if password != rePassword {
					return errors.New("两次输入的密码不相等")
				}

				salt := tools.Secret()                            //密盐
				realPassword := tools.GetPassword(password, salt) //真实密码
				(*m)["password"] = realPassword
				(*m)["salt"] = salt
			} else {
				delete(*m, "password")
				delete(*m, "re_password")
			}
		}

		if val, exists := adminRoles[uint32(roleId)]; !exists {
			return errors.New("相关角色信息查找错误")
		} else {
			(*m)["role_name"] = val.Name
		}

		// -- 检查google验证码
		if err := base_controller.CheckGoogleCode(platform, c, m); err != nil {
			return err
		}

		return nil
	},
}
