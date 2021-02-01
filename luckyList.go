package main

import (
	"awesomeProject/db"
	"awesomeProject/entity"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	//默认端口：8080
	r := gin.Default()
	r.GET("/download", HanderDownload)

	r.Run()
}

func HanderDownload(c *gin.Context) {
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "UserID")
	xlsx.SetCellValue("Sheet1", "B1", "VskitID")
	xlsx.SetCellValue("Sheet1", "C1", "Name")
	xlsx.SetCellValue("Sheet1", "D1", "手机号码")
	xlsx.SetCellValue("Sheet1", "E1", "奖品")
	xlsx.SetCellValue("Sheet1", "F1", "国家")
	xlsx.SetCellValue("Sheet1", "G1", "中奖时间(UTC)")

	luckyList := selectAll()
	var userIdList []string
	for _, v := range luckyList {
		userIdList = append(userIdList, v.UserId)
	}

	var newList []*entity.LuckRes
	userList := selectUserList(userIdList)
	for _, v := range luckyList {
		for _, vv := range userList {
			if v.UserId == vv.UserId {
				item := &entity.LuckRes{
					UserId:     v.UserId,
					VskitId:    vv.VskitId,
					Name:       vv.Name,
					PrizeName:  getPrizeName(v.PrizeId),
					Country:    v.Country,
					CreateTime: time.Unix(v.CreatedTime, 0).UTC().Format("2006/01/02 15:04:05"),
				}
				if v.Pcc != "" && v.Phone != "" {
					item.Phone = v.Pcc + "-" + v.Phone
				}
				newList = append(newList, item)
				break
			}
		}
	}

	for k, v := range newList {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", k+2), v.UserId)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", k+2), v.VskitId)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", k+2), v.Name)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", k+2), v.Phone)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", k+2), v.PrizeName)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", k+2), v.Country)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("G%d", k+2), v.CreateTime)
	}

	c.Header("Content-Type", "application/octet-stream")

	dateNow := time.Now().Format("20060102")
	filename := "LuckList_" + dateNow + ".xlsx"

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	_ = xlsx.Write(c.Writer)
}

func selectAll() []*entity.ActivityTalentRankLotteryRecord {
	db.InitDB1("prod")

	//执行查询语句
	rows, err := db.DB1.Query("SELECT * FROM activity_talent_rank_lottery_record WHERE prize_id IN (1,2,3,4,5) order by created_time desc")
	if err != nil {
		fmt.Println("查询出错了")
	}
	var luckyList []*entity.ActivityTalentRankLotteryRecord
	//循环读取结果
	for rows.Next() {
		item := &entity.ActivityTalentRankLotteryRecord{}
		//将每一行的结果都赋值到一个user对象中
		err := rows.Scan(&item.Id, &item.DeviceId, &item.UserId, &item.PrizeId, &item.Country, &item.Pcc, &item.Phone, &item.CreatedTime)
		if err != nil {
			fmt.Println("rows fail")
		}
		//将user追加到users的这个数组中
		luckyList = append(luckyList, item)
	}
	return luckyList
}

func selectUserList(userIdList []string) []*entity.UserDetail {
	db.InitDB2()

	//执行查询语句
	query := "SELECT user_id, vskit_Id, name FROM user_detail WHERE 1 = 1"
	if len(userIdList) > 0 {
		query += " and user_id in ("
		for k, v := range userIdList {
			if k+1 == len(userIdList) {
				query += "'" + v + "')"
			} else {
				query += "'" + v + "',"
			}
		}
	}

	rows, err := db.DB2.Query(query)
	if err != nil {
		fmt.Println("查询出错了")
	}
	var userList []*entity.UserDetail
	//循环读取结果
	for rows.Next() {
		item := &entity.UserDetail{}
		//将每一行的结果都赋值到一个user对象中
		err := rows.Scan(&item.UserId, &item.VskitId, &item.Name)
		if err != nil {
			fmt.Println("rows fail")
		}
		//将user追加到users的这个数组中
		userList = append(userList, item)
	}
	return userList
}

func getPrizeName(prizeId int) string {
	switch prizeId {
	case 1:
		return "1.5g flow"
	case 2:
		return "3g flow"
	case 3:
		return "4.5g flow"
	case 4:
		return "10g flow"
	case 5:
		return "a Infinix Hot 9"
	}
	return ""
}
