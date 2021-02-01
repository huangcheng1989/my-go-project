package main

import (
	"awesomeProject/common"
	"awesomeProject/db"
	"awesomeProject/dbredis"
	"awesomeProject/entity"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"sort"
	"time"
)

func main() {
	//默认端口：8080
	r := gin.Default()
	r.GET("/download", HandelDownloadCutAndDriedList)

	r.Run()
}

func HandelDownloadCutAndDriedList(c *gin.Context) {
	dbredis.InitRedisCluster("prod")
	db.InitDB1("prod")

	listPopular1 := getListPopular1()
	listPopular2 := getListPopular2()
	listEmerging1 := getListEmerging1()
	listEmerging2 := getListEmerging2()

	xlsx := excelize.NewFile()
	xlsx.NewSheet("Sheet2")
	xlsx.SetSheetName("Sheet1", "人气创作者")
	xlsx.SetSheetName("Sheet2", "新晋创作者")
	xlsx.SetCellValue("人气创作者", "A1", "UserID")
	xlsx.SetCellValue("人气创作者", "B1", "VskitID")
	xlsx.SetCellValue("人气创作者", "C1", "Name")
	xlsx.SetCellValue("人气创作者", "D1", "目前排名")

	xlsx.SetCellValue("新晋创作者", "A1", "UserID")
	xlsx.SetCellValue("新晋创作者", "B1", "VskitID")
	xlsx.SetCellValue("新晋创作者", "C1", "Name")
	xlsx.SetCellValue("新晋创作者", "D1", "目前排名")

	for k, v := range listPopular1 {
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("A%d", k+2), v.UserId)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("B%d", k+2), v.VskitId)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("C%d", k+2), v.Name)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("D%d", k+2), v.Sort)
	}

	for k, v := range listPopular2 {
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("A%d", k+len(listPopular1)+3), v.UserId)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("B%d", k+len(listPopular1)+3), v.VskitId)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("C%d", k+len(listPopular1)+3), v.Name)
		xlsx.SetCellValue("人气创作者", fmt.Sprintf("D%d", k+len(listPopular1)+3), v.Sort)
	}

	for k, v := range listEmerging1 {
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("A%d", k+2), v.UserId)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("B%d", k+2), v.VskitId)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("C%d", k+2), v.Name)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("D%d", k+2), v.Sort)
	}

	for k, v := range listEmerging2 {
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("A%d", k+len(listEmerging1)+3), v.UserId)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("B%d", k+len(listEmerging1)+3), v.VskitId)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("C%d", k+len(listEmerging1)+3), v.Name)
		xlsx.SetCellValue("新晋创作者", fmt.Sprintf("D%d", k+len(listEmerging1)+3), v.Sort)
	}

	c.Header("Content-Type", "application/octet-stream")

	dateNow := time.Now().Format("20060102")
	filename := "CutAndDriedList_" + dateNow + ".xlsx"

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	_ = xlsx.Write(c.Writer)
}

func getListPopular1() []*entity.ActivityTalentRank {
	//执行查询语句
	query := "SELECT * FROM activity_talent_rank where user_id in ("
	for k, v := range common.Top10PopularCutAndDried {
		if k == len(common.Top10PopularCutAndDried)-1 {
			query += "'" + v + "'"
		} else {
			query += "'" + v + "',"
		}
	}
	query += ")"

	rows, err := db.DB1.Query(query)
	if err != nil {
		fmt.Println("查询出错了")
	}

	var listPopular1 []*entity.ActivityTalentRank
	//循环读取结果
	for rows.Next() {
		item := &entity.ActivityTalentRank{}
		//将每一行的结果都赋值到一个item对象中
		err := rows.Scan(&item.UserId, &item.VskitId, &item.Name, &item.Type, &item.VoteCount, &item.FollowerCount, &item.Pcc, &item.Phone, &item.CreateTime, &item.UpdateTime)
		if err != nil {
			fmt.Println("rows fail")
		}
		listPopular1 = append(listPopular1, item)
	}

	for _, v := range listPopular1 {
		result, _ := dbredis.GetRedisCluster().ZRevRank(context.Background(), "popular_talent_rank", v.UserId).Result()
		v.Sort = result + 1
	}

	// 重新排序，按sort从小到大
	sort.Sort(entity.TalentItemList(listPopular1))

	return listPopular1
}

func getListPopular2() []*entity.ActivityTalentRank {
	//执行查询语句
	query := "SELECT * FROM activity_talent_rank where user_id in ("
	for k, v := range common.Top11to50PopularCutAndDried {
		if k == len(common.Top11to50PopularCutAndDried)-1 {
			query += "'" + v + "'"
		} else {
			query += "'" + v + "',"
		}
	}
	query += ")"

	rows, err := db.DB1.Query(query)
	if err != nil {
		fmt.Println("查询出错了")
	}

	var listPopular2 []*entity.ActivityTalentRank
	//循环读取结果
	for rows.Next() {
		item := &entity.ActivityTalentRank{}
		//将每一行的结果都赋值到一个item对象中
		err := rows.Scan(&item.UserId, &item.VskitId, &item.Name, &item.Type, &item.VoteCount, &item.FollowerCount, &item.Pcc, &item.Phone, &item.CreateTime, &item.UpdateTime)
		if err != nil {
			fmt.Println("rows fail")
		}
		listPopular2 = append(listPopular2, item)
	}

	for _, v := range listPopular2 {
		result, _ := dbredis.GetRedisCluster().ZRevRank(context.Background(), "popular_talent_rank", v.UserId).Result()
		v.Sort = result + 1
	}

	// 重新排序，按sort从小到大
	sort.Sort(entity.TalentItemList(listPopular2))

	return listPopular2
}

func getListEmerging1() []*entity.ActivityTalentRank {
	//执行查询语句
	query := "SELECT * FROM activity_talent_rank where user_id in ("
	for k, v := range common.Top10EmergingCutAndDried {
		if k == len(common.Top10EmergingCutAndDried)-1 {
			query += "'" + v + "'"
		} else {
			query += "'" + v + "',"
		}
	}
	query += ")"

	rows, err := db.DB1.Query(query)
	if err != nil {
		fmt.Println("查询出错了")
	}

	var listEmerging1 []*entity.ActivityTalentRank
	//循环读取结果
	for rows.Next() {
		item := &entity.ActivityTalentRank{}
		//将每一行的结果都赋值到一个item对象中
		err := rows.Scan(&item.UserId, &item.VskitId, &item.Name, &item.Type, &item.VoteCount, &item.FollowerCount, &item.Pcc, &item.Phone, &item.CreateTime, &item.UpdateTime)
		if err != nil {
			fmt.Println("rows fail")
		}
		listEmerging1 = append(listEmerging1, item)
	}

	for _, v := range listEmerging1 {
		result, _ := dbredis.GetRedisCluster().ZRevRank(context.Background(), "emerging_talent_rank", v.UserId).Result()
		v.Sort = result + 1
	}

	// 重新排序，按sort从小到大
	sort.Sort(entity.TalentItemList(listEmerging1))

	return listEmerging1
}

func getListEmerging2() []*entity.ActivityTalentRank {
	//执行查询语句
	query := "SELECT * FROM activity_talent_rank where user_id in ("
	for k, v := range common.Top11to50EmergingCutAndDried {
		if k == len(common.Top11to50EmergingCutAndDried)-1 {
			query += "'" + v + "'"
		} else {
			query += "'" + v + "',"
		}
	}
	query += ")"

	rows, err := db.DB1.Query(query)
	if err != nil {
		fmt.Println("查询出错了")
	}

	var listEmerging2 []*entity.ActivityTalentRank
	//循环读取结果
	for rows.Next() {
		item := &entity.ActivityTalentRank{}
		//将每一行的结果都赋值到一个item对象中
		err := rows.Scan(&item.UserId, &item.VskitId, &item.Name, &item.Type, &item.VoteCount, &item.FollowerCount, &item.Pcc, &item.Phone, &item.CreateTime, &item.UpdateTime)
		if err != nil {
			fmt.Println("rows fail")
		}
		listEmerging2 = append(listEmerging2, item)
	}

	for _, v := range listEmerging2 {
		result, _ := dbredis.GetRedisCluster().ZRevRank(context.Background(), "emerging_talent_rank", v.UserId).Result()
		v.Sort = result + 1
	}

	// 重新排序，按sort从小到大
	sort.Sort(entity.TalentItemList(listEmerging2))

	return listEmerging2
}
