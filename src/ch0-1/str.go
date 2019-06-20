package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/opesun/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
)

type Member struct {
	Id         int    `db:"id"`
	Nickname   string `db:"nickname"`
	Account    string `db:"account"`
	Password   string `db:"password"`
	Age        int    `db:"age"`
	CreateTime int    `db:"create_time"`
}

var db *sql.DB

var MemberInfo = make(chan map[int]map[string]interface{}, 100000)
var UrlList = make(chan string,100000)

func init() {
	var err error
	db, err = sql.Open("mysql", "root:1126254578@tcp(127.0.0.1:3306)/demo")
	if err != nil {
		fmt.Println(err)
	}


}

//寻找最长不重复子串
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//go ReScanUrl()

	go ReadListUrl()

	for i:=0;i<100 ;i++  {
		go GoGetMemberInfo()
	}


	for i:=0;i<100 ;i++  {
		go GoAddMemberRows()
	}






	go func() {
		for  {
			fmt.Println("len(MemberInfo)=",len(MemberInfo),"len(UrlList)=",len(UrlList))
			time.Sleep(1*time.Second)
		}
	}()

	for  {
		time.Sleep(10*time.Second)
	}



}

func GoAddMemberRows()  {
	for{
		var data map[int]map[string]interface{}
		for i:=0;i<1000;i++  {
			info:=<-MemberInfo
			data = map[int]map[string]interface{}{
				i: info[0],
			}
		}
		AddMemberRows(data)
	}
}

func GoGetMemberInfo()  {
	for{
		url:=<-UrlList
		GetMemberInfo(DownloadPage(url))
	}
}
var client = &http.Client{}
func DownloadPage(url string) *http.Response {
	//var url = "http://album.zhenai.com/u/24756291"
	//client := &http.Client{}
	req, err := http.NewRequest("get", url, strings.NewReader(""))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("下载页面失败：", err)
		return nil
	}
	return resp
}

/**
读取member info
 */
func GetMemberInfo(resp *http.Response)  {
	if resp==nil {
		return
	}
	p, err := goquery.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		t := p.Find("script")
		for i := 0; i < t.Length(); i++ {
			d := t.Eq(i).Text()
			if ok, err := regexp.MatchString("^window.__INITIAL_STATE__=*", d); ok && err==nil {
				start := strings.Index(d, "objectInfo")
				end := strings.Index(d, ",\"interest")
				if start<0 || end<0 {
					break
				}
				var arr map[string]interface{}
				b := []byte(d[start+12:end])
				if b == nil {
					fmt.Println("用户信息隐藏")
					break
				}
				err := json.Unmarshal(b, &arr)
				if err != nil {
					fmt.Println(err)
				}
				memberid:=int(arr["memberID"].(float64))
				members := map[int]map[string]interface{}{
					0: {"memberID": memberid, "info": d[start+12 : end],},
				}
				MemberInfo <- members
				break
			}
		}
	}
}

/**
将用户信息写入数据库
 */
func AddMemberRows(data map[int]map[string]interface{}) (nums int) {
	stmt, err := db.Prepare("insert into s_member (memberID,info,create_time) values (?,?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("预处理错误：", err)
		return
	}
	for _, v := range data {

		//if IssetMember(v["memberID"].(int)) {
			//fmt.Println("数据重复",v["memberID"])
			//break
		//}

		row, err := stmt.Exec(v["memberID"], v["info"],time.Now().Unix() )
		if err != nil {
			fmt.Println("数据插入错误：", err)
		}
		if count, err := row.RowsAffected(); err == nil {
			nums += int(count)
			//fmt.Println("数据插入：",v["memberID"])
		}
	}

	return
}


func IssetMember(memberId int) bool {
	rows:=db.QueryRow("select id from s_member where memberID=?",memberId)
	var id int
	err:= rows.Scan(&id)
	if id==0 || err!=nil {
		return false
	}else{
		return true
	}
}

func GetListInfo(resp *http.Response)  {
	if resp==nil {
		return
	}
	p, err := goquery.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		t := p.Find("a")
		var data = make(map[int]map[string]interface{})
		for i := 0; i < t.Length(); i++ {
			url := t.Eq(i).Attr("href")
			if ok, _ := regexp.MatchString("^http://www.zhenai.com/zhenghun/*", url); ok {
				url_map := map[string]interface{}{"url":url}
				data[i]=url_map
			}
		}
		AddListRows(data)
	}
}

func AddListRows(data map[int]map[string]interface{}) (nums int) {
	stmt, err := db.Prepare("insert into s_url_list (url,create_time) values (?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("预处理错误：", err)
		return
	}
	for _, v := range data {
		if IssetList(v["url"].(string)) {
			//fmt.Println("数据重复",v["url"])
			continue
		}
		row, err := stmt.Exec(v["url"],time.Now().Unix() )
		if err != nil {
			fmt.Println("数据插入错误：", err)
		}
		if count, err := row.RowsAffected(); err == nil {
			nums += int(count)
			//fmt.Println("数据插入：",v["url"])
		}
	}


	return
}


func IssetList(url string) bool {
	rows:=db.QueryRow("select id from s_url_list where url=?",url)
	var id int
	err:= rows.Scan(&id)
	if id==0 || err!=nil {
		return false
	}else{
		return true
	}
}

/**
将列表链接中member url写入管道
 */
func GetUrlInfo(resp *http.Response)  {
	if resp==nil {
		return
	}
	p, err := goquery.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
	} else {
		t := p.Find("a")
		for i := 0; i < t.Length(); i++ {
			url := t.Eq(i).Attr("href")
			if ok, _ := regexp.MatchString("^http://album.zhenai.com/u/*", url); ok {
				id,_:=strconv.Atoi(url[26:])
				if IssetMember(id) {
					continue
				}
				//fmt.Println(url)
				UrlList<-url
			}
		}
	}
}

/**
每秒取一个列表连接
收集列表中所有member信息
 */
func ReadListUrl()  {
	for  {
		rows,err:=db.Query("select * from s_url_list where id > 1000")
		if err != nil {
			fmt.Println(err)
		}else{
			var id int
			var url string
			var create_time string
			var state int
			for rows.Next() {
				err := rows.Scan(&id,&url,&create_time,&state)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(id)
				go func() {
					for i:=2;i<8;i++{
						GetUrlInfo(DownloadPage(url+"/"+strconv.Itoa(i)))
					}
				}()
				GetUrlInfo(DownloadPage(url+"/8"))
			}

			//data:=FetchRows(rows)
			//for _,v:=range data{
			//	fmt.Println(v["id"])
			//	GetUrlInfo(DownloadPage(v["url"]+"/2"))
			//	GetUrlInfo(DownloadPage(v["url"]+"/3"))
			//	GetUrlInfo(DownloadPage(v["url"]+"/4"))
			//	GetUrlInfo(DownloadPage(v["url"]+"/5"))
			//	GetUrlInfo(DownloadPage(v["url"]+"/6"))
			//}
		}
		fmt.Println("------------------------------------------------------")
	}
}

func ReScanUrl()  {
	rows,err:=db.Query("select * from s_url_list")
	if err != nil {
		fmt.Println(err)
	}

	var id int
	var url string
	var create_time string
	var state int
	for rows.Next() {
		err := rows.Scan(&id,&url,&create_time,&state)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go GetListInfo(DownloadPage(url))
	}

	//data:=FetchRows(rows)
	//for _,v:=range data{
	//	GetListInfo(DownloadPage(v["url"]))
	//	time.Sleep(time.Second)
	//}
}


func FetchRows(rows *sql.Rows) map[int]map[string]string {
	cols, _ := rows.Columns()
	scans := make([]interface{}, len(cols))
	vals := make([][]byte, len(cols))
	for i, _ := range vals {
		scans[i] = &vals[i]
	}

	var i = 0
	data := make(map[int]map[string]string)
	for rows.Next() {
		err := rows.Scan(scans...)
		if err != nil {
			fmt.Println(err)
		}
		row := make(map[string]string, len(cols))
		for i, v := range vals {
			row[cols[i]] = string(v)
		}
		data[i] = row
		i++
	}
	return data
}


func Utf8Body(resp *http.Response) *transform.Reader {
	reader := bufio.NewReader(resp.Body)
	bytes, e := reader.Peek(1024)
	if e != nil {
		fmt.Println(e)
	}
	encodig, _, _ := charset.DetermineEncoding(bytes, "")

	utf8Reader := transform.NewReader(resp.Body, encodig.NewDecoder())
	return utf8Reader
}