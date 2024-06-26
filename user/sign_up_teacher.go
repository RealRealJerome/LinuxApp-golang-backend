package user

import (
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"log"
	"net/http"
	"strings"
)

func HandleSignUpTeacher(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	phoneNum := parts[2]
	// 解析请求体中的用户信息
	var teacherInfo TeacherInfo
	err := json.NewDecoder(r.Body).Decode(&teacherInfo)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	if !ValidatePassword(teacherInfo.Password) {
		// 构建响应
		response := struct {
			Code    int
			Message string
		}{
			Code:    400,
			Message: "密码不符合规范",
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
	// 查询数据
	selectQuery := "SELECT id FROM user_detail WHERE uuid = ?"
	selectedId := 0
	err = mysqlUtil.DB.QueryRow(selectQuery, teacherInfo.TeacherId).Scan(&selectedId)
	// 工号存在
	if err != nil && selectedId != 0 {
		// 构建响应
		response := struct {
			Code    int
			Message string
		}{
			Code:    400,
			Message: "教职工号已存在！",
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
	// 查询数据
	selectQuery = "SELECT phone_num FROM user WHERE phone_num = ?"
	selectedPhoneNum := ""
	err = mysqlUtil.DB.QueryRow(selectQuery, phoneNum).Scan(&selectedPhoneNum)
	// 手机号存在
	if err != nil && selectedPhoneNum != "" {
		// 构建响应
		response := struct {
			Code    int
			Message string
		}{
			Code:    400,
			Message: "手机号已存在！",
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
	// 插入数据
	// 插入数据的 SQL 语句
	query := "INSERT INTO user_detail (uuid, name,college,school) VALUES (?, ?, ?, ?)"
	// 执行插入操作
	result, err := mysqlUtil.DB.Exec(query, teacherInfo.TeacherId,
		teacherInfo.Name, teacherInfo.College, teacherInfo.School)
	if err != nil {
		log.Fatal(err)
	}
	// 获取插入的记录的 ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	query = "INSERT INTO user (user_detail_id, phone_num,password) VALUES (?, ?, ?)"
	// 执行插入操作
	result, err = mysqlUtil.DB.Exec(query, lastInsertID,
		phoneNum, Sha256(teacherInfo.Password))
	if err != nil {
		log.Fatal(err)
	}
	// 构建响应
	response := struct {
		Code    int
		Message string
	}{
		Code:    200,
		Message: "注册成功",
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
