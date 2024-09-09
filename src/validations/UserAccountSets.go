package validations

import "sports-common/validation"

var UpAndDowns = func(platform string, data map[string]interface{}) error {
	return validation.New(data).
		Field("money").Uint("钱必须是大于0的整数").
		Field("bet_times").Int("投注倍数必须是正整数").
		Validate()
}
