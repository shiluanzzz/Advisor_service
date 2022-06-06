package model

import "fmt"

// 订单状态
const (
	Pending = iota
	Rush
	Expired
	Completed
)

var statusName = map[int]string{
	Pending:   "Pending",
	Rush:      "Rush",
	Expired:   "Expired",
	Completed: "Completed",
}

func GetOrderStatusNameById(id int) string {
	if res := statusName[id]; res != "" {
		return res
	} else {
		return fmt.Sprintf("状态%d不存在", id)
	}
}
func GetOrderEnableReplyId() []int {
	return []int{
		Rush, Pending,
	}
}

// 用户性别枚举
const (
	Unknown = iota
	Male
	Female
)

var genderName = map[int]string{
	Male:    "Male",
	Female:  "Female",
	Unknown: "Not Specified",
}

func GetGenderNameById(id int) string {
	if res := genderName[id]; res != "" {
		return res
	} else {
		return fmt.Sprintf("性别%d不存在", id)
	}
}

// 订单回复状态枚举
const (
	NotComment = iota
	Commented
)

// 顾问服务状态枚举
const (
	AdvisorServiceNotOpen = iota
	AdvisorServiceOpen
)
