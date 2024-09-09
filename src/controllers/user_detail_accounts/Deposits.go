package user_detail_accounts

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *UserDetailAccounts) Deposits(c *gin.Context) {
	cond := builder.NewCond()
	if status := c.DefaultQuery("status", ""); status != "" {
		if value, err := strconv.Atoi(status); err == nil && value >= 0 {
			cond = cond.And(builder.Eq{"status": value})
		}
	}
	if Type := c.DefaultQuery("channel_type", ""); Type != "" {
		if value, err := strconv.Atoi(Type); err == nil && value >= 0 {
			cond = cond.And(builder.Eq{"channel_type": value})
		}
	}
	if Types := c.DefaultQuery("type", ""); Types != "" {
		if value, err := strconv.Atoi(Types); err == nil && value >= 0 {
			cond = cond.And(builder.Eq{"type": value})
		}
	}
	if create := c.Query("created"); create != "" { //对时间进行处理
		areas := strings.Split(create, " - ")
		startAt, _ := time.Parse("2006-01-02 15:04:05", areas[0]+" 00:00:00")
		endAt, _ := time.Parse("2006-01-02 15:04:05", areas[1]+" 23:59:59")
		cond = cond.And(builder.Gte{"created": startAt.UnixMicro() - tools.SecondToMicro(3600*8)}).And(builder.Lte{"created": endAt.UnixMicro() - tools.SecondToMicro(3600*8)})
	}
	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	rows := []models.UserDeposit{}
	cond = cond.And(builder.Eq{"user_id": id})
	limit, offset := request.GetOffsets(c)
	total, err := dbSession.Table("user_deposits").Where(cond).OrderBy("created DESC").Limit(limit, offset).FindAndCount(&rows)
	if err != nil {
		log.Err(err.Error())
		return
	}

	// --------------------------------- 以下, 计算每次存款的下次存款时间 -----------------------------------
	type CreatedTime struct {
		ThisTime int64 `josn:"this_time"`
		NextTime int64 `josn:"next_time"`
	}
	sArr := []CreatedTime{}
	maxIndex := len(rows) - 1
	for i := maxIndex; i >= 0; i -= 1 {
		if i == 0 {
			sArr = append(sArr, CreatedTime{ThisTime: int64(rows[0].Created / 1000000), NextTime: tools.Now()})
			break
		}
		if rows[i].Status == 2 {
			thisTime := int64(rows[i].Created)
			sArr = append(sArr, CreatedTime{
				ThisTime: thisTime,
				NextTime: func() int64 {
					for j := i; j >= 0; j -= 1 {
						if rows[j].Status == 2 { // 从倒序里面查
							nextTime := int64(rows[j].Created)
							if nextTime > thisTime {
								return nextTime
							}
						}
					}
					return tools.Now()
				}(),
			})
		}
	}

	ss := new(UserDepositSumStruct)
	sumTotal, err := dbSession.Table("user_deposits").Where(cond).Sums(ss, "money", "arrive_money", "top_money", "discount")
	if err != nil {
		log.Err(err.Error())
	}

	page := c.Query("created")
	viewData := pongo2.Context{
		"rows":              rows,
		"total":             total,
		"id":                id,
		"deposits_money":    sumTotal[0],
		"deposits_arrive":   sumTotal[1],
		"deposits_top":      sumTotal[2],
		"deposits_discount": sumTotal[3],
		"currentTime":       tools.Now(), // 当前时间
		"createdTimes":      sArr,        // 所有存款时间之后的下笔存款时间
	}
	if page == "" {
		response.Render(c, "users/detail_deposits.html", viewData)
	} else {
		response.Render(c, "users/_detail_deposits.html", viewData)
	}
}
