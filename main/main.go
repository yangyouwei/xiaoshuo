package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yangyouwei/xiaoshuo/conf"
	"github.com/yangyouwei/xiaoshuo/getbookinfo"
	"log"
	"os"
)

var Db *sql.DB
var err error

//读取配置文件，初始化数据库连接
func init()  {
	var datasourcename string = conf.Mysql_conf_str.Username + ":" + conf.Mysql_conf_str.Password + "@tcp(" + conf.Mysql_conf_str.Ipaddress + ":" + conf.Mysql_conf_str.Port + ")/" + conf.Mysql_conf_str.DatabaseName
	Db, err = sql.Open("mysql", datasourcename)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

//main
func main()  {
	workmode := conf.Main_str.Mode
	switch workmode {
	case "getbookinfo":
		getbookinfo.GetBookinfo(Db)
	//case "getchapterinfo":
	//	getchapterinfo.GetChapterInfo(Db)
	//case "getcontent":
	//	getcontent.GetContent(Db)
	default:
		fmt.Println("workmod error.")
		os.Exit(1)
	}
}
