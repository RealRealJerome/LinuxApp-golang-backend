package sms

import (
	"encoding/json"
	redisUtil "github.com/A5-golang-backend/redis"
	"io"
	"net/http"
	"strings"
)

type verifySmsCodeReq struct {
	SmsCode string `json:"smsCode"`
}

func VerifySmsCode(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	var req verifySmsCodeReq
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "参数解析错误", http.StatusBadRequest)
		return
	}
	code := req.SmsCode
	// 获取请求路径中的 phoneNum 参数
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	phoneNum := parts[2]
	if !Verify(phoneNum, code) {
		// 3. 构建返回体
		response := struct {
			Code    int    `json:"code"` // 是否成功响应
			Message string `json:"message"`
		}{
			Code:    400,
			Message: "验证码错误或已过期",
		}
		// 4. 序列化返回体为 JSON 格式
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
		return
	}
	// 3. 构建返回体
	response := struct {
		Code    int    `json:"code"` // 是否成功响应
		Message string `json:"message"`
	}{
		Code:    200,
		Message: "验证码校验通过",
	}
	// 4. 序列化返回体为 JSON 格式
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
func Verify(phoneNum, smsCode string) bool {
	res := redisUtil.RDB.Get(phoneNum)
	if res.Err() != nil || res.Val() != smsCode {
		return false
	}
	return true
}
