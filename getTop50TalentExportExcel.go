package main

import (
	"awesomeProject/db"
	"awesomeProject/dbredis"
	"awesomeProject/entity"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	//默认端口：8080
	r := gin.Default()
	r.GET("/download", GetTop50TalentList)

	r.Run()
}

func GetTop50TalentList(c *gin.Context) {
	dbredis.InitRedisCluster("prod")
	db.InitDB2()

	listPopular := getListPopular()
	listEmerging := getListEmerging()

	xlsx := excelize.NewFile()
	xlsx.NewSheet("Sheet2")
	xlsx.SetSheetName("Sheet1", "人气创作者")
	xlsx.SetSheetName("Sheet2", "新晋创作者")
	xlsx.SetCellValue("人气创作者", "A1", "UserID")
	xlsx.SetCellValue("人气创作者", "B1", "VskitID")
	xlsx.SetCellValue("人气创作者", "C1", "Name")
	xlsx.SetCellValue("人气创作者", "D1", "目前排名")
	xlsx.SetCellValue("人气创作者", "E1", "投票数")

	xlsx.SetCellValue("新晋创作者", "A1", "UserID")
	xlsx.SetCellValue("新晋创作者", "B1", "VskitID")
	xlsx.SetCellValue("新晋创作者", "C1", "Name")
	xlsx.SetCellValue("新晋创作者", "D1", "目前排名")
	xlsx.SetCellValue("新晋创作者", "E1", "投票数")

	for k, v := range listPopular {
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("A%d", k+2), v.UserId)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("B%d", k+2), v.VskitId)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("C%d", k+2), v.UserName)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("D%d", k+2), v.Sort)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("E%d", k+2), v.VoteNum)
	}

	for k, v := range listEmerging {
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("A%d", k+2), v.UserId)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("B%d", k+2), v.VskitId)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("C%d", k+2), v.UserName)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("D%d", k+2), v.Sort)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("E%d", k+2), v.VoteNum)
	}

	c.Header("Content-Type", "application/octet-stream")

	dateNow := time.Now().Format("20060102")
	filename := "Top50TalentList_" + dateNow + ".xlsx"

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	_ = xlsx.Write(c.Writer)
}

func getListPopular() []*entity.TalentItem {
	var listPopular []*entity.TalentItem
	result, err := dbredis.GetRedisCluster().ZRevRangeWithScores(context.Background(), "popular_talent_rank", 0, 49).Result()
	if err != nil {
		return nil
	}

	var userIds []string
	for i := 0; i < len(result); i++ {
		item := &entity.TalentItem{
			UserId:  result[i].Member.(string),
			Sort:    i + 1,
			VoteNum: int(result[i].Score),
		}
		listPopular = append(listPopular, item)
		userIds = append(userIds, item.UserId)
	}

	userInfos := getUserInfo(userIds)
	for _, v := range listPopular {
		for _, vv := range userInfos {
			if v.UserId == vv.UserId {
				v.VskitId = vv.VskitId
				v.UserName = vv.Name
				break
			}
		}
	}

	return listPopular
}

func getListEmerging() []*entity.TalentItem {
	var listEmerging []*entity.TalentItem
	result, err := dbredis.GetRedisCluster().ZRevRangeWithScores(context.Background(), "emerging_talent_rank", 0, 49).Result()
	if err != nil {
		return nil
	}

	var userIds []string
	for i := 0; i < len(result); i++ {
		item := &entity.TalentItem{
			UserId:  result[i].Member.(string),
			Sort:    i + 1,
			VoteNum: int(result[i].Score),
		}
		listEmerging = append(listEmerging, item)
		userIds = append(userIds, item.UserId)
	}

	userInfos := getUserInfo(userIds)
	for _, v := range listEmerging {
		for _, vv := range userInfos {
			if v.UserId == vv.UserId {
				v.VskitId = vv.VskitId
				v.UserName = vv.Name
				break
			}
		}
	}

	return listEmerging
}

func getUserInfo(userIds []string) []*entity.UserDetail {
	//执行查询语句
	query := "SELECT user_id,vskit_id,name FROM user_detail where user_id in ("
	for k, v := range userIds {
		if k == len(userIds)-1 {
			query += "'" + v + "'"
		} else {
			query += "'" + v + "',"
		}
	}
	query += ")"

	rows, err := db.DB2.Query(query)
	if err != nil {
		fmt.Println("查询出错了")
	}

	var userInfos []*entity.UserDetail
	//循环读取结果
	for rows.Next() {
		item := &entity.UserDetail{}
		//将每一行的结果都赋值到一个item对象中
		err := rows.Scan(&item.UserId, &item.VskitId, &item.Name)
		if err != nil {
			fmt.Println("rows fail")
		}
		userInfos = append(userInfos, item)
	}

	return userInfos
}
