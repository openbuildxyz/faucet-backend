package utils

import "time"

func IsWithinLast24Hours(t time.Time) bool {
	// 获取当前时间
	now := time.Now()
	// 计算当前时间与给定时间的差值
	duration := now.Sub(t)

	// 如果时间差小于24小时，则该时间在24小时内
	return duration <= 24*time.Hour
}
