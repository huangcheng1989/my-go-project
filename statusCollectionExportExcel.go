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
	r.GET("/download", HanderDownload1)

	r.Run()
}

func HanderDownload1(c *gin.Context) {
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "StatusID")
	xlsx.SetCellValue("Sheet1", "B1", "Type（类型）")
	xlsx.SetCellValue("Sheet1", "C1", "Source（来源）")
	xlsx.SetCellValue("Sheet1", "D1", "Link（链接）")
	xlsx.SetCellValue("Sheet1", "E1", "Reason（原因）")
	xlsx.SetCellValue("Sheet1", "F1", "Phone number（手机号码）")
	xlsx.SetCellValue("Sheet1", "G1", "Country（国家）")
	xlsx.SetCellValue("Sheet1", "H1", "提交时间")
	xlsx.SetCellValue("Sheet1", "I1", "Selected（选择）")

	list1 := selectAll1()

	var newList []*entity.StatusRes
	for _, v := range list1 {
		item := &entity.StatusRes{
			Id:         v.Id,
			Type:       getType(v.Type, v.TypeExtra),
			Source:     getSource(v.Source, v.SourceExtra),
			Link:       v.FileLink,
			Reason:     v.UploadReason,
			Phone:      v.Pcc + " " + v.Phone,
			Country:    v.Country,
			CreateTime: v.CreateTime,
			Selected:   v.Selected,
		}
		newList = append(newList, item)
	}

	for k, v := range newList {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", k+2), v.Id)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", k+2), v.Type)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", k+2), v.Source)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", k+2), v.Link)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", k+2), v.Reason)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", k+2), v.Phone)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("G%d", k+2), v.Country)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d", k+2), time.Unix(int64(v.CreateTime), 0).Format("2006-01-02 15:04:05"))
		//xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d", k+2), v.Selected)
	}

	c.Header("Content-Type", "application/octet-stream")

	dateNow := time.Now().Format("20060102")
	filename := "Status-Collection_" + dateNow + ".xlsx"

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	_ = xlsx.Write(c.Writer)
}

func selectAll1() []*entity.ActivityStatusCollection {
	db.InitDB1("prod")

	var count int
	_ = db.DB1.QueryRow("SELECT count(1) FROM activity_status_collection").Scan(&count)

	var list1 []*entity.ActivityStatusCollection
	for i := 0; i < count; i += 200 {
		//执行查询语句
		rows, err := db.DB1.Query(fmt.Sprintf("SELECT * FROM activity_status_collection order by create_time desc limit %d,200", i))
		if err != nil {
			fmt.Println("查询出错了")
		}

		var list []*entity.ActivityStatusCollection
		//循环读取结果
		for rows.Next() {
			item := &entity.ActivityStatusCollection{}
			//将每一行的结果都赋值到一个user对象中
			err := rows.Scan(&item.Id, &item.UserId, &item.DeviceId, &item.Type, &item.TypeExtra, &item.Source, &item.SourceExtra,
				&item.FileLink, &item.UploadReason, &item.Pcc, &item.Phone, &item.Country, &item.Selected, &item.CreateTime, &item.UpdateTime)
			if err != nil {
				fmt.Println("rows fail")
			}
			//将user追加到users的这个数组中
			list = append(list, item)
		}

		list1 = append(list1, list...)
	}
	return list1
}

func getType(ttype int, typeExtra string) string {
	switch ttype {
	case 1:
		return "Meme"
	case 2:
		return "News"
	case 3:
		return "Record life"
	case 4:
		return "Blessings"
	}
	return "Others: " + typeExtra
}

func getSource(source int, sourceExtra string) string {
	switch source {
	case 1:
		return "Whatsapp"
	case 2:
		return "Instagram"
	case 3:
		return "Snapchat"
	}
	return "Others: " + sourceExtra
}
