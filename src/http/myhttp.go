package http

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var alivesum, titlesum int

//socket进行对网站的连接
func MyHttpSocket(ip string) string {
	//fmt.Println(ip)
	var result string

	//socket tcp连接,超时时间
	conn, err := net.DialTimeout("tcp", ip, 3*time.Second)

	if err != nil {

		//fmt.Println(err)
		return ""
	}

	//发送内容
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + ip + "\r\n\r\n"))

	if err != nil {
		return ""
	}

	//读取时间2秒超时
	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	if err != nil {
		return ""
	}

	//最多只读8192位,一般来说有title就肯定已经有了
	reply := make([]byte, 8192)
	_, err = conn.Read(reply)

	if err != nil {
		return ""
	}

	err = conn.Close()
	if err != nil {
		return ""
	}
	html := string(reply)

	//获取状态码
	status := GetStatusCode(html)

	//如果是400可能是因为没有用https
	if status == "400" {
		result = SystemHttp(ip)
		return result
	}

	//正则匹配title
	r, _ := regexp.Compile("<title>(.*)</title>")

	res := r.FindStringSubmatch(html)

	if len(res) < 2 {

		if (strings.Count(html, "") - 1) > 20 {
			result = "[+]" + ip + "  open ---------" + string([]byte(html)[:14])
			alivesum++
			return result
		} else {
			result = "[+]" + ip + "  open ---------" + html
			alivesum++
			return result
		}
		//fmt.Println(result)

	}

	result = "[+]" + ip + "  open ---------" + res[1]

	alivesum++
	titlesum++

	return result

}

//使用封装好了http
func SystemHttp(ip string) string {
	var result string
	ip = "https://" + ip

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   2 * time.Second,
	}
	resp, err := c.Get(ip)

	if err != nil {
		return ""
	}

	reply, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	html := string(reply)

	if err != nil {

		return ""
	}

	r, _ := regexp.Compile("<title>(.*)</title>")

	res := r.FindStringSubmatch(html)

	if len(res) < 2 {
		result = "[+]" + ip + "  open ---------"
		//fmt.Println(result)
		alivesum++
		return result
	}

	result = "[+]" + ip + "  open ---------" + res[1]

	alivesum++
	titlesum++

	return result
}

func GetStatusCode(html string) string {
	http1 := strings.Split(html, "\n")[0]
	statusC := strings.Split(http1, " ")
	if len(statusC) > 2 {
		statusCode := statusC[1]
		return statusCode
	}

	return ""
}

func OutputAliveSum() {
	fmt.Println("AliveSum: " + strconv.Itoa(alivesum))
}

func OutputTitleSum() {
	fmt.Println("TitleSum: " + strconv.Itoa(titlesum))
}
