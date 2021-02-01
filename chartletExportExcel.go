package main

import (
	"awesomeProject/db"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func main() {
	//默认端口：8080
	r := gin.Default()
	r.GET("/download", HandelDownloadChartlet)

	r.Run()
}

func HandelDownloadChartlet(c *gin.Context) {
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "auto_id")
	xlsx.SetCellValue("Sheet1", "B1", "贴纸id")
	xlsx.SetCellValue("Sheet1", "C1", "贴纸英文名称")
	xlsx.SetCellValue("Sheet1", "D1", "视频id")
	xlsx.SetCellValue("Sheet1", "E1", "视频标题")
	xlsx.SetCellValue("Sheet1", "F1", "视频url")
	xlsx.SetCellValue("Sheet1", "G1", "作者userId")
	xlsx.SetCellValue("Sheet1", "H1", "作者vskitId")
	xlsx.SetCellValue("Sheet1", "I1", "真实播放总量")
	xlsx.SetCellValue("Sheet1", "J1", "真实点赞总量")
	xlsx.SetCellValue("Sheet1", "K1", "真实评论总量")
	xlsx.SetCellValue("Sheet1", "L1", "真实分享总量")
	xlsx.SetCellValue("Sheet1", "M1", "视频一级标签名称")
	xlsx.SetCellValue("Sheet1", "N1", "完播率")

	list := getList()
	for k, v := range list {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", k+2), getAutoIdByUuid(v.Id))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", k+2), v.Id)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", k+2), v.ChartletTitle)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", k+2), v.VideoId)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", k+2), v.VideoTitle)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", k+2), v.VideoUrl)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("G%d", k+2), v.UserId)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d", k+2), v.VskitId)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("I%d", k+2), v.Views-v.ManualView)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("J%d", k+2), v.likes-v.ManualLike)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("K%d", k+2), v.CommentCount-v.ManualCommentCount)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("L%d", k+2), v.Shares-v.ManualShare)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("M%d", k+2), v.TagTitle)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("N%d", k+2), v.RatePc)
	}

	c.Header("Content-Type", "application/octet-stream")

	dateNow := time.Now().Format("20060102")
	filename := "ChartletList_" + dateNow + ".xlsx"

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")

	//回写到web 流媒体 形成下载
	_ = xlsx.Write(c.Writer)
}

func getList() []*chartlet {
	db.InitDB2()

	chartletList := getChartletList()
	var chartletIds []string
	for _, v := range chartletList {
		chartletIds = append(chartletIds, v.Id)
	}

	videoList := getVideoList(chartletIds)
	var videoIds []string
	for _, v := range videoList {
		videoIds = append(videoIds, v.VideoId)
		for _, vv := range chartletList {
			if v.Id == vv.Id {
				v.ChartletTitle = vv.ChartletTitle
				break
			}
		}
	}

	newVideoIds := removeRepeatedElement1(videoIds)
	tagList := getTagList(newVideoIds)

	for _, v := range videoList {
		for _, vv := range tagList {
			if v.VideoId == vv.VideoId {
				v.TagTitle = vv.TagTitle
				break
			}
		}
	}

	var rateList []*videoRatePc
	for i := 0; i < len(newVideoIds); i += 100 {
		newVideoIds1 := newVideoIds[i : i+100]
		res := getRatePc(newVideoIds1)
		if len(res) > 0 {
			rateList = append(rateList, res...)
		}
	}

	for _, v := range videoList {
		for _, vv := range rateList {
			if v.VideoId == vv.VideoId {
				v.RatePc = vv.RatePc
				break
			}
		}
	}
	return videoList
}

type tag struct {
	VideoId  string `json:"video_id"`
	TagTitle string `json:"tag_title"`
}

func getTagList(newVideoIds []string) []*tag {
	//执行查询语句
	query := "SELECT video_id,title as tag_title FROM video_tags_info where video_id in ("
	for k, v := range newVideoIds {
		if k == len(newVideoIds)-1 {
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

	var list []*tag
	//循环读取结果
	for rows.Next() {
		item := &tag{}
		//将每一行的结果都赋值到一个item对象中
		err := rows.Scan(&item.VideoId, &item.TagTitle)
		if err != nil {
			fmt.Println("rows fail")
		}
		list = append(list, item)
	}
	return list
}

func getVideoList(chartletIds []string) []*chartlet {
	//执行查询语句
	query := "SELECT b.vskit_id,a.chartlet_id as id,a.video_id,a.title as video_title,a.video_url,a.user_id,a.views,a.manual_view,a.likes,a.manual_like," +
		"a.comment_count,a.manual_comment_count,a.shares,a.manual_share" +
		" FROM video_detail a left join user_detail b on a.user_id = b.user_id where a.video_status = 4 and a.chartlet_id in ("
	for k, v := range chartletIds {
		if k == len(chartletIds)-1 {
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

	var list []*chartlet
	//循环读取结果
	for rows.Next() {
		item := &chartlet{}
		//将每一行的结果都赋值到一个item对象中
		err := rows.Scan(&item.VskitId, &item.Id, &item.VideoId, &item.VideoTitle, &item.VideoUrl, &item.UserId, &item.Views, &item.ManualView, &item.likes, &item.ManualLike, &item.CommentCount, &item.ManualCommentCount, &item.Shares, &item.ManualShare)
		if err != nil {
			fmt.Println("rows fail")
		}
		list = append(list, item)
	}
	return list
}

func getChartletList() []*chartlet {
	//执行查询语句
	rows, err := db.DB2.Query("SELECT title,uuid FROM chartlet_info where auto_id in (1130, 1124, 1143, 1142, 1140, 1135)")
	if err != nil {
		fmt.Println("查询出错了")
	}

	var list []*chartlet
	//循环读取结果
	for rows.Next() {
		item := &chartlet{}
		//将每一行的结果都赋值到一个item对象中
		err := rows.Scan(&item.ChartletTitle, &item.Id)
		if err != nil {
			fmt.Println("rows fail")
		}
		list = append(list, item)
	}
	return list
}

func getRatePc(videoIds []string) []*videoRatePc {
	getUrl := fmt.Sprintf(GetRatePcUrl, "2020-12-01", time.Now().Format("2006-01-02"))
	for _, v := range videoIds {
		if v != "" {
			getUrl += "&video_id=" + v
		}
	}

	res := Get(getUrl)
	var result *result
	err := json.Unmarshal([]byte(res), &result)
	if err != nil {
		return nil
	}

	var list []*videoRatePc
	for _, v := range result.Data {
		item := &videoRatePc{
			VideoId: v.VideoId,
			RatePc:  v.RatePc,
		}
		list = append(list, item)
	}
	return list
}

//发送GET请求
//url:请求地址
//response:请求返回的内容
func Get(url string) (response string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, error := client.Get(url)
	defer resp.Body.Close()
	if error != nil {
		panic(error)
	}

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	response = result.String()
	return
}

func removeRepeatedElement1(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

const GetRatePcUrl = "http://vskit-data-api.vskit-data.grp/video_show?dateFrom=%s&dateTo=%s"

type result struct {
	Code string  `json:"code"`
	Data []*item `json:"data"`
}

type item struct {
	Comment int     `json:"comment"`
	Nlike   int     `json:"n_like"`
	Npc     float64 `json:"n_pc"`
	Nshare  int     `json:"n_share"`
	Nview   int     `json:"n_view"`
	RatePc  float64 `json:"rate_pc"`
	VideoId string  `json:"videoId"`
}

type videoRatePc struct {
	VideoId string
	RatePc  float64
}

type chartlet struct {
	Id                 string  `json:"id"`
	ChartletTitle      string  `json:"chartlet_title"`
	VideoId            string  `json:"video_id"`
	VideoTitle         string  `json:"video_title"`
	VideoUrl           string  `json:"video_url"`
	UserId             string  `json:"user_id"`
	VskitId            string  `json:"vskit_id"`
	Views              int     `json:"views"`
	ManualView         int     `json:"manual_view"`
	likes              int     `json:"likes"`
	ManualLike         int     `json:"manual_like"`
	CommentCount       int     `json:"comment_count"`
	ManualCommentCount int     `json:"manual_comment_count"`
	Shares             int     `json:"shares"`
	ManualShare        int     `json:"manual_share"`
	TagTitle           string  `json:"tag_title"`
	RatePc             float64 `json:"rate_pc"`
}

func getAutoIdByUuid(uuid string) string {
	switch uuid {
	case "f0994b28-d159-42d6-ae62-f1b209854298":
		return "1124"
	case "4fc4d504-e489-429a-8008-3b8ea8c117b5":
		return "1130"
	case "d7e449ee-cd45-4271-a0f1-d7471993fa58":
		return "1135"
	case "79da3b8d-e170-4a86-94c8-1f28b8156b2a":
		return "1140"
	case "dfca450a-c5c4-45d7-ba59-d30d64c27a71":
		return "1142"
	case "557e0439-e559-485b-bf94-e0a03a5e0772":
		return "1143"
	}
	return ""
}
