package http_tools

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// PostFormFile 准备一个您将提交的表单该网址。
func PostFormFile(url, filename string) (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	//添加您的镜像文件
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	fw, err := w.CreateFormFile("filename", f.Name())
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}

	//添加其他字段
	if fw, err = w.CreateFormField("key"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("KEY")); err != nil {
		return
	}
	//不要忘记关闭multipart writer。
	//如果你不关闭它,你的请求将丢失终止边界。
	w.Close()

	//现在你有一个表单,你可以提交它给你的处理程序。
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}

	// Don不要忘记设置内容类型,这将包含边界。
	req.Header.Set("Content-Type", w.FormDataContentType())

	//提交请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err == nil {
		return
	}

	//检查响应
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("status:％s", res.Status)
	}
	return
}
