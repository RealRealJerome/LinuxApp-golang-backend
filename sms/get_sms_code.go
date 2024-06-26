package sms

import (
	"encoding/json"
	redisUtil "github.com/A5-golang-backend/redis"
	"io"
	"net/http"
	"strings"
	"time"
)

type smsCodeResp struct {
	phoneNum string
	code     string
}

func rpcGetSmsCode(phoneNum string) (*smsCodeResp, error) {
	// 设置要调用的API的URL
	url := "http://47.103.119.74:8888/sms?phoneNum=" + phoneNum
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// 读取响应的body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	strBody := string(body)
	res := strings.Split(strBody, ":")
	return &smsCodeResp{
		phoneNum: res[0],
		code:     res[1],
	}, err
}
func GetSmsCode(w http.ResponseWriter, r *http.Request) {
	// 获取请求路径中的 phoneNum 参数
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	phoneNum := parts[2]
	res, err := rpcGetSmsCode(phoneNum)
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	// 3. 构建返回体
	response := struct {
		Code    int    `json:"code"` // 是否成功响应
		Message string `json:"message"`
	}{
		Code:    200,
		Message: "短信验证码已发送，可能会有延后，请耐心等待",
	}
	// 4. 序列化返回体为 JSON 格式
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	redisUtil.RDB.Set(res.phoneNum, res.code, 60*time.Second)
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
