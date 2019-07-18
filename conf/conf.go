package conf

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"path/filepath"
	"strconv"
)

type mainS struct {
	Concurrent	int
	Mode string
	Filepath string
}

type Mysql_conf struct {
	Username string
	Password string
	Ipaddress string
	Port	string
	DatabaseName	string
}

type Chapter1 struct {
	Rules []string
}
type Chapter2 struct {
	Rules []string
}

type Headerrule struct {
	Hr string
}

var HR Headerrule
var Chapterrules1 Chapter1
var Chapterrules1slice *[]string
var Chapterrules2 Chapter2
var Chapterrules2slice *[]string
var Mysql_conf_str Mysql_conf
var Main_str mainS

func init() {
	cfg, err := goconfig.LoadConfigFile("conf")
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
		panic(err)
	}
	HR.Headerrules(cfg,err)
	Chapterrules1.Getchapterrules(cfg)
	Chapterrules2.Getchapterrules2(cfg)
	Mysql_conf_str.Mysql_fun(cfg,err)
	Main_str.main_fun(cfg,err)
}

func (this *Headerrule)Headerrules(c *goconfig.ConfigFile,err error )  {
	r,err := c.GetValue("headerrule","rules")
	if err != nil {
		log.Fatalf("无法获取键值section（%s）：%s", "headerrule", err)
		panic(err)
	}
	this.Hr = r
}

func (this *mainS)main_fun(c *goconfig.ConfigFile,err error)  {
	n,err := c.GetValue("main","concurrent")
	if err != nil {
		log.Fatalf("无法获取键值section（%s）：%s", "concurrent", err)
		panic(err)
	}
	this.Concurrent,err = strconv.Atoi(n)
	if err != nil {
		log.Fatalf("%s）：%s,无效", "concurrent", err)
		panic(err)
	}

	a,err := c.GetValue("main","filepath")
	this.Filepath,err = filepath.Abs(a)
	if err != nil {
		log.Fatalf("%s）：%s,无效", "filepath", err)
		panic(err)
	}

	this.Mode,err = c.GetValue("main","mode")
	if err != nil {
		log.Fatalf("%s）：%s,无效", "mode", err)
		panic(err)
	}
}



func (this *Chapter1) Getchapterrules(c *goconfig.ConfigFile) {
	this.Rules = c.GetKeyList("chapter_rules1")
	for _,v := range this.Rules {
		fmt.Println(v)
		r,_ := c.GetValue("chapter_rules1", v)
		var b []string
		b = append(b,r)

	}
}

func (this *Chapter2) Getchapterrules2(c *goconfig.ConfigFile) {
	this.Rules = c.GetKeyList("chapter_rules2")
	for _,v := range this.Rules {
		r,_ := c.GetValue("chapter_rules2", v)
		a := Chapterrules2slice
		*a = append(*a,r)
	}
}

func (this *Mysql_conf)Mysql_fun(c *goconfig.ConfigFile,err error) {
	this.Username, err = c.GetValue("main", "username")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "username", err)
		panic(err)
	}

	this.Password, err = c.GetValue("main", "password")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "password", err)
		panic(err)
	}

	this.Ipaddress, err = c.GetValue("main", "addr")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "addr", err)
		panic(err)
	}

	this.Port, err = c.GetValue("main", "port")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "prot", err)
		panic(err)
	}

	this.DatabaseName, err = c.GetValue("main", "databasename")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "databasename", err)
		panic(err)
	}
}


