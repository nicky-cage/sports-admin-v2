package validations

import "sports-common/validation"

var ReceiveBankCard = func(platform string, data map[string]interface{}) error {
	return validation.New(data).
		Field("bank_card").BankCard("银行卡号格式不正确").
		Field("min_money_limit").NumericGt0("最小存款额必须是大于0的数").
		Field("max_money_limit").NumericGt0("最大存款额必须是大于0的数").
		Field("day_money_limit").NumericGt0("日限制金额必须是大于0的数").
		//Field("day_used_money").NumericEq0("当日已使用金额必须是大于等于0的数").
		Field("day_times_limit").Uint("日限制次数必须是大于0的整数").
		//Field("day_use_times").Uint0("当日已使用次数必须是整数").
		Validate()
}
