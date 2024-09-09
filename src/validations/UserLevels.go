package validations

import (
	"sports-common/validation"
)

// UserLevels Validation说明
var UserLevels = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("name").Length(2, 10, "会员等级名称不能为空").
		Field("upgrade_deposit").Numeric("升级存款必须应是数字").
		Field("hold_stream").Numeric("保级流水要求应是数字").
		Field("upgrade_stream").Numeric("升级注水要求应是数字").
		Field("upgrade_bonus").Numeric("升级红利应是数字").
		Field("birth_bonus").Numeric("生日红利应是数字").
		Field("month_bonus").Numeric("月度红利应是数字").
		Field("day_withdraw_count").Uint("每日提款次数应是数字").
		Field("day_withdraw_total").Numeric("每日提款限额应是数字").
		Validate()
}
