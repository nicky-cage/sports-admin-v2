package controllers

import (
	"sports-admin/controllers/ads"
)

// Ads app启动
var Ads = ads.Ads{
	ActionList:   ads.ActionList,
	ActionCreate: ads.ActionCreate,
	ActionUpdate: ads.ActionUpdate,
	ActionSave:   ads.ActionSave,
	ActionDelete: ads.ActionDelete,
	ActionState:  ads.ActionState,
}
