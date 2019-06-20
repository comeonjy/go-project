package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/opesun/goquery"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var db *sql.DB

var IdList = make(chan string, 10000)

func init() {
	var err error
	db, err = sql.Open("mysql", "root:1126254578@tcp(127.0.0.1:3306)/demo")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	//go func() {
	//	for {
	//		GetUid()
	//		time.Sleep(5 * time.Second)
	//	}
	//}()

	for i := 0; i < 10; i++ {
		go Worker()
	}

	for   {
		time.Sleep(1*time.Hour)
	}

}

func GetUid() {
	url := "https://work.3d66.com/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	p, err := goquery.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		t := p.Find(".lis .work-info .info-link .contact")
		for i := 0; i < t.Length(); i++ {
			member_id := t.Eq(i).Attr("data-user_id")
			IdList <- member_id
		}
	}
}

func IssetMember(member_id string) bool {
	rows := db.QueryRow("select id from s_threed66 where member_id=?", member_id)
	var id int
	err := rows.Scan(&id)
	if id == 0 || err != nil {
		return false
	} else {
		return true
	}
}

func Worker() {
	for {
		member_id := <-IdList
		GetInfo(member_id)
		fmt.Println(len(IdList))
	}
}

func GetInfo(uid string) {
	if IssetMember(uid) {
		return
	}
	url := "https://work.3d66.com/index/work_introduce/index/user_id/" + uid
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	p, err := goquery.ParseString(string(b))
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		t := p.Find(".excell-list")
		qq := t.Find(".excell-list1 p").Text()
		mobile := t.Find(".excell-list2 p").Text()
		url := t.Find(".excell-list3 p").Text()
		name := p.Find(".studio-name h3").Text()

		_, _ = db.Exec("insert into s_threed66 (member_id,name,mobile,qq,url,create_time) values (?,?,?,?,?,?)", uid,name, mobile, qq, url, time.Now().Unix())
	}
}

func DownloadPage(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
