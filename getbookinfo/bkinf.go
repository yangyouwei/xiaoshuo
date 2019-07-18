package getbookinfo

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/yangyouwei/xiaoshuo/conf"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)
type Bookinfo struct {
	Bookid          int64  `db:"id"`
	Bookname        string `db:"booksName"`
	Sourcesfilename string `db:"sourcesfilename"`
	RegexRules      string  `db:"regexRules"`
}
var db *sql.DB
var filenamech = make(chan string, 10)
var cr = `^\s*.*第(\d+|[一二三四五六七八九十百千万]+)章.*[^。]*$`
var cr1 = conf.Chapterrules1.Rules
var cr2 = conf.Chapterrules2.Rules

func GetBookinfo(Db *sql.DB)  {
	db = Db
	pathname, err := filepath.Abs(conf.Main_str.Filepath)
	if err != nil {
		fmt.Println("path error")
		return
	}
	concurrenc := conf.Main_str.Concurrent
	wg := sync.WaitGroup{} //控制主程序等待，以便goroutines运行完
	wg.Add(concurrenc + 1)
	go func(wg *sync.WaitGroup, filenamech chan string) {
		GetAllFile(pathname, filenamech)
		close(filenamech) //关闭通道，以便读取通道的程序知道通道已经关闭。
		wg.Done()         //一定在函数的内部的最后一行运行。否则可能函数没有执行完毕。
	}(&wg, filenamech)
	for i := 0; i < concurrenc; i++ {
		go func(wg *sync.WaitGroup, filenamech chan string) {
			for {
				filename, isclose := <-filenamech
				if !isclose { //判断通道是否关闭，关闭则退出循环
					break
				}
				dosomewrork(filename)
			}
			wg.Done()
		}(&wg, filenamech)
	}
	wg.Wait()
}

//获取文件名
func GetAllFile(pathname string, fn_ch chan string) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			GetAllFile(fullDir, fn_ch)
			if err != nil {
				fmt.Println("read dir fail:", err)
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			fn_ch <- fullName
		}
	}
}

//文档处理函数
func dosomewrork(fp string) {
	b := Bookinfo{}
	b.getinfo(fp)
	b.insert(db)
}

func (b *Bookinfo) getinfo(fp string) {
	//bookname
	bn := strings.Split(filepath.Base(fp), ".")
	bookname := bn[2]
	b.Bookname = bookname
	b.Sourcesfilename = fp

	//匹配章节规则
	var isok bool
	var rules string
	fi, err := os.Open(fp)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer fi.Close()
	n := 0
	br := bufio.NewReader(fi)
	for i := 0; i < 1000; i++ {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		//fmt.Println(string(a))
		isok,rules = getrule(string(a))
		if isok {
			a := &n
			*a = *a + 1
		}
	}
	if n >= 3 {
		b.RegexRules = rules
	}
}

func getrule(s string) (isok bool,r string) {
	//匹配  xxx第xx章xxx  类似章节名称
	isok1 , err := regexp.Match(cr,[]byte(s))
	if err != nil {
		fmt.Println(err)
	}
	//如果能匹配，则匹配更详细的规则作为该小说的规则
	if isok1 {
		//fmt.Println(s,":",cr)
		for _,v := range *conf.Chapterrules1slice {

			isok2 , _ := regexp.Match(v,[]byte(s))
			if isok2 {
				fmt.Println(v)
				return true, v
			}
		}
	}
	//如果都匹配不上，使用规则2中的规则。
	for _,v := range *conf.Chapterrules2slice {
		isok2 , _ := regexp.Match(v,[]byte(s))
		if isok2 {
			fmt.Println(v)
			return true, v
		}
	}
	return false, ""
}

//book信息写入数据库
func (b *Bookinfo) insert(db *sql.DB) {
	stmt, err := db.Prepare(`INSERT books ( booksName,sourcesfilename,regexRules) VALUES (?,?,?)`)
	check(err)

	res, err := stmt.Exec(b.Bookname,b.Sourcesfilename,b.RegexRules)
	check(err)

	id, err := res.LastInsertId() //必须是自增id的才可以正确返回。
	check(err)
	defer stmt.Close()

	idstr := fmt.Sprintf("%v", id)
	fmt.Println(idstr)
	stmt.Close()
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func userules(r string,s string) string {
	reg := regexp.MustCompile(r)
	result := reg.FindAllStringSubmatch(s,-1)
	s = result[0][0]
	return s
}

