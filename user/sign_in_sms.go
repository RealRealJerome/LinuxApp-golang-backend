package user

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"github.com/A5-golang-backend/sms"
	"net/http"
	"strings"
)

func HandleSignInSMS(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	phoneNum := parts[2]

	// 解析请求体中的用户信息
	var smsCode SMSCode
	err := json.NewDecoder(r.Body).Decode(&smsCode)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// 先判断存不存在
	selectQuery := "SELECT password FROM user WHERE phone_num = ?"
	selectedPW := ""
	err = mysqlUtil.DB.QueryRow(selectQuery, phoneNum).Scan(&selectedPW)
	if err == sql.ErrNoRows {
		// 构建响应
		response := struct {
			Code    int
			Message string
			Token   string `json:"token"`
		}{
			Code:    400,
			Message: "手机号不存在，请先注册",
			Token:   "",
		}
		// 4. 序列化返回体为 JSON 格式
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 返回JSON响应
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
		return
	}
	isMatched := sms.Verify(phoneNum, smsCode.SMSCode)
	if !isMatched {
		// 构建响应
		response := struct {
			Code    int
			Message string
			Token   string `json:"token"`
		}{
			Code:    400,
			Message: "验证码有误或已过期，请重新获取",
			Token:   "",
		}
		// 4. 序列化返回体为 JSON 格式
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 返回JSON响应
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
		return
	}
	// 生成一个模拟的访问令牌
	token := generateToken(phoneNum)
	// 构建响应
	response := struct {
		Code    int
		Message string
		Token   string `json:"token"`
	}{
		Code:    200,
		Message: "登录成功",
		Token:   token,
	}
	// 4. 序列化返回体为 JSON 格式
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
