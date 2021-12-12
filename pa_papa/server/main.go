package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"pa_papa/db"
	"pa_papa/jiexi_html"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	s_login = ""
	m_urls = make(map[string][]string)
	p_urls = make(map[string][]string)
	admin = false
	s_down = ""
	//filr_path:=filepath.Dir(os.Args[0])
	//fmt.Println(filr_path)
	path_s, _ := filepath.Abs(os.Args[0])
	path_ss := strings.Split(path_s, "/")
	path = strings.Join(path_ss[:len(path_ss)-2], "/")
	fmt.Println(path)
}

func main() {
	db.DB_init()
	s_login = ""
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(path+"/tmpl"))))
	http.HandleFunc("/login", login)
	http.HandleFunc("/login_que", login_que)
	http.HandleFunc("/pa", pa)
	http.HandleFunc("/pa_bv_post", pa_bv_download)
	http.HandleFunc("/pa_bv_list", pa_bv_list)
	http.HandleFunc("/pa_picture", pa_picture)
	http.HandleFunc("/pa_cv_post", pa_cv_download)
	//http.HandleFunc("/pa_cv_post2", Delt)
	http.HandleFunc("/pa_cv_list", pa_cv_list)
	//http.HandleFunc("/papa", papa)
	http.ListenAndServe(":9999", nil)
}

var s_login string
var s_down string
var m_urls map[string][]string
var p_urls map[string][]string
var p_cv_num string
var admin bool
var path string

//登陆页面
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(path)
	tmpl, err := template.ParseFiles(path + "/tmpl/Yamada.tmpl")
	if err != nil {
		fmt.Println("解析失败")
	}
	tmpl.Execute(w, s_login)
}

//判断登陆页面表单
func login_que(w http.ResponseWriter, r *http.Request) {
	zhanghao := r.FormValue("zhanghao")
	mima := r.FormValue("mima")
	if len(zhanghao) == 0 || len(mima) == 0 {
		s_login = "账号密码不能为空"
		fmt.Println("账号密码不能为空")
		http.Redirect(w, r, "/login", http.StatusFound)
		//tmpl.Execute(w, "账号密码不能为空")
		return
	}
	mima_zhen, err := db.GetMima(db.DB, zhanghao)
	if err != nil {
		s_login = "账号不存在"
		fmt.Println("账号不存在")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if mima != mima_zhen {
		s_login = "密码错误"
		fmt.Println("密码错误")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
		return
	}
	admin = true
	http.Redirect(w, r, "/pa", http.StatusFound)
}

func pa(w http.ResponseWriter, r *http.Request) {
	if !admin {
		fmt.Println("不是管理员不给进")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	tmpl, err := template.ParseFiles(path + "/tmpl/files.tmpl")
	if err != nil {
		fmt.Println("解析失败")
	}
	s_down = ""
	tmpl.Execute(w, s_down)
}
func pa_picture(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(path + "/tmpl/pictrue_files.tmpl")
	if err != nil {
		fmt.Println("解析失败")
	}
	s_down = ""
	tmpl.Execute(w, s_down)
}

func pa_bv_list(w http.ResponseWriter, r *http.Request) {
	if !admin {
		fmt.Println("不是管理员不给进")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	bv := r.FormValue("BV")
	if len(bv) < 12 {
		s_down = "BV号格式错误"
		http.Redirect(w, r, "/pa", http.StatusFound)
		return
	}
	if bv[0:2] != "BV" {
		s_down = "BV号格式错误"
		http.Redirect(w, r, "/pa", http.StatusFound)
		return
	}
	reque, _ := http.NewRequest("GET", "https://www.bilibili.com/video/"+bv, nil)
	client := http.Client{}
	reque.Header.Add("referer", "https://www.bilibili.com")
	reque.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36 Edg/94.0.992.38")
	rep, err := client.Do(reque)
	if err != nil {
		s_down = "解析出错"
		http.Redirect(w, r, "/pa", http.StatusFound)
		return
	}
	if rep.StatusCode != http.StatusOK {
		s_down = "解析出错"
		http.Redirect(w, r, "/pa", http.StatusFound)
		return
	}
	pages, err := jiexi_html.Get_aid_cid_pages("https://www.bilibili.com/video/" + bv)
	if err != nil {
		s_down = "解析出错"
		http.Redirect(w, r, "/pa", http.StatusFound)
		return
	}
	video_url_list, err := jiexi_html.Get_video_url_liet(bv, pages)
	if err != nil {
		s_down = "解析出错"
		http.Redirect(w, r, "/pa", http.StatusFound)
		return
	}
	m_urls[bv] = video_url_list
	s_list := ""
	for i := 0; i < len(video_url_list); i++ {
		s_list += `<div class="bv">        <form class="form" method="post" action="/pa_bv_post" onsubmit="return true">            <label>P` +
			strconv.Itoa(i+1) + `</label>            <input type="hidden" name="page" value="` +
			strconv.Itoa(i) + `" />` + `<input type="hidden" name="BV" value="` + bv + `" />          <button type="submit">下载</button>        </form>    </div>`
	}
	fmt.Println(s_list)
	tmpl, err := template.ParseFiles(path + "/tmpl/video_list.tmpl")
	if err != nil {
		s_down = "解析失败"
		http.Redirect(w, r, "/pa", http.StatusFound)
	}
	s_down = ""
	tmpl.Execute(w, map[string]interface{}{"s_list": template.HTML(s_list)})
	//tmpl.Execute(w, template.HTML(s_list))
}

func pa_bv_download(w http.ResponseWriter, r *http.Request) {
	if !admin {
		fmt.Println("不是管理员不给进")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	p_s := r.FormValue("page")
	p, _ := strconv.Atoi(p_s)
	bv := r.FormValue("BV")
	if len(bv) < 12 {
		s_down = "BV号格式错误"
		http.Redirect(w, r, "/pa", http.StatusFound)
		return
	}
	if bv[0:2] != "BV" {
		s_down = "BV号格式错误"
		http.Redirect(w, r, "/pa", http.StatusFound)
		return
	}
	p_url := m_urls[bv][p]
	//fmt.Println("p_url:",p_url)
	p_url = strings.Replace(p_url, `\u0026`, "&", -1)
	client := http.Client{}
	resp, err := http.NewRequest("GET", p_url, nil)
	resp.Header.Add("referer", "https://www.bilibili.com")
	resp.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36 Edg/95.0.1020.30")
	res, err := client.Do(resp)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("连接失败,网页状态码：：", res.StatusCode)
		s_down = "解析失败"
		http.Redirect(w, r, "/pa", http.StatusFound)
	}
	fmt.Println("连接成功,网页状态码：", res.StatusCode)
	w.Header().Set("Content-Disposition", "attachment; filename="+bv+p_s+".mp4")
	w.Header().Set("Content-Type", "video/flv")
	//w.Header().Set("Content-Length", res.Header.Get("Content-Length"))
	defer res.Body.Close()
	io.Copy(w, res.Body)
}

func pa_cv_list(w http.ResponseWriter, r *http.Request) {
	cv := r.FormValue("CV")
	p_cv_num = cv
	if cv[0:2] == "cv" || cv[0:2] == "CV" {
	} else {
		s_down = "CV号格式错误"
		http.Redirect(w, r, "/pa_picture", http.StatusFound)
		return
	}
	s := "https://www.bilibili.com/read/" + cv
	result := jiexi_html.GetURL(s)
	image_url_list := jiexi_html.Get_image_url_liet(result)
	//for _, s2 := range image_url_list {
	//	fmt.Println(s2)
	//}
	p_urls[cv] = image_url_list
	s_list := ""
	for i := 0; i < len(image_url_list); i++ {
		s_list += `<div class="cv">        <form class="form" method="post" action="/pa_cv_post" onsubmit="return true">            <label>P` +
			strconv.Itoa(i+1) + `</label>            <input type="hidden" name="page" value="` +
			strconv.Itoa(i) + `" />` + `<input type="hidden" name="CV" value="` + cv + `" />          <button type="submit">下载</button>        </form>    </div>`
	}
	fmt.Println(s_list)
	tmpl, err := template.ParseFiles(path + "/tmpl/Piture_list.tmpl")
	if err != nil {
		s_down = "解析失败"
		http.Redirect(w, r, "/pa_picture", http.StatusFound)
	}
	s_down = ""
	tmpl.Execute(w, map[string]interface{}{"s_list": template.HTML(s_list)})
	//tmpl.Execute(w, template.HTML(s_list))
}

// GetFileName 获得图片文件名
func GetFileName(ImageUrl string) string {
	re := regexp.MustCompile(`/(\w+\.((jpg)|(png)|(gif)|(bmp)|(webp)|(swf)))`)
	result := re.FindAllStringSubmatch(ImageUrl, -1)
	if len(result) > 0 {
		return result[0][2]
	} else {
		return "jpg"
	}
}

func pa_cv_download(w http.ResponseWriter, r *http.Request) {
	p_s := r.FormValue("page")
	p, _ := strconv.Atoi(p_s)
	cv := r.FormValue("CV")
	//if len(cv)<12{
	//	s_down="CV号格式错误"
	//	http.Redirect(w,r,"/pa_picture",http.StatusFound)
	//	return
	//}
	if cv[0:2] == "cv" || cv[0:2] == "CV" {
	} else {
		s_down = "CV号格式错误"
		http.Redirect(w, r, "/pa_picture", http.StatusFound)
		return
	}
	p_url := p_urls[cv][p]
	//fmt.Println("p_url:",p_url)
	filelast := GetFileName(p_url)
	if p_url[0:6] != "https:" {
		p_url = "https:" + p_url
	}
	client := http.Client{}
	resp, err := http.NewRequest("GET", p_url, nil)
	resp.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36 Edg/95.0.1020.30")
	res, err := client.Do(resp)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("连接失败,网页状态码：：", res.StatusCode)
		s_down = "解析失败"
		http.Redirect(w, r, "/pa_picture", http.StatusFound)
	}
	fmt.Println("连接成功,网页状态码：", res.StatusCode)
	w.Header().Set("Content-Disposition", "attachment; filename="+cv+p_s+"."+filelast)
	//w.Header().Set("Content-Type", "image/webp")
	//w.Header().Set("Content-Length", res.Header.Get("Content-Length"))
	defer res.Body.Close()
	io.Copy(w, res.Body)
}

//var i =0
//func Delt(w http.ResponseWriter,r *http.Request) {
//	Delt1(w,r)
//}
//func Delt1(w http.ResponseWriter,r *http.Request)  (http.ResponseWriter,http.Request) {
//		cv:=p_cv_num
//		if i<len(p_urls) {
//			imagURL:=p_urls[cv]
//			fmt.Println(imagURL[i])
//			pa_cv_downloadAll(w,r,imagURL[i],i)
//			i++
//		}
//		return Delt1(w,r)
//}
//func pa_cv_downloadAll(w http.ResponseWriter, r *http.Request,p_url string,p_s int) {
//	cv:=p_cv_num
//	k:=string(p_s)
//	filelast:=GetFileName(p_url)
//	if p_url[0:6]!="https:" {
//		p_url = "https:"+p_url
//	}
//	client := http.Client{}
//	resp, err := http.NewRequest("GET", p_url, nil)
//	resp.Header.Add("user-agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36 Edg/95.0.1020.30")
//	res, err := client.Do(resp)
//	if err != nil ||res.StatusCode!=200{
//		fmt.Println("连接失败,网页状态码：：", res.StatusCode)
//		s_down="解析失败"
//		http.Redirect(w,r,"/pa_picture",http.StatusFound)
//	}
//	fmt.Println("连接成功,网页状态码：", res.StatusCode)
//	defer res.Body.Close()
//	w.Header().Set("Content-Disposition", "attachment; filename="+cv+k+"."+filelast)
//	//w.Header().Set("Content-Type", "image/webp")
//	//w.Header().Set("Content-Length", res.Header.Get("Content-Length"))
//	io.Copy(w,res.Body)
//}
