package validations

import (
	"errors"
	"sports-common/validation"
	models "sports-models"

	"xorm.io/builder"
)

// UserCards Validation说明
var UserCards = func(platform string, data map[string]interface{}) error {
	return validation.New(data).
		Field("user_name").Check("此用户信息不存在",
		// 检测用户名称是否存在
		func(oriUserName interface{}) error {
			userName := oriUserName.(string)
			cond := builder.NewCond().And(builder.Eq{"username": userName})
			user := models.User{}
			exists, err := models.Users.Find(platform, &user, cond)
			if err != nil {
				return err
			}
			if !exists {
				return errors.New("")
			}
			return nil
		}).
		Field("address").Length(2, 50, "缺少地址信息").
		Field("bank_id").Uint("缺少银行信息").
		Field("province_id").Uint("缺少省份信息").
		Field("city_id").Uint("缺少城市信息").
		Field("district_id").Uint("缺少县区信息").
		Field("branch_name").Length(2, 50, "缺少支行信息").
		Field("real_name").Length(2, 50, "缺少开户姓名").
		Field("card_number").Length(16, 19, "银行卡长度不对").
		Validate()
}
