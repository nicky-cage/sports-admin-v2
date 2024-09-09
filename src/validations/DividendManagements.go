package validations

import (
	"sports-common/validation"
)

// CheckDividendMoney 红利检测
var CheckDividendMoney = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("money").Uint("钱必须是大于0的整数").
		Validate()
}

var CheckDividendFlowMultiple = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("flow_multiple").Uint0("流水倍数必须是正整数").
		Validate()
}
