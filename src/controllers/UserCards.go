package controllers

import (
	"fmt"
	"sports-admin/caches"
	"sports-admin/validations"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// UserCards 用户绑定银行卡列表
var UserCards = struct {
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
	Detail func(c *gin.Context)
	List   func(C *gin.Context)
}{
	ActionCreate: &ActionCreate{
		Model:    models.UserCards,
		ViewFile: "user_cards/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			return pongo2.Context{
				"banks": caches.Banks.All(platform),
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.UserCards,
		ViewFile: "user_cards/edit.html",
		Row: func() interface{} {
			return &models.UserCard{}
		},
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			return pongo2.Context{
				"banks": caches.Banks.All(platform),
			}
		},
	},
	ActionSave: &ActionSave{
		Model:     models.UserCards,
		Validator: validations.UserCards,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			name := (*m)["user_name"].(string)
			platform := request.GetPlatform(c)
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			var cards []models.UserCard
			if num, err := dbSession.Where("user_name=?", name).FindAndCount(&cards); err != nil {
				return err
			} else if num+1 > 3 {
				return errors.New("添加银行卡已达到上限")
			}
			var users []models.User
			if uErr := dbSession.Where("username = ?", name).Find(&users); uErr != nil {
				return uErr
			}
			(*m)["user_id"] = users[0].Id
			return nil
		},
		SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
			platform := request.GetPlatform(c)
			caches.Banks.Load(platform)
			name := (*m)["real_name"].(string)
			userName := (*m)["user_name"].(string)

			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			idsql := "select user_id from user_cards where user_name='" + userName + "'"
			res, ierr := dbSession.QueryString(idsql)
			if ierr != nil {
				log.Err(ierr.Error())
			}
			sql := "update users set realname =? where id=?"
			_, err := dbSession.Exec(sql, name, res[0]["user_id"])
			if err != nil {
				log.Err(err.Error())
			}
			csql := "update user_cards set real_name=? where user_id=?"
			_, cerr := dbSession.Exec(csql, name, res[0]["user_id"])
			if cerr != nil {
				log.Err(cerr.Error())
			}
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.UserCards,
		Row: func() interface{} {
			return &models.UserCard{}
		},
		DeleteAfter: func(c *gin.Context, i interface{}) {
			platform := request.GetPlatform(c)
			row := i.(*models.UserCard)
			_ = models.BlockedCards.Add(platform, row.CardNumber, row.UserName)
		},
	},
	Detail: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		name := c.Query("user_name")
		sql := "select real_name from user_cards where user_name='" + name + "' order by updated desc"
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
		}
		if len(res) > 0 {
			response.Result(c, res[0]["real_name"])
			return
		} else {
			response.Ok(c)
		}
	},
	List: func(c *gin.Context) {
		limit, offset := request.GetOffsets(c)
		sqlWhere := func() string {
			condWhere := " WHERE 1 = 1 "
			name := c.DefaultQuery("user_name", "")
			if name != "" {
				condWhere += " AND a.user_name LIKE '%" + name + "%' "
			}
			number := c.DefaultQuery("card_number", "")
			if number != "" {
				condWhere += " AND a.card_number LIKE '%" + number + "%' "
			}
			realName := c.DefaultQuery("real_name", "")
			if realName != "" {
				condWhere += " AND a.real_name LIKE '%" + realName + "%' "
			}
			return condWhere
		}()

		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		orderBy := " ORDER BY id DESC"
		sqlLimit := fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		fromAndJoin := " from user_cards a " +
			" left join banks b on a.bank_id = b.id " +
			" left join provinces c on a.province_id = c.id" +
			" left join cities d on a.city_id = d.id " +
			" left join districts f on a.district_id = f.id "
		sql := "select a.*,b.name as name ,c.name as p_name,d.name as c_name,f.name as d_name " + fromAndJoin + sqlWhere + orderBy + sqlLimit
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		sqlTotal := "SELECT COUNT(*) AS total " + fromAndJoin + sqlWhere
		total := 0
		rows, err := dbSession.QueryString(sqlTotal)
		if err != nil || len(rows) == 0 {
			total = 0
		} else {
			if totalInt, err := strconv.Atoi(rows[0]["total"]); err == nil {
				total = totalInt
			}
		}

		SetLoginAdmin(c)
		if request.IsAjax(c) {
			response.Render(c, "user_cards/_list.html", pongo2.Context{"rows": res, "total": total})
			return
		}
		//response.Render(c, "user_cards/list.html", pongo2.Context{"rows": res, "total": total})
		response.Render(c, "user_cards/index.html", pongo2.Context{"rows": res, "total": total})
	},
}
