package controllers

import (
	"sports-admin/controllers/blacklists"
	"sports-admin/controllers/blocked_cards"
	"sports-admin/controllers/blocked_devices"
	"sports-admin/controllers/blocked_ips"
	"sports-admin/controllers/blocked_mails"
	"sports-admin/controllers/blocked_phones"
)

// BlackLists 黑名单
var BlackLists = blacklists.BlackLists{}

// BlockedCards 黑名单 - 银行卡号
var BlockedCards = blocked_cards.BlockedCards{
	ActionList:   blocked_cards.ActionList,
	ActionCreate: blocked_cards.ActionCreate,
	ActionUpdate: blocked_cards.ActionUpdate,
	ActionSave:   blocked_cards.ActionSave,
	ActionDelete: blocked_cards.ActionDelete,
}

// BlockedDevices 黑名单 - 设备
var BlockedDevices = blocked_devices.BlockedDevices{
	ActionList:   blocked_devices.ActionList,
	ActionCreate: blocked_devices.ActionCreate,
	ActionUpdate: blocked_devices.ActionUpdate,
	ActionSave:   blocked_devices.ActionSave,
	ActionDelete: blocked_devices.ActionDelete,
}

// BlockedIps 黑名单 - 设备
var BlockedIps = blocked_ips.BlockedIps{
	ActionList:   blocked_ips.ActionList,
	ActionCreate: blocked_ips.ActionCreate,
	ActionUpdate: blocked_ips.ActionUpdate,
	ActionSave:   blocked_ips.ActionSave,
	ActionDelete: blocked_ips.ActionDelete,
}

// BlockedMails 黑名单 - 电子邮件
var BlockedMails = blocked_mails.BlockedMails{
	ActionList:   blocked_mails.ActionList,
	ActionCreate: blocked_mails.ActionCreate,
	ActionUpdate: blocked_mails.ActionUpdate,
	ActionSave:   blocked_mails.ActionSave,
	ActionDelete: blocked_mails.ActionDelete,
}

// BlockedPhones 黑名单 - 手机号码
var BlockedPhones = blocked_phones.BlockedPhones{
	ActionList:   blocked_phones.ActionList,
	ActionCreate: blocked_phones.ActionCreate,
	ActionUpdate: blocked_phones.ActionUpdate,
	ActionSave:   blocked_phones.ActionSave,
	ActionDelete: blocked_phones.ActionDelete,
}
