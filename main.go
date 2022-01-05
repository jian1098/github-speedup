package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("正在解析域名...")
	//需要解析的域名
	var host = "github.com"
	var filename = "C:\\Windows\\System32\\drivers\\etc\\hosts"

	//dns服务器id
	var guid = []string{
		"29d2a14f-accc-43b8-9444-fd6b9e7902bc",
		"0e215450-a287-486e-b758-49b00e432bd4",
		"38806ef7-4638-4808-96e4-85047dfa5853",
		"53e79941-312c-4343-8739-6be3cd105805",
		"e00f2cc8-3f66-452d-b57f-804f84b774ca",
		"97a2da5b-4aa4-4811-aff7-2eb038b6f742",
		"7ddde139-f06e-40aa-b585-25d2ee6fad5f",
		"d9041619-7d90-42ea-9811-2b2fe11cb2b0",
		"02a01d5d-5111-481f-aade-e999a584d8a4",
		"80a828bd-19ed-48c3-a035-e69f6468da03",
		"fc7b8db4-f81d-4432-8d27-fa43dd13df3c",
		"91937e5b-1db0-47b5-b114-c9294694f377",
		"1f4c5976-8cf3-47e7-be10-aa9270461477",
		"3c1c826d-3444-4350-849b-0b9b9755df78",
		"7a2c1fe7-f9de-4fee-b797-bfe343d49f15",
		"eac78784-07a1-4869-be7c-3870a8dcebfc",
		"87c200e0-0059-479f-8103-e9e504f735d0",
		"dc440a55-1148-480f-90a7-9d1e0269b682",
		"08117724-8437-4ebb-88ae-93e50f660867",
		"5fb9012d-b47c-4087-84a2-0b0dfa8c94ab",
		"4250e220-157f-4831-8e6b-ad7cead81ca0",
		"1e375923-e5ee-491e-ba21-621a95ef9de9",
		"39bed414-9402-4266-aa9a-8252e958558f",
		"a0be885d-24ad-487d-bbb0-c94cd02a137d",
		"29d2a14f-accc-43b8-9444-fd6b9e7902bc",
	}
	var ipArr []string
	var wg sync.WaitGroup
	for _, val := range guid {
		data := url.Values{}
		data.Set("guid", val)
		data.Set("host", host)
		wg.Add(1)
		go func() {
			res := httpPost("http://mping.chinaz.com/Handle/AjaxHandler.ashx?action=Ping", data)
			ip := getIp(res)
			if ip != "" {
				ipArr = append(ipArr, ip+" "+host)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	ipArr = arrayUnique(ipArr)

	content, err := readFile(filename, host)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, val := range ipArr {
		content += val + "\n"
	}
	//将获取到的ip写入hosts文件
	err = writeFile(filename, content)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入数据成功!")
	time.Sleep(time.Duration(time.Second * 3))
}

//post请求
func httpPost(link string, data url.Values) string {
	resp, err := http.PostForm(link, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

//获取ip地址
func getIp(str string) string {
	gexp := "^(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9])$"
	arr := strings.Split(str, "'")
	var ip string
	for _, str := range arr {
		match, _ := regexp.MatchString(gexp, str)
		if match {
			ip = str
			break
		}
	}
	return ip
}

//数组去重
func arrayUnique(arr []string) []string {
	size := len(arr)
	result := make([]string, 0, size)
	temp := map[string]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[arr[i]]; ok != true {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}
	return result
}

//读取文件内容
func readFile(filename, hosts string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content string //用来保存文件内容
	for scanner.Scan() {
		// 读取当前行内容
		line := scanner.Text()
		if strings.Contains(line, hosts) {
			continue
		}
		content += line + "\n"
	}
	return content, nil
}

//写入文件
func writeFile(filename, content string) error {
	//覆盖写入
	err := ioutil.WriteFile(filename, []byte(content), 0664)
	if err != nil {
		return err
	}
	return nil
}
