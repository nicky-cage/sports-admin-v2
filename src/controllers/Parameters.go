package controllers

import (
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// Parameters 参数
var Parameters = struct {
	//*ActionList
	*ActionUpdate
	*ActionSave
	*ActionCreate
	List func(c *gin.Context)
}{
	//ActionList: &ActionList{
	//	Model:    models.Parameters,
	//	ViewFile: "parameters/list.html",
	//	Rows: func() interface{} {
	//		return &[]models.Parameter{}
	//	},
	//	GetQueryCond: func(c *gin.Context) builder.Cond {
	//		cond := builder.NewCond()
	//		if idStr, exists := c.GetQuery("id"); exists {
	//			if id, err := strconv.Atoi(idStr); err == nil {
	//				cond = cond.And(builder.Eq{"group_id": id})
	//			}
	//		}
	//
	//		return cond
	//	},
	//	QueryCond: map[string]interface{}{
	//		"name": "%",
	//	},
	//	OrderBy: func(c *gin.Context) string {
	//		return "sort DESC"
	//	},
	//},
	List: func(c *gin.Context) {
		id := c.Query("id")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select * from parameters where group_id= " + id
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		response.Render(c, "parameters/list.html", pongo2.Context{"rows": res})
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.Parameters,
		ViewFile: "parameters/edit.html",
		Row: func() interface{} {
			return &models.Parameter{}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Parameters,
		UpdateAfter: func(c *gin.Context, m *map[string]interface{}) {
			key := (*m)["name"]
			cacheKey := consts.ParamCacheKey + key.(string)
			platform := request.GetPlatform(c)
			redis := common.Redis(platform)
			defer common.RedisRestore(platform, redis)
			redis.Del(cacheKey)
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.Parameters,
		ViewFile: "parameters/add.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			sql := "select id,title from parameter_groups "
			res, err := dbSession.QueryString(sql)
			if err != nil {
				log.Err(err.Error())
				return nil
			}
			return pongo2.Context{"rows": res}
		},
	},
}
