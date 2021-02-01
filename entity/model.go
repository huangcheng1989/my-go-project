package entity

type LuckRes struct {
	UserId     string
	VskitId    string
	Name       string
	Phone      string
	PrizeName  string
	Country    string
	CreateTime string
}

type ActivityTalentRankLotteryRecord struct {
	Id          string `json:"id,omitempty" xorm:"not null pk comment('抽奖id') VARCHAR(40)"`
	DeviceId    string `json:"device_id,omitempty" xorm:"not null default '' comment('device_id') VARCHAR(40)"`
	UserId      string `json:"user_id,omitempty" xorm:"not null default '' comment('用户id') VARCHAR(40)"`
	PrizeId     int    `json:"prize_id,omitempty" xorm:"not null default 0 comment('奖品id 0:未中奖') INT(2)"`
	Country     string `json:"country,omitempty" xorm:"not null default '' comment('二位国家码') VARCHAR(5)"`
	Pcc         string `json:"pcc,omitempty" xorm:"default '' comment('手机区号') VARCHAR(20)"`
	Phone       string `json:"phone,omitempty" xorm:"default '' comment('用户手机号') VARCHAR(20)"`
	CreatedTime int64  `json:"created_time,omitempty" xorm:"not null default 0 comment('创建时间') BIGINT(20)"`
}

type UserDetail struct {
	UserId  string `json:"user_id,omitempty" xorm:"not null pk comment('用户id') VARCHAR(100)"`
	VskitId string `json:"vskit_id,omitempty" xorm:"not null comment('vskit id') VARCHAR(30)"`
	Name    string `json:"name,omitempty" xorm:"comment('用户昵称') VARCHAR(128)"`
}

type ActivityStatusCollection struct {
	Id           string `json:"id,omitempty" xorm:"not null pk default '' comment('id') VARCHAR(100)"`
	UserId       string `json:"user_id,omitempty" xorm:"not null default '' comment('user_id') VARCHAR(100)"`
	DeviceId     string `json:"device_id,omitempty" xorm:"not null default '' comment('device_id') VARCHAR(100)"`
	Type         int    `json:"type,omitempty" xorm:"not null default 0 comment('类型: 1.Meme; 2.News; 3.Record life; 4.Blessings ; 5.Others') TINYINT(1)"`
	TypeExtra    string `json:"type_extra,omitempty" xorm:"not null default '' comment('额外类型') VARCHAR(50)"`
	Source       int    `json:"source,omitempty" xorm:"not null default 0 comment('来源: 1.Whatsapp; 2.Instagram; 3.Snapchat; 4.Others') TINYINT(1)"`
	SourceExtra  string `json:"source_extra,omitempty" xorm:"not null default '' comment('额外来源') VARCHAR(50)"`
	FileLink     string `json:"file_link,omitempty" xorm:"not null default '' comment('文件链接') VARCHAR(255)"`
	UploadReason string `json:"upload_reason,omitempty" xorm:"not null default '' comment('上传原因') VARCHAR(255)"`
	Pcc          string `json:"pcc,omitempty" xorm:"not null default '' comment('国家电话码') VARCHAR(10)"`
	Phone        string `json:"phone,omitempty" xorm:"not null default '' comment('手机号码') VARCHAR(20)"`
	Country      string `json:"country,omitempty" xorm:"not null default '' comment('国家') VARCHAR(20)"`
	Selected     int    `json:"selected,omitempty" xorm:"not null default 0 comment('是否被选中：0:否，1:是') TINYINT(1)"`
	CreateTime   int    `json:"create_time,omitempty" xorm:"not null default 0 comment('创建时间') INT(13)"`
	UpdateTime   int    `json:"update_time,omitempty" xorm:"not null default 0 comment('更新时间') INT(13)"`
}

type StatusRes struct {
	Id         string
	Type       string
	Source     string
	Link       string
	Reason     string
	Phone      string
	Country    string
	CreateTime int
	Selected   int
}

type TalentItem struct {
	UserId   string `json:"user_id"`
	VoteNum  int    `json:"vote_num"`
	Sort     int    `json:"sort"`
	UserName string `json:"user_name"`
	VskitId  string `json:"vskit_id"`
}

type ActivityTalentRank struct {
	UserId        string `json:"user_id,omitempty" xorm:"not null pk default '' comment('用户ID') VARCHAR(100)"`
	VskitId       string `json:"vskit_id,omitempty" xorm:"not null default '' comment('vskit id') VARCHAR(30)"`
	Name          string `json:"name,omitempty" xorm:"not null default '' comment('用户昵称') VARCHAR(100)"`
	Type          int    `json:"type,omitempty" xorm:"not null default 0 comment('榜单类型：1.人气创作者；2.新晋创作者') TINYINT(1)"`
	VoteCount     int    `json:"vote_count,omitempty" xorm:"not null default 0 comment('投票数') INT(15)"`
	FollowerCount int    `json:"follower_count,omitempty" xorm:"not null default 0 comment('粉丝数') INT(15)"`
	Pcc           string `json:"pcc,omitempty" xorm:"not null default '' comment('国家电话码') VARCHAR(5)"`
	Phone         string `json:"phone,omitempty" xorm:"not null default '' comment('用户手机号') VARCHAR(20)"`
	CreateTime    int    `json:"create_time,omitempty" xorm:"not null default 0 comment('创建时间') INT(10)"`
	UpdateTime    int    `json:"update_time,omitempty" xorm:"not null default 0 comment('更新时间') INT(10)"`
	Sort          int64
}

type TalentItemList []*ActivityTalentRank

func (til TalentItemList) Len() int {
	return len(til)
}

func (til TalentItemList) Swap(i, j int) {
	til[i], til[j] = til[j], til[i]
}

// 按照Sort（排序从小到大）排序
func (til TalentItemList) Less(i, j int) bool {
	if til[i].Sort < til[j].Sort {
		return true
	}
	return false
}
