package service

import (
	qb "github.com/didi/gendry/builder"
	"github.com/fatih/structs"
	"log"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
)

var ADVISORTABLE = "advisor"

// NewAdvisor 新增一个顾问信息
func NewAdvisor(model *model.Advisor) int {
	// 转化数据并生成sql语句
	var data []map[string]interface{}
	data = append(data, structs.Map(model))
	cond, vals, err := qb.BuildInsert(ADVISORTABLE, data)
	if err != nil {
		log.Println("gendry SQL生成错误", err)
		return errmsg.ERROR_SQL_BUILD
	}

	// 执行sql语句
	_, err = utils.DbConn.Exec(cond, vals...)
	if err != nil {
		log.Println("数据新增错误", err)
		return errmsg.ERROR_SQL_BUILD
	}
	return errmsg.SUCCESS
}

// UpdateAdvisor 修改顾问的相关信息
func UpdateAdvisor(model *model.Advisor) int {
	where := map[string]interface{}{
		"phone": model.Phone,
	}
	// 把新的角色直接转化为map,去掉其中的value为空的key和其他相关值
	// phone,password,coin不可直接更新
	updates := structs.Map(model)
	delete(updates, "phone")
	delete(updates, "password")
	delete(updates, "coin")
	delete(updates, "total_order_num")
	delete(updates, "rank")
	delete(updates, "rank_num")
	cond, vals, err := qb.BuildUpdate(ADVISORTABLE, where, updates)
	if err != nil {
		log.Println("gendry SQL生成错误", err)
		return errmsg.ERROR_SQL_BUILD
	}
	_, err = utils.DbConn.Exec(cond, vals...)
	if err != nil {
		log.Println("数据库更新数据出错", err)
		return errmsg.ERROR_MYSQL
	}
	return errmsg.SUCCESS
}

func GetAdvisorInfo(phone string) (int, model.Advisor) {
	where := map[string]interface{}{
		"phone": phone,
	}
	selects := []string{
		"name", "phone", "coin", "total_order_num", "status",
		"rank", "rank_num", "work_experience", "bio", "about",
	}
	cond, values, err := qb.BuildSelect(ADVISORTABLE, where, selects)
	if err != nil {
		log.Println("gendry SQL生成错误", err)
		return errmsg.ERROR_SQL_BUILD, model.Advisor{}
	}
	row := utils.DbConn.QueryRow(cond, values...)
	res := model.Advisor{}
	// 能不能直接自动赋值到结构体对应的字段?
	err = row.Scan(
		&res.Name, &res.Phone, &res.Coin, &res.TotalOrderNum, &res.Status,
		&res.Rank, &res.RankNum, &res.WorkExperience, &res.Bio, &res.About,
	)
	if err != nil {
		log.Println("数据库查询用户信息出错", err)
		return errmsg.ERROR_MYSQL, model.Advisor{}
	}
	return errmsg.SUCCESS, res
}

func GetAdvisorList(page int) {
}
