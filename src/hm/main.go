package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/opesun/goquery"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
)

type Result struct {
	Count interface{} `json:"count"`
	MemberInfo []interface{} `json:"res"`
}

type MemberInfo struct {
	level int
	nickname string
	phone string
	pk_user_main int
	popular int
	total int
	avatar string
}
var db *sql.DB
func init()  {
	var err error
	db, err = sql.Open("mysql", "root:1126254578@tcp(127.0.0.1:3306)/demo")
	if err != nil {
		fmt.Println(err)
	}
}

func main()  {
	//GetHmMember()

	GetHmMemberFromUid()


}

/**
通过ID暴力抓取用户QQ、手机信息
 */
func GetHmMemberFromUid(){
	for uid:=0;uid<7000;uid++{
		url:="http://720.huimwang.com/people?uid="+strconv.Itoa(uid)
		resp:=DownloadPage(url)
		p, err := goquery.Parse(resp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			t := p.Find(".aut-text span")
			var qq int
			var phone string
			for i:=0;i<t.Length() ;i++  {
				d:=t.Eq(i).Text()
				if ok, err := regexp.MatchString("QQ：*", d); ok && err==nil {
					qq,_=strconv.Atoi(d[5:])
				}
				if ok, err := regexp.MatchString("电话：*", d); ok && err==nil {
					phone=d[9:]
				}
			}
			if qq!=0 || phone!=""  {
				if IssetMember(uid) {
					_, err = db.Exec("update s_hm_member set qq=? where pk_user_main=?",qq,uid)
				}else{
					_, err = db.Exec("insert into s_hm_member (pk_user_main,phone,qq) values (?,?,?)",uid,phone,qq)
					fmt.Println("数据插入",uid)
				}
				if err != nil {
					fmt.Println("sql错误：", err)
					return
				}

			}
			if err != nil {
				fmt.Println("数据插入错误：", err)
			}

		}
	}
}

func IssetMember(pk_user_main int) bool {
	rows:=db.QueryRow("select pk_user_main from s_hm_member where pk_user_main=?",pk_user_main)
	var id int
	err:= rows.Scan(&id)
	if id==0 || err!=nil {
		return false
	}else{
		return true
	}
}

/**
通过json获取用户信息
 */
func GetHmMember()  {
	var url="http://720.huimwang.com/people?ajax=1&page="
	for i:=1;i<154;i++ {
		resp := DownloadPage(url + strconv.Itoa(i))
		b, _ := ioutil.ReadAll(resp.Body)
		arr := Result{}
		err := json.Unmarshal(b, &arr)
		if err != nil {
			fmt.Println(err)
		}
		var num int
		for _, data := range arr.MemberInfo {
			stmt, err := db.Prepare("insert into s_hm_member (pk_user_main,phone,nickname,level,total,popular) values (?,?,?,?,?,?)")
			defer stmt.Close()
			if err != nil {
				fmt.Println("预处理错误：", err)
				return
			}
			v := data.(map[string]interface{})

			_, err = stmt.Exec(v["pk_user_main"], v["phone"], v["nickname"], v["level"], v["total"], v["popular"])
			if err != nil {
				fmt.Println("数据插入错误：", err)
			}
			num++
		}
		fmt.Println(num)
	}
}

func DownloadPage(url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("get", url, strings.NewReader(""))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("下载页面失败：", err)
		return nil
	}
	return resp

}