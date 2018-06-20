// duanzi.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	//	test()

	var start, end int
	fmt.Printf("请输入开始页：")
	fmt.Scan(&start)
	fmt.Printf("请输入结束页：")
	fmt.Scan(&end)

	DoWork(start, end)
}

func DoWork(start, end int) {
	fmt.Printf("正在爬去 %d 到 %d 的页面\n", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}

	for i := start; i <= end; i++ {

		fmt.Printf("第%v个页面爬取完成\n", <-page)
	}
}

func SpiderPage(i int, page chan int) {
	url := "https://www.pengfu.com/xiaohua_" + strconv.Itoa(i) + ".html"
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("httpGet err = ", err)
		return
	}
	//	fmt.Println(result)
	//<h1 class="dp-b"><a href="https://www.pengfu.com/content_1829547_1.html"
	//
	//	tmpStr := "<h1 class='dp-b'><a href='(?s:(.*?))'" //
	r1 := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))"`) //正则表达式
	if r1 == nil {
		fmt.Println("regexp.MustCompile err")
		return
	}
	joyUrls := r1.FindAllStringSubmatch(result, -1)

	fileTitle := make([]string, 0)
	fileContent := make([]string, 0)

	for _, data := range joyUrls {
		//		fmt.Println(data[1])
		title, content, err := SpiderOneJoy(data[1])
		if err != nil {
			fmt.Println("SpiderOneJoy err = ", err)
		}
		fileTitle = append(fileTitle, title)
		fileContent = append(fileContent, content)

	}

	StoreJoyToFile(i, fileTitle, fileContent)
	page <- i
}

func StoreJoyToFile(i int, fileTitle, fileContent []string) {
	f, err := os.Create(strconv.Itoa(i) + ".txt")
	if err != nil {
		fmt.Println("os.Create err = ", err)
	}
	defer f.Close()

	n := len(fileContent)
	for i = 0; i < n; i++ {
		f.WriteString(fileTitle[i] + "\n")
		f.WriteString(fileContent[i] + "\n")
		f.WriteString("================================================================\n")
	}

}

func SpiderOneJoy(url string) (title, content string, err error) {
	result, err1 := HttpGet(url)
	if err1 != nil {
		err = err1
		return
	}
	r1 := regexp.MustCompile(`<h1>(?s:(.*?))</h1>`) //正则表达式
	if r1 == nil {
		err = fmt.Errorf("%s", "regexp.MustCompile err")
		return
	}
	tmpTitle := r1.FindAllStringSubmatch(result, 1)
	for _, data := range tmpTitle {
		title = data[1]
		title = strings.Replace(title, " ", "", -1)
		title = strings.Replace(title, "\t", "", -1)
		title = strings.Replace(title, "\n", "", -1)
		break
	}

	r2 := regexp.MustCompile(`<div class="content-txt pt10">(?s:(.*?))<a id="prev"`) //正则表达式
	if r2 == nil {
		err = fmt.Errorf("%s", "regexp.MustCompile err")
		return
	}
	tmpContent := r2.FindAllStringSubmatch(result, -1)
	for _, data := range tmpContent {
		content = data[1]
		content = strings.Replace(content, " ", "", -1)
		content = strings.Replace(content, "\t", "", -1)
		content = strings.Replace(content, "\n", "", -1)
		break
	}
	return
}

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4*1024)

	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}

	return

}
