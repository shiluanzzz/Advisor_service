package model

import "fmt"

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

func GetStatusNameByCode(id int) string {
	if statusName[id] != "" {
		return statusName[id]
	} else {
		return fmt.Sprintf("状态%d不存在", id)
	}
}
func GetOrderEnableReplyId() []int {
	return []int{
		Rush, Pending,
	}
}
