package jiexi_html

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

//https://api.bilibili.com/x/player/playurl?avid=548692744&cid=428179835

//获取html主体提取aid和json
func Get_json(url string) (json_s []byte, err error) {
	reque_html, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("获取html_创建请求失败:", err)
		return
	}
	client := http.Client{}
	reque_html.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36 Edg/94.0.992.38")
	resp, err := client.Do(reque_html)
	if err != nil {
		fmt.Println("获取html_连接出错:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("获取html_连接失败:", resp.StatusCode)
		err = errors.New("获取html_连接失败")
		return
	}
	fmt.Println("获取html_连接成功:", resp.StatusCode)
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("获取html_读取网页主体失败:", err)
		return
	}
	//fmt.Printf("%s\n", res)
	re_m_s := "window.__INITIAL_STATE__=([^;]+);"
	re := regexp.MustCompile(re_m_s)
	//re_s := re.FindAllSubmatch(res, 1)

	//fmt.Printf("匹配字符串：%s\n", re_s[0][1])

	//re_m_aid := `"aid":([0-9]+),`
	//re = regexp.MustCompile(re_m_aid)
	//aid = re.FindAllSubmatch(res, 1)[0][1]
	//fmt.Printf("aid:%s\n", aid)

	re_m_json := `"pages":([^]]+)],`
	re = regexp.MustCompile(re_m_json)
	json_s = append(re.FindAllSubmatch(res, 1)[0][1], []byte{']'}...)
	//fmt.Printf("json:%s\n", json_s)
	return
}

//解析页数的结构体
type Pages struct {
	Cid        int         `json:"cid"`
	Page       int         `json:"page"`
	From       string      `json:"from"`
	Part       string      `json:"part"`
	Duration   int         `json:"duration"`
	Vid        string      `json:"vid"`
	Weblink    string      `json:"weblink"`
	Dimension  interface{} `json:"dimension"`
	FirstFrame string      `json:"first_frame"`
}

//获取页数相对应的cid切片
func Get_aid_cid_pages(url string) (pages []Pages, err error) {
	jso, _ := Get_json(url)
	//fmt.Printf("jso:%s\n", jso)
	pages = make([]Pages, 0)
	err = json.Unmarshal(jso, &pages)
	if err != nil {
		fmt.Println("json解析失败", err)
		return
	}
	fmt.Println("json解析成功")
	return
}

func Get_video_url(bvid string, pages []Pages, i int) (video_url string, err error) {
	fmt.Println("使用的cid:", pages[i].Cid)
	url := "https://api.bilibili.com/x/player/playurl?" + "cid=" + strconv.Itoa(pages[i].Cid) + "&bvid=" + bvid
	reque, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求头失败：", err)
		return
	}
	client := http.Client{}
	reque.Header.Add("referer", "https://www.bilibili.com")
	reque.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36 Edg/94.0.992.38")
	rep, err := client.Do(reque)
	if err != nil {
		fmt.Println("连接失败", err)
		return
	}
	if rep.StatusCode != http.StatusOK {
		fmt.Println("连接失败,状态码：", rep.StatusCode)
		return
	}
	fmt.Println("连接成功,状态码：", rep.StatusCode)
	rep_byte, err := io.ReadAll(rep.Body)
	if err != nil {
		fmt.Println("读取返回的json失败：", err)
		return
	}
	re_url_s := `"url":"([^"]+)"`
	re := regexp.MustCompile(re_url_s)
	video_url = string(re.FindAllSubmatch(rep_byte, 1)[0][1])
	return
}

//获取视频下载链接列表
func Get_video_url_liet(bvid string, pages []Pages) (video_url_list []string, err error) {
	for i := 0; i < len(pages); i++ {
		v_url, err := Get_video_url(bvid, pages, i)
		if err != nil {
			return []string{}, err
		}
		video_url_list = append(video_url_list, v_url)
	}
	return
}
