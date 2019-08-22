package fetcher

import (
	"fmt"
	"github.com/gpmgo/gopm/modules/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

//var rateLimit = time.Tick(1000 * time.Microsecond)
//TODO 不同返回值的fetcher
func Fetcher(urls string) ([]byte, error) {
	//<-rateLimit
	fmt.Println(urls)
	resp := getResp(urls, returnIp())
	return ioutil.ReadAll(resp.Body)
}

/**
* 返回response
*/
func getResp(urls string, ip string) *http.Response {
	var response *http.Response
	for {
		//ip="local"
		ip = returnIp()
		request, _ := http.NewRequest("GET", urls, nil)
		//随机返回User-Agent 信息
		request.Header.Set("User-Agent", getAgent())
		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		request.Header.Set("Connection", "keep-alive")
		proxy, err := url.Parse(ip)
		if err != nil {
			log.Error("URL Parse:", err)
			continue
		}
		//设置超时时间
		timeout := time.Duration(10 * time.Second)
		//fmt.Printf("使用代理:%s\n", proxy)
		client := &http.Client{}
		if ip != "local" {
			client = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxy),
				},
				Timeout: timeout,
			}
		}

		response, err = client.Do(request)
		if err != nil || response.StatusCode != 200 {

			log.Error("遇到了错误-并切换ip %s\n", err)
			continue

		}
		break
	}
	return response
}

/**
* 随机返回一个User-Agent
*/
func getAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(agent)
	return agent[r.Intn(len)]
}

func returnIp() string {
	agent := [...]string{
		"local",
		//"http://113.120.60.243:22752",
		//"http://120.83.102.255:808",
		//"http://171.41.120.91:9999",
		//"http://58.218.205.42:17768",
		//"http://114.228.6.146:4376",
		//"http://182.34.194.4:4376",
		//"http://58.218.205.42:16479",
		//"http://58.218.205.42:13275",
		//"http://117.63.140.237:4376",
		//"http://182.34.192.251:4376",
		//"http://180.116.159.82:4375",
		//"http://180.97.250.136:18555",
		//"http://180.97.250.136:18528",
		//"http://180.115.116.53:4375",
		//"http://180.97.250.136:18562",
		//"http://182.101.203.13:4313",
		//"http://182.109.88.65:4321",
		//"http://106.5.28.73:4321",
		//"http://180.97.250.136:16538",
		//"http://180.97.250.136:13542",
		//"http://124.152.185.142:4329",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return agent[r.Intn(len(agent))]
}
