package user

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"net/http"
	"strings"
)

var codeMsgMap = map[int]string{
	1: "手机号不存在",
	2: "旧密码不匹配",
	3: "新密码不符合规范",
	4: "插入数据失败：服务器内部错误",
	5: "密码修改成功！请重新登录",
}

func HandleChangePW(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	phoneNum := parts[2]
	vars := r.URL.Query()
	neededOld := vars.Get("needed_old")
	// 解析请求体中的用户信息
	var oldNewPW OldNewPW
	err := json.NewDecoder(r.Body).Decode(&oldNewPW)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	oldPW := oldNewPW.OldPassword
	newPW := oldNewPW.NewPassword
	if neededOld == "true" {
		res := ChPW(phoneNum, oldPW, newPW)
		if res == 4 {
			http.Error(w, codeMsgMap[4], http.StatusInternalServerError)
			return
		}
		if res == 5 {
			// 构建响应
			response := struct {
				Code    int
				Message string
			}{
				Code:    200,
				Message: codeMsgMap[5],
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
		// 构建响应
		response := struct {
			Code    int
			Message string
		}{
			Code:    400,
			Message: codeMsgMap[res],
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
	} else {
		res := ForgotPW(phoneNum, newPW)
		if res == 4 {
			http.Error(w, codeMsgMap[4], http.StatusInternalServerError)
			return
		}
		if res == 5 {
			// 构建响应
			response := struct {
				Code    int
				Message string
			}{
				Code:    200,
				Message: "找回密码成功，请重新登录",
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
		// 构建响应
		response := struct {
			Code    int
			Message string
		}{
			Code:    400,
			Message: codeMsgMap[res],
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
}

// 找回密码
func ForgotPW(phoneNum, newPW string) int {
	// 查询数据
	selectQuery := "SELECT password FROM user WHERE phone_num = ?"
	selectedPW := ""
	err := mysqlUtil.DB.QueryRow(selectQuery, phoneNum).Scan(&selectedPW)
	if err == sql.ErrNoRows {
		return 1
	}
	if !ValidatePassword(newPW) {
		return 3
	}
	query := "UPDATE user SET password = ? WHERE phone_num = ?"

	// 准备更新的数据
	password := Sha256(newPW)
	// 执行更新操作
	_, err = mysqlUtil.DB.Exec(query, password, phoneNum)
	if err != nil {
		return 4
	}
	return 5
}

// 1--手机号不存在，2--旧密码不匹配，3--新密码不符合规范，4--插入失败：服务器内部错误，5--OK
func ChPW(phoneNum, oldPW, newPW string) int {
	// 查询数据
	selectQuery := "SELECT password FROM user WHERE phone_num = ?"
	selectedPW := ""
	err := mysqlUtil.DB.QueryRow(selectQuery, phoneNum).Scan(&selectedPW)
	if err == sql.ErrNoRows {
		return 1
	}
	if Sha256(oldPW) != selectedPW {
		return 2
	}
	if !ValidatePassword(newPW) {
		return 3
	}
	query := "UPDATE user SET password = ? WHERE phone_num = ?"

	// 准备更新的数据
	password := Sha256(newPW)
	// 执行更新操作
	_, err = mysqlUtil.DB.Exec(query, password, phoneNum)
	if err != nil {
		return 4
	}
	return 5
}
