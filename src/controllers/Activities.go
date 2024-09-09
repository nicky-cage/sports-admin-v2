package controllers

import (
	"sports-admin/controllers/activities"
	"sports-admin/controllers/activity_applies"
	"sports-admin/controllers/activity_audits"
	"sports-admin/controllers/activity_histories"
	"sports-admin/controllers/activity_managements"
	models "sports-models"
)

// Activities 常规活动
var Activities = activities.Activities{
	ActionDelete: &ActionDelete{Model: models.Activities},
	ActionState:  &ActionState{Model: models.Activities},
}

// ActivitiesAudits 活动礼金-审核列表
var ActivitiesAudits = activity_audits.ActivityAudits{}

// ActivitiesHrs 活动礼金-历史列表
var ActivitiesHrs = activity_histories.ActivityHrs{}

// ActivitiesManagements 活动礼金申请
var ActivitiesManagements = activity_managements.ActivityManagements{}

// ActivityApplies 活动申请
var ActivityApplies = activity_applies.ActivityApplies{}
