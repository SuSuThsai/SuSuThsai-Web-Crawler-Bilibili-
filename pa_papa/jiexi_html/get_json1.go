package jiexi_html

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

func GetURL(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf("ERROR:get url:%s", req.URL)
		return ""
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.128 Safari/537.36 Edg/89.0.774.77")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("ERROR:模仿浏览器登陆失败！URL:%s", req.URL)
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("状态码出错:%v", resp.StatusCode)
		return ""
	}
	fmt.Println("状态码为：", resp.StatusCode)
	//defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("读取出错!", err)
	}
	return string(result)
}

func Get_image_url_liet(context string) []string {
	re := regexp.MustCompile(`<img[\s\S]+?src="([^"]+)"`)
	result := re.FindAllStringSubmatch(context, -1)
	imageURL := make([]string, 0)
	fmt.Printf("一共：%d 张.\n", len(result))
	for _, m := range result {
		//fmt.Printf("%s %s\n",m[1],GetFileName(m[1]))
		imageURL = append(imageURL, m[1])
		//imageAlt=append(imageAlt,m)
	}
	return imageURL
}
