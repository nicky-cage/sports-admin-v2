package controllers

import (
	"sports-admin/controllers/ad_sponsors"
)

// AdSponsors 赞助配置
var AdSponsors = ad_sponsors.AdSponsors{
	ActionList:   ad_sponsors.ActionList,
	ActionCreate: ad_sponsors.ActionCreate,
	ActionUpdate: ad_sponsors.ActionUpdate,
	ActionSave:   ad_sponsors.ActionSave,
	ActionDelete: ad_sponsors.ActionDelete,
	ActionState:  ad_sponsors.ActionState,
}
