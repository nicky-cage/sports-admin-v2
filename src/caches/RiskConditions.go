package caches

import (
	models "sports-models"
)

const keyRiskConditions = "risk_conditions"

// RiskConditions 风控条件
var RiskConditions = struct {
	Load func(string)
	All  func(string) models.RiskCondition
}{
	Load: func(platform string) {
		riskConditions := map[uint32]models.RiskCondition{}
		rs := models.RiskCondition{}

		_, err := models.RiskConditions.FindById(platform, 1, &rs)
		if err != nil {
			return
		}

		riskConditions[1] = rs

		_ = setCache(platform, keyRiskConditions, riskConditions)

	},
	All: func(platform string) models.RiskCondition {
		riskConditions := map[uint32]models.RiskCondition{}
		_ = getCache(platform, keyRiskConditions, &riskConditions)

		return riskConditions[1]
	},
}
