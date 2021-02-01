package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//Db数据库连接池
var DB1 *sql.DB
var DB2 *sql.DB

//注意方法名大写，就是public
func InitDB1(env string) {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"

	var path string
	if env == "prod" {
		//正式
		path = "palmstore_vskit:PAlmstore2018@tcp(palmstore-vskit-aurora-master.csidlk2hdfqg.eu-west-1.rds.amazonaws.com:3306)/vskit_activity?interpolateParams=False&charset=utf8mb4"
	} else if env == "test" {
		//测试
		path = "vshow:aosjdkfAdfjijaDFIJAsjdkhf2837492asjhdf@tcp(db.mylichking.com:3306)/vskit_activity?interpolateParams=False&charset=utf8mb4"
	}

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB1, _ = sql.Open("mysql", path)
}

func InitDB2() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := "palmstore_vskit:PAlmstore2018@tcp(palmstore-vskit-db.csidlk2hdfqg.eu-west-1.rds.amazonaws.com:3306)/palmstore_vskit?interpolateParams=False&charset=utf8mb4"

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB2, _ = sql.Open("mysql", path)
}
