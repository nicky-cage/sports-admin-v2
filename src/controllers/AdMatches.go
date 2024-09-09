package controllers

import (
	"sports-admin/controllers/ad_matches"
)

// AdMatches 轮播图
var AdMatches = ad_matches.AdMatches{
	ActionList:   ad_matches.ActionList,
	ActionCreate: ad_matches.ActionCreate,
	ActionUpdate: ad_matches.ActionUpdate,
	ActionSave:   ad_matches.ActionSave,
	ActionDelete: ad_matches.ActionDelete,
	ActionState:  ad_matches.ActionState,
}
