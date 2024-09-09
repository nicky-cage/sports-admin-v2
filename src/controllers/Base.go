package controllers

import (
	"sports-admin/controllers/base_controller"
	"sports-admin/controllers/users"
	"sports-common/response"
)

type ViewData = response.ViewData

type IdRecord = base_controller.IdRecord

type ActionCreate = base_controller.ActionCreate

type ActionDelete = base_controller.ActionDelete

type ActionDetail = base_controller.ActionDetail

type ActionList = base_controller.ActionList

type ActionSave = base_controller.ActionSave

type ActionState = base_controller.ActionState

type ActionUpdate = base_controller.ActionUpdate

type ActionExport = base_controller.ActionExport
type ExportHeader = base_controller.ExportHeader

type Register = users.Register
