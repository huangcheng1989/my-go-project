package main

import (
	"fmt"
	"time"
)

func main() {
	addVote := (float64(12) - 1) / float64(getHourDiffer(time.Now().Format("2006-01-02 15:04:05"), "2021-01-24 23:59:59"))
	date := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(getHourDiffer(date, "2021-01-24 23:59:59"), addVote)
}

//获取相差时间
func getHourDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}
