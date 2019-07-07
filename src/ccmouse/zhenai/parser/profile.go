package parser

import (
	"bytes"
	"ccmouse/engine"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

type Profile struct {
	Name string
	City string
	Age int32
	Edu string
	MaritalStatus string
	Height string
	Income string
	Spouse []string
}

func ParseProfile(contents []byte) engine.ParseResult {
	var result engine.ParseResult

	nodes, err:= goquery.NewDocumentFromReader(bytes.NewReader(contents))
	if err!=nil {
		panic(err)
	}
	var profile Profile

	nickname:=nodes.Find("h1.nickName").Text()
	profile.Name=nickname

	info:=nodes.Find(".des.f-cl").Text()
	arr:=strings.Split(info,` | `)
	if len(arr)!=6 {
		log.Print("Profile Error:",info)
		return result
	}
	profile.City=arr[0]

	runes:=[]rune(arr[1])
	age, err := strconv.Atoi(string(runes[:len(runes)-1]))
	if err!=nil {
		age=0
	}
	profile.Age=int32(age)

	profile.Edu=arr[2]
	profile.MaritalStatus=arr[3]
	profile.Height=arr[4][:len(arr[4])-2]
	profile.Income=arr[5]

	var Spouse []string
	nodes.Find(".gray-btns > .m-btn").Each(func(i int, selection *goquery.Selection) {
		Spouse=append(Spouse,selection.Text())
	})
	profile.Spouse=Spouse

	result.Items=append(result.Items,profile)


	return result
}
